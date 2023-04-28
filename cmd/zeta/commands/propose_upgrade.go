// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeta-protocol/zeta/core/config"
	"github.com/zeta-protocol/zeta/core/protocolupgrade"
	"github.com/zeta-protocol/zeta/core/txn"
	vgcrypto "github.com/zeta-protocol/zeta/libs/crypto"
	vgjson "github.com/zeta-protocol/zeta/libs/json"
	"github.com/zeta-protocol/zeta/logging"
	"github.com/zeta-protocol/zeta/paths"
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"

	"github.com/blang/semver"
	"github.com/jessevdk/go-flags"
)

type ProposeUpgradeCmd struct {
	config.ZetaHomeFlag
	config.OutputFlag
	config.Passphrase `long:"passphrase-file"`

	ZetaReleaseTag     string `short:"v" long:"zeta-release-tag" required:"true" description:"A valid zeta core release tag for the upgrade proposal"`
	UpgradeBlockHeight uint64 `short:"h" long:"height" required:"true" description:"The block height at which the upgrade should be made"`
}

var proposeUpgradeCmd ProposeUpgradeCmd

func (opts *ProposeUpgradeCmd) Execute(_ []string) error {
	log := logging.NewLoggerFromConfig(logging.NewDefaultConfig())
	defer log.AtExit()

	output, err := opts.GetOutput()
	if err != nil {
		return err
	}

	registryPass, err := opts.Get("node wallet", false)
	if err != nil {
		return err
	}

	zetaPaths := paths.New(opts.ZetaHome)

	_, conf, err := config.EnsureNodeConfig(zetaPaths)
	if err != nil {
		return err
	}

	if !conf.IsValidator() {
		return errors.New("node is not a validator")
	}

	cmd := commandspb.ProtocolUpgradeProposal{
		ZetaReleaseTag:     opts.ZetaReleaseTag,
		UpgradeBlockHeight: opts.UpgradeBlockHeight,
	}

	commander, blockData, cfunc, err := getNodeWalletCommander(log, registryPass, zetaPaths)
	if err != nil {
		return fmt.Errorf("failed to get commander: %w", err)
	}

	if opts.UpgradeBlockHeight <= blockData.Height {
		return fmt.Errorf("upgrade block earlier than current block height")
	}

	_, err = semver.Parse(protocolupgrade.TrimReleaseTag(opts.ZetaReleaseTag))
	if err != nil {
		return fmt.Errorf("invalid protocol version for upgrade received: version (%s), %w", opts.ZetaReleaseTag, err)
	}

	defer cfunc()

	tid := vgcrypto.RandomHash()
	powNonce, _, err := vgcrypto.PoW(blockData.Hash, tid, uint(blockData.SpamPowDifficulty), vgcrypto.Sha3)
	if err != nil {
		return fmt.Errorf("failed to get proof of work: %w", err)
	}

	pow := &commandspb.ProofOfWork{
		Tid:   tid,
		Nonce: powNonce,
	}

	var txHash string
	ch := make(chan struct{})
	commander.CommandWithPoW(
		context.Background(),
		txn.ProtocolUpgradeCommand,
		&cmd,
		func(h string, e error) {
			txHash, err = h, e
			close(ch)
		}, nil, pow)

	<-ch

	if err != nil {
		return err
	}

	if output.IsHuman() {
		fmt.Printf("Upgrade proposal sent.\ntxHash: %s", txHash)
	} else if output.IsJSON() {
		return vgjson.Print(struct {
			TxHash string `json:"txHash"`
		}{
			TxHash: txHash,
		})
	}
	return err
}

func ProposeProtocolUpgrade(ctx context.Context, parser *flags.Parser) error {
	proposeUpgradeCmd = ProposeUpgradeCmd{}

	var (
		short = "Propose a protocol upgrade"
		long  = "Propose a protocol upgrade"
	)
	_, err := parser.AddCommand("protocol_upgrade_proposal", short, long, &proposeUpgradeCmd)
	return err
}
