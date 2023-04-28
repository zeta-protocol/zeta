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

	"github.com/zeta-protocol/zeta/cmd/data-node/commands/start"
	"github.com/zeta-protocol/zeta/datanode/config"
	"github.com/zeta-protocol/zeta/logging"
	"github.com/zeta-protocol/zeta/paths"

	"github.com/jessevdk/go-flags"
)

type UnsafeResetAllCmd struct {
	config.ZetaHomeFlag
	*config.Config
}

//nolint:unparam
func (cmd *UnsafeResetAllCmd) Execute(_ []string) error {
	ctx, cfunc := context.WithCancel(context.Background())
	defer cfunc()
	log := logging.NewLoggerFromConfig(
		logging.NewDefaultConfig(),
	)
	defer log.AtExit()

	zetaPaths := paths.New(cmd.ZetaHome)

	cfgLoader, err := config.InitialiseLoader(zetaPaths)
	if err != nil {
		return fmt.Errorf("couldn't initialise configuration loader: %w", err)
	}

	cmd.Config, err = cfgLoader.Get()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	connConfig := cmd.Config.SQLStore.ConnectionConfig
	if start.ResetDatabaseAndNetworkHistory(ctx, log, zetaPaths, connConfig); err != nil {
		return fmt.Errorf("failed to reset database and network history: %w", err)
	}

	return nil
}

var unsafeResetCmd UnsafeResetAllCmd

func UnsafeResetAll(_ context.Context, parser *flags.Parser) error {
	unsafeResetCmd = UnsafeResetAllCmd{}

	_, err := parser.AddCommand("unsafe_reset_all", "(unsafe) Remove all application state", "(unsafe) Remove all datanode application state (wipes the database and removes all network history)", &unsafeResetCmd)
	return err
}
