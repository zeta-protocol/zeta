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

package nodewallets

import (
	"errors"

	"github.com/zeta-protocol/zeta/core/nodewallets/eth"
	"github.com/zeta-protocol/zeta/core/nodewallets/zeta"
)

var (
	ErrZetaWalletIsMissing       = errors.New("the Zeta node wallet is missing")
	ErrEthereumWalletIsMissing   = errors.New("the Ethereum node wallet is missing")
	ErrTendermintPubkeyIsMissing = errors.New("the Tendermint pubkey is missing")
)

type TendermintPubkey struct {
	Pubkey string
}

type NodeWallets struct {
	Zeta       *zeta.Wallet
	Ethereum   *eth.Wallet
	Tendermint *TendermintPubkey
}

func (w *NodeWallets) SetEthereumWallet(ethWallet *eth.Wallet) {
	w.Ethereum = ethWallet
}

func (w *NodeWallets) Verify() error {
	if w.Zeta == nil {
		return ErrZetaWalletIsMissing
	}
	if w.Ethereum == nil {
		return ErrEthereumWalletIsMissing
	}
	if w.Tendermint == nil {
		return ErrTendermintPubkeyIsMissing
	}
	return nil
}
