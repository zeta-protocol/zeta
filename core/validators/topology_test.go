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

package validators_test

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	bmocks "github.com/zeta-protocol/zeta/core/broker/mocks"
	"github.com/zeta-protocol/zeta/core/events"
	"github.com/zeta-protocol/zeta/core/nodewallets"
	"github.com/zeta-protocol/zeta/core/validators"
	"github.com/zeta-protocol/zeta/core/validators/mocks"
	"github.com/zeta-protocol/zeta/libs/crypto"
	vgcrypto "github.com/zeta-protocol/zeta/libs/crypto"
	vgtesting "github.com/zeta-protocol/zeta/libs/testing"
	"github.com/zeta-protocol/zeta/logging"
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	types1 "github.com/tendermint/tendermint/proto/tendermint/types"
)

var tmPubKey = "tm-pub-key"

type NodeWallets struct {
	zeta             validators.Wallet
	tendermintPubkey string
	ethereumAddress  string
	ethereum         validators.Signer
}

func (n *NodeWallets) GetZeta() validators.Wallet {
	return n.zeta
}

func (n *NodeWallets) GetTendermintPubkey() string {
	return n.tendermintPubkey
}

func (n *NodeWallets) GetEthereumAddress() string {
	return n.ethereumAddress
}

func (n *NodeWallets) GetEthereum() validators.Signer {
	return n.ethereum
}

type DummyMultiSigTopology struct{}

func (*DummyMultiSigTopology) IsSigner(address string) bool {
	return true
}

func (*DummyMultiSigTopology) ExcessSigners(addresses []string) bool {
	return false
}

func (*DummyMultiSigTopology) GetThreshold() uint32 {
	return 666
}

func (*DummyMultiSigTopology) GetSigners() []string {
	return []string{}
}

type testTop struct {
	*validators.Topology
	ctrl        *gomock.Controller
	wallet      *mocks.MockWallet
	broker      *bmocks.MockBroker
	timeService *mocks.MockTimeService
	multisigTop *mocks.MockMultiSigTopology
}

func getTestTopologyWithNodeWallet(
	t *testing.T, wallet *mocks.MockWallet, nw *NodeWallets, ctrl *gomock.Controller,
) *testTop {
	t.Helper()

	broker := bmocks.NewMockBroker(ctrl)
	timeService := mocks.NewMockTimeService(ctrl)
	broker.EXPECT().Send(gomock.Any()).AnyTimes()
	multisigTop := mocks.NewMockMultiSigTopology(ctrl)

	commander := mocks.NewMockCommander(gomock.NewController(t))

	top := validators.NewTopology(logging.NewTestLogger(), validators.NewDefaultConfig(), nw, broker, true, commander, multisigTop, timeService)
	return &testTop{
		Topology:    top,
		ctrl:        ctrl,
		wallet:      wallet,
		broker:      broker,
		timeService: timeService,
		multisigTop: multisigTop,
	}
}

func getTestTopology(t *testing.T) *testTop {
	t.Helper()
	ctrl := gomock.NewController(t)
	dummyPubKey := "iamapubkey"
	pubKey := crypto.NewPublicKey(dummyPubKey, []byte(dummyPubKey))

	wallet := mocks.NewMockWallet(ctrl)
	wallet.EXPECT().PubKey().Return(pubKey).AnyTimes()
	wallet.EXPECT().ID().Return(pubKey).AnyTimes()

	nw := &NodeWallets{
		zeta:             wallet,
		tendermintPubkey: "rlg/jtPcVSdV23oFX8828sYFD84d7QsPt12YpiQH3Zw=",
		ethereumAddress:  "0x5cd0ec63687588817044794bf15d4e37991efab3",
	}

	return getTestTopologyWithNodeWallet(t, wallet, nw, ctrl)
}

func getTestTopWithDefaultValidator(t *testing.T) *testTop {
	t.Helper()

	top := getTestTopology(t)

	// Add Tendermint public key to validator set

	defaultTmPubKey := "default-tm-public-key"
	defaultTmPubKeyBase64 := base64.StdEncoding.EncodeToString([]byte(defaultTmPubKey))

	state := struct {
		Validators map[string]validators.ValidatorData
	}{
		Validators: map[string]validators.ValidatorData{
			defaultTmPubKeyBase64: {
				ID:              top.wallet.PubKey().Hex(),
				ZetaPubKey:      top.wallet.PubKey().Hex(),
				TmPubKey:        "asdasd",
				EthereumAddress: "0x123456",
				InfoURL:         "n0.xyz.zeta/node/123",
				Country:         "GB",
			},
		},
	}

	buf, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("error marshalling state %v", err)
	}

	if err := top.LoadValidatorsOnGenesis(context.Background(), buf); err != nil {
		t.Fatalf("error loading validators on genesis: %v", err)
	}

	return top
}

func getTestTopologyWithSelfValidatorData(
	t *testing.T, self validators.ValidatorData,
) *testTop {
	t.Helper()

	ctrl := gomock.NewController(t)
	pubKey := crypto.NewPublicKey(self.ZetaPubKey, []byte(self.ZetaPubKey))
	id := crypto.NewPublicKey(self.ID, []byte(self.ZetaPubKey))

	wallet := mocks.NewMockWallet(ctrl)
	wallet.EXPECT().PubKey().Return(pubKey).AnyTimes()
	wallet.EXPECT().ID().Return(id).AnyTimes()
	nw := &NodeWallets{
		zeta:             wallet,
		tendermintPubkey: self.TmPubKey,
		ethereumAddress:  self.EthereumAddress,
	}

	return getTestTopologyWithNodeWallet(t, wallet, nw, ctrl)
}

func loadGenesisValidators(
	t *testing.T, top *testTop, data ...validators.ValidatorData,
) error {
	t.Helper()
	state := struct {
		Validators map[string]validators.ValidatorData
	}{
		Validators: map[string]validators.ValidatorData{},
	}

	for _, v := range data {
		state.Validators[v.TmPubKey] = v
	}

	buf, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("error marshalling state %v", err)
	}

	return top.LoadValidatorsOnGenesis(context.Background(), buf)
}

func TestValidatorTopology(t *testing.T) {
	t.Run("add node registration - success", testAddNewNodeSuccess)
	t.Run("add node registration - failure", testAddNewNodeFailure)
	t.Run("test add node registration send event to broker", testAddNewNodeSendsValidatorUpdateEventToBroker)
	t.Run("topology validators length is equal to number of added validators", testGetLen)
	t.Run("added validator exists in topology", testExists)
	t.Run("test get by key", testGetByKey)
	t.Run("test validators validations", testValidatorsValidation)
}

func testValidatorsValidation(t *testing.T) {
	self := validators.ValidatorData{
		ID:              "f42b834d75f9ecb7b8167277fdae6ff664085d69588c508ada655d7876961558",
		ZetaPubKey:      "6a8325087e5bdf57b60cf06c3764e3c6a32840079fdc432a437ce32cd99316b5",
		TmPubKey:        "rlg/jtPcVSdV23oFX8828sYFD84d7QsPt12YpiQH3Zw=",
		EthereumAddress: "0x5cd0ec63687588817044794bf15d4e37991efab3",
	}
	otherValidators := []validators.ValidatorData{
		{
			ID:              "4f69b1784656174e89eb094513b7136e88670b42517ed0e48cb6fd3062eb8478",
			ZetaPubKey:      "f4686749895bf51c6df4092ef6be4279c384a3c380c24ea7a2fd20afc602a35d",
			TmPubKey:        "uBr9FP/M/QyVtOa3j18+hjksXra7qxCa7e25/FVW5c0=",
			EthereumAddress: "0xF3920d9Ab483177C99846503A118fa84A557bB27",
		},
		{
			ID:              "74023df02b8afc9eaf3e3e2e8b07eab1d2122ac3e74b1b0222daf4af565ad3dd",
			ZetaPubKey:      "10b06fec6398d9e9d542d7b7d36933a1e6f0bb0631b0e532681c05123d4bd5aa",
			TmPubKey:        "hz528OlxLZoV+476oJP2lzrhAZwZNjjLAfvpd2wLvcg=",
			EthereumAddress: "0x1b79814f66773df25ba126E8d1A557ab2676246f",
		},
	}

	testSuite := []struct {
		name          string
		self          validators.ValidatorData
		others        []validators.ValidatorData
		expectFailure bool
	}{
		{
			name:          "node setup correct",
			self:          self,
			others:        append(otherValidators, self),
			expectFailure: false,
		},
		{
			name: "node setup incorrect, invalid ID",
			self: validators.ValidatorData{
				TmPubKey: "rlg/jtPcVSdV23oFX8828sYFD84d7QsPt12YpiQH3Zw=",

				ID:              "INVALID-ID",
				ZetaPubKey:      "6a8325087e5bdf57b60cf06c3764e3c6a32840079fdc432a437ce32cd99316b5",
				EthereumAddress: "0x5cd0ec63687588817044794bf15d4e37991efab3",
			},
			others:        append(otherValidators, self),
			expectFailure: true,
		},
		{
			name: "node setup correct, invalid pubkey",
			self: validators.ValidatorData{
				ID:              "f42b834d75f9ecb7b8167277fdae6ff664085d69588c508ada655d7876961558",
				ZetaPubKey:      "INVALID-PUBKEY",
				TmPubKey:        "rlg/jtPcVSdV23oFX8828sYFD84d7QsPt12YpiQH3Zw=",
				EthereumAddress: "0x5cd0ec63687588817044794bf15d4e37991efab3",
			},
			others:        append(otherValidators, self),
			expectFailure: true,
		},
		{
			name: "node setup incorrect, invalid ethereum address",
			self: validators.ValidatorData{
				ID:              "f42b834d75f9ecb7b8167277fdae6ff664085d69588c508ada655d7876961558",
				ZetaPubKey:      "6a8325087e5bdf57b60cf06c3764e3c6a32840079fdc432a437ce32cd99316b5",
				TmPubKey:        "rlg/jtPcVSdV23oFX8828sYFD84d7QsPt12YpiQH3Zw=",
				EthereumAddress: "0xNOPE",
			},
			others:        append(otherValidators, self),
			expectFailure: true,
		},
		{
			name: "node setup inccorrect, all invalid",
			self: validators.ValidatorData{
				ID:              "WRONG",
				ZetaPubKey:      "BAD",
				TmPubKey:        "rlg/jtPcVSdV23oFX8828sYFD84d7QsPt12YpiQH3Zw=",
				EthereumAddress: "0xNOPE",
			},
			others:        append(otherValidators, self),
			expectFailure: true,
		},
	}
	for _, set := range testSuite {
		t.Run(set.name, func(t *testing.T) {
			// one validator -> self, 2 non validators
			top := getTestTopologyWithSelfValidatorData(t, set.self)
			if set.expectFailure {
				assert.Panics(t, func() {
					loadGenesisValidators(t, top, set.others...)
				})
			} else {
				assert.NotPanics(t, func() {
					err := loadGenesisValidators(t, top, set.others...)
					assert.NoError(t, err)
				})
			}
			top.ctrl.Finish()
		})
	}
}

func testAddNewNodeSuccess(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()

	nr := commandspb.AnnounceNode{
		Id:              "zeta-master-pubkey",
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      "zeta-key",
		EthereumAddress: "eth-address",
	}
	ctx := context.Background()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)
}

func testAddNewNodeFailure(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()

	nr := commandspb.AnnounceNode{
		Id:              "zeta-master-pubkey",
		ChainPubKey:     "tm-pub-key-1",
		ZetaPubKey:      "zeta-key",
		EthereumAddress: "eth-address",
	}
	ctx := context.Background()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	// Add node with existing ZetaPubKey
	nr = commandspb.AnnounceNode{
		Id:              "zeta-master-pubkey",
		ChainPubKey:     "tm-pub-key-2",
		ZetaPubKey:      "zeta-key",
		EthereumAddress: "eth-address-2",
	}
	err = top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.Error(t, err)
}

func testGetLen(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()

	// first the len is 1 since the default validator loaded from genenesis
	assert.Equal(t, 1, top.Len())

	nr := commandspb.AnnounceNode{
		Id:              "zeta-master-pubkey",
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      "zeta-key",
		EthereumAddress: "eth-address",
	}
	ctx := context.Background()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	assert.Equal(t, 2, top.Len())
}

func testExists(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()

	assert.False(t, top.IsValidatorZetaPubKey("zeta-key"))
	assert.False(t, top.IsValidatorNodeID("zeta-master-pubkey"))

	nr := commandspb.AnnounceNode{
		Id:              "zeta-master-pubkey",
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      "zeta-key",
		EthereumAddress: "eth-address",
	}
	ctx := context.Background()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	assert.True(t, top.IsValidatorZetaPubKey("zeta-key"))
	assert.True(t, top.IsValidatorNodeID("zeta-master-pubkey"))
}

func testGetByKey(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()

	assert.False(t, top.IsValidatorZetaPubKey("zeta-key"))
	assert.False(t, top.IsValidatorNodeID("zeta-master-pubkey"))

	nr := commandspb.AnnounceNode{
		Id:              "zeta-master-pubkey",
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      "zeta-key",
		EthereumAddress: "eth-address",
		InfoUrl:         "n0.xyz.zeta/node/url/random",
		Country:         "CZ",
	}
	ctx := context.Background()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	expectedData := &validators.ValidatorData{
		ID:              "zeta-master-pubkey",
		ZetaPubKey:      nr.ZetaPubKey,
		EthereumAddress: "eth-address",
		TmPubKey:        nr.ChainPubKey,
		InfoURL:         nr.InfoUrl,
		Country:         nr.Country,
	}

	actualData := top.Get(nr.Id)
	assert.NotNil(t, actualData)

	assert.Equal(t, expectedData, actualData)
}

func testAddNewNodeSendsValidatorUpdateEventToBroker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	_, err := nodewallets.GenerateZetaWallet(zetaPaths, "pass", "pass", false)
	require.NoError(t, err)
	wallet, err := nodewallets.GetZetaWallet(zetaPaths, "pass")
	require.NoError(t, err)

	nw := &NodeWallets{
		zeta: wallet,
	}

	broker := bmocks.NewMockBroker(ctrl)
	timeService := mocks.NewMockTimeService(ctrl)
	commander := mocks.NewMockCommander(gomock.NewController(t))
	top := validators.NewTopology(logging.NewTestLogger(), validators.NewDefaultConfig(), nw, broker, true, commander, &DummyMultiSigTopology{}, timeService)

	ctx := context.Background()
	nr := commandspb.AnnounceNode{
		Id:              "zeta-master-pubkey",
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      "zeta-key",
		EthereumAddress: "eth-address",
		InfoUrl:         "n0.xyz.zeta/node/url/random",
		Country:         "CZ",
		Name:            "validator",
		AvatarUrl:       "http://n0.xyz/avatar",
	}

	updateEvent := events.NewValidatorUpdateEvent(
		ctx,
		nr.Id,
		nr.ZetaPubKey,
		nr.ZetaPubKeyIndex,
		nr.EthereumAddress,
		nr.ChainPubKey,
		nr.InfoUrl,
		nr.Country,
		nr.Name,
		nr.AvatarUrl,
		nr.FromEpoch,
		true,
		0,
	)

	rankingEvent := events.NewValidatorRanking(ctx, "0", nr.Id, "0", "0", "0", "pending", "tendermint", 10)

	// one for the validator update, one for the ranking score
	broker.EXPECT().Send(updateEvent).Times(1)
	broker.EXPECT().Send(rankingEvent).Times(1)

	assert.NoError(t, top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint))
}

func TestValidatorTopologyKeyRotate(t *testing.T) {
	t.Run("add key rotate - success", testAddKeyRotateSuccess)
	t.Run("add key rotate - fails when node does not exists", testAddKeyRotateSuccessFailsOnNonExistingNode)
	t.Run("add key rotate - fails when target block height is less then current block height", testAddKeyRotateSuccessFailsWhenTargetBlockHeightIsLessThenCurrentBlockHeight)
	t.Run("add key rotate - fails when new key index is less then current current key index", testAddKeyRotateSuccessFailsWhenNewKeyIndexIsLessThenCurrentKeyIndex)
	t.Run("add key rotate - fails when key rotation for node already exists", testAddKeyRotateSuccessFailsWhenKeyRotationForNodeAlreadyExists)
	t.Run("add key rotate - fails when current pub key hash does not match", testAddKeyRotateSuccessFailsWhenCurrentPubKeyHashDoesNotMatch)
	t.Run("beginning of block - success", testBeginBlockSuccess)
	t.Run("beginning of block - notify key change", testBeginBlockNotifyKeyChange)
}

func testAddKeyRotateSuccess(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()
	top.timeService.EXPECT().GetTimeNow().AnyTimes()

	id := "zeta-master-pubkey"
	zetaPubKey := "zeta-key"
	newZetaPubKey := fmt.Sprintf("new-%s", zetaPubKey)

	nr := commandspb.AnnounceNode{
		Id:              id,
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      hex.EncodeToString([]byte(zetaPubKey)),
		EthereumAddress: "eth-address",
	}
	ctx := context.TODO()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	kr := &commandspb.KeyRotateSubmission{
		NewPubKeyIndex:    1,
		TargetBlock:       15,
		NewPubKey:         newZetaPubKey,
		CurrentPubKeyHash: hashKey(zetaPubKey),
	}

	err = top.AddKeyRotate(ctx, id, 10, kr)
	assert.NoError(t, err)
}

func testAddKeyRotateSuccessFailsOnNonExistingNode(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()
	top.timeService.EXPECT().GetTimeNow().AnyTimes()

	id := "zeta-master-pubkey"
	newZetaPubKey := "new-ega-key"

	ctx := context.TODO()

	err := top.AddKeyRotate(ctx, id, 10, newKeyRotationSubmission("", newZetaPubKey, 1, 10))
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to add key rotate for non existing node \"zeta-master-pubkey\"")
}

func testAddKeyRotateSuccessFailsWhenTargetBlockHeightIsLessThenCurrentBlockHeight(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()
	top.timeService.EXPECT().GetTimeNow().AnyTimes()

	id := "zeta-master-pubkey"
	zetaPubKey := "zeta-key"
	newZetaPubKey := fmt.Sprintf("new-%s", zetaPubKey)

	nr := commandspb.AnnounceNode{
		Id:              id,
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      hex.EncodeToString([]byte(zetaPubKey)),
		EthereumAddress: "eth-address",
	}
	ctx := context.TODO()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	err = top.AddKeyRotate(ctx, id, 15, newKeyRotationSubmission(zetaPubKey, newZetaPubKey, 1, 10))
	assert.ErrorIs(t, err, validators.ErrTargetBlockHeightMustBeGreaterThanCurrentHeight)
}

func testAddKeyRotateSuccessFailsWhenNewKeyIndexIsLessThenCurrentKeyIndex(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()
	top.timeService.EXPECT().GetTimeNow().AnyTimes()

	id := "zeta-master-pubkey"
	zetaPubKey := "zeta-key"
	newZetaPubKey := fmt.Sprintf("new-%s", zetaPubKey)

	nr := commandspb.AnnounceNode{
		Id:              id,
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      hex.EncodeToString([]byte(zetaPubKey)),
		EthereumAddress: "eth-address",
		ZetaPubKeyIndex: 2,
	}
	ctx := context.TODO()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	// test less then
	err = top.AddKeyRotate(ctx, id, 10, newKeyRotationSubmission(zetaPubKey, newZetaPubKey, 1, 15))
	assert.ErrorIs(t, err, validators.ErrNewZetaPubKeyIndexMustBeGreaterThenCurrentPubKeyIndex)
}

func testAddKeyRotateSuccessFailsWhenKeyRotationForNodeAlreadyExists(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()
	top.timeService.EXPECT().GetTimeNow().AnyTimes()

	id := "zeta-master-pubkey"
	zetaPubKey := "zeta-key"
	newZetaPubKey := fmt.Sprintf("new-%s", zetaPubKey)

	nr := commandspb.AnnounceNode{
		Id:              id,
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      hex.EncodeToString([]byte(zetaPubKey)),
		EthereumAddress: "eth-address",
		ZetaPubKeyIndex: 1,
	}
	ctx := context.TODO()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	// add first
	err = top.AddKeyRotate(ctx, id, 10, newKeyRotationSubmission(zetaPubKey, newZetaPubKey, 2, 12))
	assert.NoError(t, err)

	// add second
	err = top.AddKeyRotate(ctx, id, 10, newKeyRotationSubmission(zetaPubKey, newZetaPubKey, 2, 13))
	assert.ErrorIs(t, err, validators.ErrNodeAlreadyHasPendingKeyRotation)
}

func testAddKeyRotateSuccessFailsWhenCurrentPubKeyHashDoesNotMatch(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()
	top.timeService.EXPECT().GetTimeNow().AnyTimes()

	id := "zeta-master-pubkey"
	zetaPubKey := "zeta-key"
	newZetaPubKey := fmt.Sprintf("new-%s", zetaPubKey)

	nr := commandspb.AnnounceNode{
		Id:              id,
		ChainPubKey:     tmPubKey,
		ZetaPubKey:      hex.EncodeToString([]byte(zetaPubKey)),
		EthereumAddress: "eth-address",
		ZetaPubKeyIndex: 1,
	}
	ctx := context.TODO()
	err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
	assert.NoError(t, err)

	err = top.AddKeyRotate(ctx, id, 10, newKeyRotationSubmission("random-key", newZetaPubKey, 2, 12))
	assert.ErrorIs(t, err, validators.ErrCurrentPubKeyHashDoesNotMatch)
}

func hashKey(key string) string {
	return hex.EncodeToString(vgcrypto.Hash([]byte(key)))
}

func newKeyRotationSubmission(currentPubKey, newZetaPubKey string, keyIndex uint32, targetBlock uint64) *commandspb.KeyRotateSubmission {
	return &commandspb.KeyRotateSubmission{
		NewPubKeyIndex:    keyIndex,
		TargetBlock:       targetBlock,
		NewPubKey:         newZetaPubKey,
		CurrentPubKeyHash: hashKey(currentPubKey),
	}
}

func testBeginBlockSuccess(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()
	top.timeService.EXPECT().GetTimeNow().AnyTimes()

	chainValidators := []string{"tm-pubkey-1", "tm-pubkey-2", "tm-pubkey-3", "tm-pubkey-4"}

	ctx := context.TODO()
	for i := 0; i < len(chainValidators); i++ {
		j := i + 1
		id := fmt.Sprintf("zeta-master-pubkey-%d", j)
		nr := commandspb.AnnounceNode{
			Id:              id,
			ChainPubKey:     chainValidators[i],
			ZetaPubKey:      hex.EncodeToString([]byte(fmt.Sprintf("zeta-key-%d", j))),
			EthereumAddress: fmt.Sprintf("eth-address-%d", j),
		}

		err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
		assert.NoErrorf(t, err, "failed to add node registation %s", id)
	}

	// add key rotations
	err := top.AddKeyRotate(ctx, "zeta-master-pubkey-1", 10, newKeyRotationSubmission("zeta-key-1", "new-zeta-key-1", 1, 11))
	assert.NoError(t, err)
	err = top.AddKeyRotate(ctx, "zeta-master-pubkey-2", 10, newKeyRotationSubmission("zeta-key-2", "new-zeta-key-2", 1, 11))
	assert.NoError(t, err)
	err = top.AddKeyRotate(ctx, "zeta-master-pubkey-3", 10, newKeyRotationSubmission("zeta-key-3", "new-zeta-key-3", 1, 13))
	assert.NoError(t, err)
	err = top.AddKeyRotate(ctx, "zeta-master-pubkey-4", 10, newKeyRotationSubmission("zeta-key-4", "new-zeta-key-4", 1, 13))
	assert.NoError(t, err)

	// when
	top.BeginBlock(ctx, abcitypes.RequestBeginBlock{Header: types1.Header{Height: 11}})
	// then
	data1 := top.Get("zeta-master-pubkey-1")
	assert.NotNil(t, data1)
	assert.Equal(t, "new-zeta-key-1", data1.ZetaPubKey)
	data2 := top.Get("zeta-master-pubkey-2")
	assert.NotNil(t, data2)
	assert.Equal(t, "new-zeta-key-2", data2.ZetaPubKey)
	data3 := top.Get("zeta-master-pubkey-3")
	assert.NotNil(t, data3)
	assert.Equal(t, hex.EncodeToString([]byte("zeta-key-3")), data3.ZetaPubKey)
	data4 := top.Get("zeta-master-pubkey-4")
	assert.NotNil(t, data4)
	assert.Equal(t, hex.EncodeToString([]byte("zeta-key-4")), data4.ZetaPubKey)

	// when
	top.BeginBlock(ctx, abcitypes.RequestBeginBlock{Header: types1.Header{Height: 13}})
	// then
	data3 = top.Get("zeta-master-pubkey-3")
	assert.NotNil(t, data3)
	assert.Equal(t, "new-zeta-key-3", data3.ZetaPubKey)
	data4 = top.Get("zeta-master-pubkey-4")
	assert.NotNil(t, data4)
	assert.Equal(t, "new-zeta-key-4", data4.ZetaPubKey)
}

type Callback struct {
	mock.Mock
}

func (m *Callback) Call(ctx context.Context, a, b string) {
	m.Called(ctx, a, b)
}

func newCallback(times int) *Callback {
	c := Callback{}
	c.On("Call", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Times(times)
	return &c
}

func testBeginBlockNotifyKeyChange(t *testing.T) {
	top := getTestTopWithDefaultValidator(t)
	defer top.ctrl.Finish()
	top.timeService.EXPECT().GetTimeNow().AnyTimes()

	chainValidators := []string{"tm-pubkey-1", "tm-pubkey-2"}

	ctx := context.TODO()
	for i := 0; i < len(chainValidators); i++ {
		j := i + 1
		id := fmt.Sprintf("zeta-master-pubkey-%d", j)
		nr := commandspb.AnnounceNode{
			Id:              id,
			ChainPubKey:     chainValidators[i],
			ZetaPubKey:      hex.EncodeToString([]byte(fmt.Sprintf("zeta-key-%d", j))),
			EthereumAddress: fmt.Sprintf("eth-address-%d", j),
		}

		err := top.AddNewNode(ctx, &nr, validators.ValidatorStatusTendermint)
		assert.NoErrorf(t, err, "failed to add node registation %s", id)
	}

	// add key rotations
	err := top.AddKeyRotate(ctx, "zeta-master-pubkey-1", 10, newKeyRotationSubmission("zeta-key-1", "new-zeta-key-1", 1, 11))
	assert.NoError(t, err)
	err = top.AddKeyRotate(ctx, "zeta-master-pubkey-2", 10, newKeyRotationSubmission("zeta-key-2", "new-zeta-key-2", 1, 11))
	assert.NoError(t, err)

	// register callbacks
	c1 := newCallback(2)
	c2 := newCallback(2)
	top.NotifyOnKeyChange(c1.Call, c2.Call)

	// when
	top.BeginBlock(ctx, abcitypes.RequestBeginBlock{Header: types1.Header{Height: 11}})

	// then
	c1.AssertExpectations(t)
	c2.AssertExpectations(t)
}
