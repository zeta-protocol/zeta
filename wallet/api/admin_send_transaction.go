package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/zeta-protocol/zeta/commands"
	vgcrypto "github.com/zeta-protocol/zeta/libs/crypto"
	"github.com/zeta-protocol/zeta/libs/jsonrpc"
	apipb "github.com/zeta-protocol/zeta/protos/zeta/api/v1"
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
	walletpb "github.com/zeta-protocol/zeta/protos/zeta/wallet/v1"
	"github.com/zeta-protocol/zeta/wallet/api/node"
	wcommands "github.com/zeta-protocol/zeta/wallet/commands"
	"github.com/golang/protobuf/jsonpb"
	"github.com/mitchellh/mapstructure"
)

type AdminSendTransactionParams struct {
	Wallet      string      `json:"wallet"`
	PublicKey   string      `json:"publicKey"`
	Network     string      `json:"network"`
	NodeAddress string      `json:"nodeAddress"`
	Retries     uint64      `json:"retries"`
	SendingMode string      `json:"sendingMode"`
	Transaction interface{} `json:"transaction"`
}

type ParsedAdminSendTransactionParams struct {
	Wallet         string
	PublicKey      string
	Network        string
	NodeAddress    string
	Retries        uint64
	SendingMode    apipb.SubmitTransactionRequest_Type
	RawTransaction string
}

type AdminSendTransactionResult struct {
	ReceivedAt time.Time               `json:"receivedAt"`
	SentAt     time.Time               `json:"sentAt"`
	TxHash     string                  `json:"transactionHash"`
	Tx         *commandspb.Transaction `json:"transaction"`
	Node       AdminNodeInfoResult     `json:"node"`
}

type AdminNodeInfoResult struct {
	Host string `json:"host"`
}

type AdminSendTransaction struct {
	walletStore         WalletStore
	networkStore        NetworkStore
	nodeSelectorBuilder NodeSelectorBuilder
}

func (h *AdminSendTransaction) Handle(ctx context.Context, rawParams jsonrpc.Params) (jsonrpc.Result, *jsonrpc.ErrorDetails) {
	params, err := validateAdminSendTransactionParams(rawParams)
	if err != nil {
		return nil, invalidParams(err)
	}

	receivedAt := time.Now()

	if exist, err := h.walletStore.WalletExists(ctx, params.Wallet); err != nil {
		return nil, internalError(fmt.Errorf("could not verify the wallet exists: %w", err))
	} else if !exist {
		return nil, invalidParams(ErrWalletDoesNotExist)
	}

	alreadyUnlocked, err := h.walletStore.IsWalletAlreadyUnlocked(ctx, params.Wallet)
	if err != nil {
		return nil, internalError(fmt.Errorf("could not verify whether the wallet is already unlock or not: %w", err))
	}
	if !alreadyUnlocked {
		return nil, requestNotPermittedError(ErrWalletIsLocked)
	}

	w, err := h.walletStore.GetWallet(ctx, params.Wallet)
	if err != nil {
		return nil, internalError(fmt.Errorf("could not retrieve the wallet: %w", err))
	}
	request := &walletpb.SubmitTransactionRequest{}
	if err := jsonpb.Unmarshal(strings.NewReader(params.RawTransaction), request); err != nil {
		return nil, invalidParams(fmt.Errorf("the transaction does not use a valid Zeta command: %w", err))
	}

	request.PubKey = params.PublicKey
	request.Propagate = true
	if errs := wcommands.CheckSubmitTransactionRequest(request); !errs.Empty() {
		return nil, invalidParams(errs)
	}

	currentNode, errDetails := h.getNode(ctx, params)
	if errDetails != nil {
		return nil, errDetails
	}

	lastBlockData, errDetails := h.getLastBlockDataFromNetwork(ctx, currentNode)
	if errDetails != nil {
		return nil, errDetails
	}

	marshaledInputData, err := wcommands.ToMarshaledInputData(request, lastBlockData.BlockHeight)
	if err != nil {
		return nil, internalError(fmt.Errorf("could not marshal the input data: %w", err))
	}

	signature, err := w.SignTx(params.PublicKey, commands.BundleInputDataForSigning(marshaledInputData, lastBlockData.ChainID))
	if err != nil {
		return nil, internalError(fmt.Errorf("could not sign the transaction: %w", err))
	}

	// Build the transaction.
	tx := commands.NewTransaction(params.PublicKey, marshaledInputData, &commandspb.Signature{
		Value:   signature.Value,
		Algo:    signature.Algo,
		Version: signature.Version,
	})

	// Generate the proof of work for the transaction.
	txID := vgcrypto.RandomHash()
	powNonce, _, err := vgcrypto.PoW(lastBlockData.BlockHash, txID, uint(lastBlockData.ProofOfWorkDifficulty), lastBlockData.ProofOfWorkHashFunction)
	if err != nil {
		return nil, internalError(fmt.Errorf("could not compute the proof-of-work: %w", err))
	}

	tx.Pow = &commandspb.ProofOfWork{
		Nonce: powNonce,
		Tid:   txID,
	}

	sentAt := time.Now()
	txHash, err := currentNode.SendTransaction(ctx, tx, params.SendingMode)
	if err != nil {
		return nil, networkErrorFromTransactionError(err)
	}

	return AdminSendTransactionResult{
		ReceivedAt: receivedAt,
		SentAt:     sentAt,
		TxHash:     txHash,
		Tx:         tx,
		Node: AdminNodeInfoResult{
			Host: currentNode.Host(),
		},
	}, nil
}

func (h *AdminSendTransaction) getNode(ctx context.Context, params ParsedAdminSendTransactionParams) (node.Node, *jsonrpc.ErrorDetails) {
	var hosts []string
	var retries uint64
	if len(params.Network) != 0 {
		exists, err := h.networkStore.NetworkExists(params.Network)
		if err != nil {
			return nil, internalError(fmt.Errorf("could not determine if the network exists: %w", err))
		} else if !exists {
			return nil, invalidParams(ErrNetworkDoesNotExist)
		}

		n, err := h.networkStore.GetNetwork(params.Network)
		if err != nil {
			return nil, internalError(fmt.Errorf("could not retrieve the network configuration: %w", err))
		}

		if err := n.EnsureCanConnectGRPCNode(); err != nil {
			return nil, invalidParams(ErrNetworkConfigurationDoesNotHaveGRPCNodes)
		}
		hosts = n.API.GRPC.Hosts
		retries = n.API.GRPC.Retries
	} else {
		hosts = []string{params.NodeAddress}
		retries = params.Retries
	}

	nodeSelector, err := h.nodeSelectorBuilder(hosts, retries)
	if err != nil {
		return nil, internalError(fmt.Errorf("could not initialize the node selector: %w", err))
	}

	currentNode, err := nodeSelector.Node(ctx, noNodeSelectionReporting)
	if err != nil {
		return nil, nodeCommunicationError(ErrNoHealthyNodeAvailable)
	}

	return currentNode, nil
}

func (h *AdminSendTransaction) getLastBlockDataFromNetwork(ctx context.Context, node node.Node) (*AdminLastBlockData, *jsonrpc.ErrorDetails) {
	lastBlock, err := node.LastBlock(ctx)
	if err != nil {
		return nil, nodeCommunicationError(ErrCouldNotGetLastBlockInformation)
	}

	if lastBlock.ChainID == "" {
		return nil, nodeCommunicationError(ErrCouldNotGetChainIDFromNode)
	}

	return &AdminLastBlockData{
		BlockHash:               lastBlock.BlockHash,
		ChainID:                 lastBlock.ChainID,
		BlockHeight:             lastBlock.BlockHeight,
		ProofOfWorkHashFunction: lastBlock.ProofOfWorkHashFunction,
		ProofOfWorkDifficulty:   lastBlock.ProofOfWorkDifficulty,
	}, nil
}

func NewAdminSendTransaction(walletStore WalletStore, networkStore NetworkStore, nodeSelectorBuilder NodeSelectorBuilder) *AdminSendTransaction {
	return &AdminSendTransaction{
		walletStore:         walletStore,
		networkStore:        networkStore,
		nodeSelectorBuilder: nodeSelectorBuilder,
	}
}

func validateAdminSendTransactionParams(rawParams jsonrpc.Params) (ParsedAdminSendTransactionParams, error) {
	if rawParams == nil {
		return ParsedAdminSendTransactionParams{}, ErrParamsRequired
	}

	params := AdminSendTransactionParams{}
	if err := mapstructure.Decode(rawParams, &params); err != nil {
		return ParsedAdminSendTransactionParams{}, ErrParamsDoNotMatch
	}

	if params.Wallet == "" {
		return ParsedAdminSendTransactionParams{}, ErrWalletIsRequired
	}

	if params.PublicKey == "" {
		return ParsedAdminSendTransactionParams{}, ErrPublicKeyIsRequired
	}

	if params.Network == "" && params.NodeAddress == "" {
		return ParsedAdminSendTransactionParams{}, ErrNetworkOrNodeAddressIsRequired
	}

	if params.Network != "" && params.NodeAddress != "" {
		return ParsedAdminSendTransactionParams{}, ErrSpecifyingNetworkAndNodeAddressIsNotSupported
	}

	if params.Transaction == nil || params.Transaction == "" {
		return ParsedAdminSendTransactionParams{}, ErrTransactionIsRequired
	}

	tx, err := json.Marshal(params.Transaction)
	if err != nil {
		return ParsedAdminSendTransactionParams{}, ErrTransactionIsNotValidJSON
	}

	return ParsedAdminSendTransactionParams{
		Wallet:         params.Wallet,
		PublicKey:      params.PublicKey,
		RawTransaction: string(tx),
		Network:        params.Network,
		NodeAddress:    params.NodeAddress,
		Retries:        params.Retries,
	}, nil
}
