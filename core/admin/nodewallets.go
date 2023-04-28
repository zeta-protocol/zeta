// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.ZETA file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package admin

import (
	"fmt"
	"net/http"

	"github.com/zeta-protocol/zeta/core/nodewallets"
	"github.com/zeta-protocol/zeta/core/nodewallets/eth/clef"
	"github.com/zeta-protocol/zeta/core/nodewallets/eth/keystore"
	"github.com/zeta-protocol/zeta/core/nodewallets/registry"
	"github.com/zeta-protocol/zeta/libs/crypto"
	"github.com/zeta-protocol/zeta/logging"
	"github.com/zeta-protocol/zeta/paths"
)

type wallet interface {
	Name() string
	PubKey() crypto.PublicKey
}

type Wallet struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

func newWallet(w wallet) Wallet {
	return Wallet{
		Name:      w.Name(),
		PublicKey: w.PubKey().Hex(),
	}
}

type NodeWalletArgs struct {
	Chain string
}

type NodeWalletReloadReply struct {
	OldWallet Wallet `json:"oldWallet"`
	NewWallet Wallet `json:"newWallet"`
}

type NodeWallet struct {
	log                  *logging.Logger
	nodeWallets          *nodewallets.NodeWallets
	registryLoader       *registry.Loader
	nodeWalletPassphrase string
	zetaPaths            paths.Paths
}

func NewNodeWallet(
	log *logging.Logger,
	zetaPaths paths.Paths,
	nodeWalletPassphrase string,
	nodeWallets *nodewallets.NodeWallets,
) (*NodeWallet, error) {
	registryLoader, err := registry.NewLoader(zetaPaths, nodeWalletPassphrase)
	if err != nil {
		return nil, err
	}

	return &NodeWallet{
		log:                  log.Named(nodeWalletNamedLogger),
		nodeWallets:          nodeWallets,
		registryLoader:       registryLoader,
		nodeWalletPassphrase: nodeWalletPassphrase,
		zetaPaths:            zetaPaths,
	}, nil
}

func (nw *NodeWallet) Reload(_ *http.Request, args *NodeWalletArgs, reply *NodeWalletReloadReply) error {
	nw.log.Info("Reloading node wallet", logging.String("chain", args.Chain))

	switch args.Chain {
	case "zeta":
		oW := newWallet(nw.nodeWallets.Zeta)

		reg, err := nw.registryLoader.Get(nw.nodeWalletPassphrase)
		if err != nil {
			return fmt.Errorf("couldn't load node wallet registry: %v", err)
		}

		if err := nw.nodeWallets.Zeta.Reload(*reg.Zeta); err != nil {
			nw.log.Error("Reloading node wallet failed", logging.Error(err))
			return fmt.Errorf("failed to reload Zeta wallet: %w", err)
		}

		nW := newWallet(nw.nodeWallets.Zeta)

		reply.NewWallet = nW
		reply.OldWallet = oW

		nw.log.Info("Reloaded node wallet", logging.String("chain", args.Chain))
		return nil
	case "ethereum":
		oW := newWallet(nw.nodeWallets.Ethereum)

		reg, err := nw.registryLoader.Get(nw.nodeWalletPassphrase)
		if err != nil {
			return fmt.Errorf("couldn't load node wallet registry: %v", err)
		}

		algoType := nw.nodeWallets.Ethereum.Algo()
		_, isKeyStoreWallet := reg.Ethereum.Details.(registry.EthereumKeyStoreWallet)
		_, isClefWallet := reg.Ethereum.Details.(registry.EthereumClefWallet)

		if isKeyStoreWallet && algoType != keystore.KeyStoreAlgoType {
			w, err := nodewallets.GetEthereumWalletWithRegistry(nw.zetaPaths, reg)
			if err != nil {
				return fmt.Errorf("failed reload key: %w", err)
			}

			nw.nodeWallets.SetEthereumWallet(w)
		} else if isClefWallet && algoType != clef.ClefAlgoType {
			w, err := nodewallets.GetEthereumWalletWithRegistry(
				nw.zetaPaths,
				reg,
			)
			if err != nil {
				return fmt.Errorf("failed reload key: %w", err)
			}

			nw.nodeWallets.SetEthereumWallet(w)
		} else {
			if err := nw.nodeWallets.Ethereum.Reload(reg.Ethereum.Details); err != nil {
				nw.log.Error("Reloading node wallet failed", logging.Error(err))
				return fmt.Errorf("failed to reload Ethereum wallet: %w", err)
			}
		}

		nW := newWallet(nw.nodeWallets.Ethereum)

		reply.NewWallet = nW
		reply.OldWallet = oW

		nw.log.Info("Reloaded node wallet", logging.String("chain", args.Chain))
		return nil
	}

	return fmt.Errorf("failed to reload wallet for non existing chain %q", args.Chain)
}

func (nw *NodeWallet) Show(_ *http.Request, args *NodeWalletArgs, reply *Wallet) error {
	switch args.Chain {
	case "zeta":
		*reply = newWallet(nw.nodeWallets.Zeta)
		return nil
	case "ethereum":
		*reply = newWallet(nw.nodeWallets.Ethereum)
		return nil
	}

	return fmt.Errorf("failed to show wallet for non existing chain %q", args.Chain)
}
