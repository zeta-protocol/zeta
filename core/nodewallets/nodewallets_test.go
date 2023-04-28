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

package nodewallets_test

import (
	"testing"

	"github.com/zeta-protocol/zeta/core/nodewallets"
	ethnw "github.com/zeta-protocol/zeta/core/nodewallets/eth"
	vgnw "github.com/zeta-protocol/zeta/core/nodewallets/zeta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeWallet(t *testing.T) {
	t.Run("Verify node wallets succeeds", testVerifyNodeWalletsSucceeds)
	t.Run("Verify node wallets fails", testVerifyNodeWalletsFails)
}

func testVerifyNodeWalletsSucceeds(t *testing.T) {
	nw := &nodewallets.NodeWallets{
		Zeta:       &vgnw.Wallet{},
		Ethereum:   &ethnw.Wallet{},
		Tendermint: &nodewallets.TendermintPubkey{},
	}

	assert.NoError(t, nw.Verify())
}

func testVerifyNodeWalletsFails(t *testing.T) {
	tcs := []struct {
		name        string
		expectedErr error
		nw          *nodewallets.NodeWallets
	}{
		{
			name:        "with missing Ethereum wallet",
			expectedErr: nodewallets.ErrEthereumWalletIsMissing,
			nw: &nodewallets.NodeWallets{
				Zeta: &vgnw.Wallet{},
			},
		}, {
			name:        "with missing Zeta wallet",
			expectedErr: nodewallets.ErrZetaWalletIsMissing,
			nw: &nodewallets.NodeWallets{
				Ethereum: &ethnw.Wallet{},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			require.EqualError(tt, tc.nw.Verify(), tc.expectedErr.Error())
		})
	}
}
