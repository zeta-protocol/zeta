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
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/zeta-protocol/zeta/cmd/zeta/commands/faucet"
	"github.com/zeta-protocol/zeta/cmd/zeta/commands/genesis"
	"github.com/zeta-protocol/zeta/cmd/zeta/commands/nodewallet"
	"github.com/zeta-protocol/zeta/cmd/zeta/commands/paths"
	tools "github.com/zeta-protocol/zeta/cmd/zetatools"
	"github.com/zeta-protocol/zeta/core/config"
)

// Subcommand is the signature of a sub command that can be registered.
type Subcommand func(context.Context, *flags.Parser) error

// Register registers one or more subcommands.
func Register(ctx context.Context, parser *flags.Parser, cmds ...Subcommand) error {
	for _, fn := range cmds {
		if err := fn(ctx, parser); err != nil {
			return err
		}
	}
	return nil
}

func Main(ctx context.Context) error {
	// special case for the tendermint subcommand, so we bypass the command line
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "tendermint", "tm":
			return (&tmCmd{}).Execute(nil)
		case "wallet":
			return (&walletCmd{}).Execute(nil)
		case "datanode":
			return (&datanodeCmd{}).Execute(nil)
		case "blockexplorer":
			return (&blockExplorerCmd{}).Execute(nil)
		}
	}

	parser := flags.NewParser(&config.Empty{}, flags.Default)

	if err := Register(ctx, parser,
		faucet.Faucet,
		genesis.Genesis,
		Init,
		nodewallet.NodeWallet,
		Verify,
		Version,
		Wallet,
		Datanode,
		tools.ZetaTools,
		Watch,
		Tm,
		Tendermint,
		Query,
		Bridge,
		paths.Paths,
		UnsafeResetAll,
		AnnounceNode,
		RotateEthKey,
		ProposeProtocolUpgrade,
		Start,
		Node,
		BlockExplorer,
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return err
	}

	if _, err := parser.Parse(); err != nil {
		return err
	}
	return nil
}
