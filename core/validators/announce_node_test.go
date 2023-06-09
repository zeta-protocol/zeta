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
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/zeta-protocol/zeta/core/nodewallets"
	"github.com/zeta-protocol/zeta/core/validators"
	vgrand "github.com/zeta-protocol/zeta/libs/rand"
	vgtesting "github.com/zeta-protocol/zeta/libs/testing"
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestTendermintKey(t *testing.T) {
	notBase64 := "170ffakjde"
	require.Error(t, validators.VerifyTendermintKey(notBase64))

	validKey := "794AFpbqJvHF711mhAK3fvSLnoXuuiig2ecrdeSJ/bk="
	require.NoError(t, validators.VerifyTendermintKey(validKey))
}

func TestSignVerifyAnnounceNode(t *testing.T) {
	cmd := createSignedAnnounceCommand(t)
	require.NoError(t, validators.VerifyAnnounceNode(cmd))
}

func TestDoubleAnnounce(t *testing.T) {
	tt := getTestTopology(t)
	cmd := createSignedAnnounceCommand(t)
	ctx := context.Background()

	// Add it once
	require.NoError(t, tt.Topology.ProcessAnnounceNode(ctx, cmd))

	// Add it again
	require.ErrorIs(t, tt.Topology.ProcessAnnounceNode(ctx, cmd), validators.ErrZetaNodeAlreadyRegisterForChain)
}

func createSignedAnnounceCommand(t *testing.T) *commandspb.AnnounceNode {
	t.Helper()
	nodeWallets := createTestNodeWallets(t)
	cmd := commandspb.AnnounceNode{
		Id:              nodeWallets.Zeta.ID().Hex(),
		ZetaPubKey:      nodeWallets.Zeta.PubKey().Hex(),
		ZetaPubKeyIndex: nodeWallets.Zeta.Index(),
		ChainPubKey:     "794AFpbqJvHF711mhAK3fvSLnoXuuiig2ecrdeSJ/bk=",
		EthereumAddress: nodeWallets.Ethereum.PubKey().Hex(),
		FromEpoch:       1,
		InfoUrl:         "www.some.com",
		Name:            "that is not my name",
		AvatarUrl:       "www.avatar.com",
		Country:         "some country",
	}
	err := validators.SignAnnounceNode(&cmd, nodeWallets.Zeta, nodeWallets.Ethereum)
	require.NoError(t, err)

	// verify that the expected signature for zeta key is there
	messageToSign := cmd.Id + cmd.ZetaPubKey + fmt.Sprintf("%d", cmd.ZetaPubKeyIndex) + cmd.ChainPubKey + cmd.EthereumAddress + fmt.Sprintf("%d", cmd.FromEpoch) + cmd.InfoUrl + cmd.Name + cmd.AvatarUrl + cmd.Country
	sig, err := nodeWallets.Zeta.Sign([]byte(messageToSign))
	sigHex := hex.EncodeToString(sig)
	require.NoError(t, err)
	require.Equal(t, sigHex, cmd.ZetaSignature.Value)

	// verify that the expected signature for eth key is there
	ethSig, err := nodeWallets.Ethereum.Sign(crypto.Keccak256([]byte(messageToSign)))
	ethSigHex := hex.EncodeToString(ethSig)
	require.NoError(t, err)
	require.Equal(t, ethSigHex, cmd.EthereumSignature.Value)

	return &cmd
}

func createTestNodeWallets(t *testing.T) *nodewallets.NodeWallets {
	t.Helper()
	config := nodewallets.NewDefaultConfig()
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	if _, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletsPass, "", false); err != nil {
		panic("couldn't generate Ethereum node wallet for tests")
	}

	if _, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletsPass, false); err != nil {
		panic("couldn't generate Zeta node wallet for tests")
	}
	nw, err := nodewallets.GetNodeWallets(config, zetaPaths, registryPass)
	require.NoError(t, err)
	return nw
}
