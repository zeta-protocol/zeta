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
	"runtime/debug"

	"github.com/zeta-protocol/zeta/libs/memory"

	"github.com/zeta-protocol/zeta/cmd/data-node/commands/start"
	"github.com/zeta-protocol/zeta/datanode/config"
	"github.com/zeta-protocol/zeta/logging"
	"github.com/zeta-protocol/zeta/paths"
	"github.com/zeta-protocol/zeta/version"

	"github.com/jessevdk/go-flags"
)

type StartCmd struct {
	config.ZetaHomeFlag

	config.Config
}

var startCmd StartCmd

const namedLogger = "datanode"

func (cmd *StartCmd) Execute(args []string) error {
	ctx, cfunc := context.WithCancel(context.Background())
	defer cfunc()
	log := logging.NewLoggerFromConfig(
		logging.NewDefaultConfig()).Named(namedLogger)
	defer log.AtExit()

	// we define this option to parse the cli args each time the config is
	// loaded. So that we can respect the cli flag precedence.
	parseFlagOpt := func(cfg *config.Config) error {
		_, err := flags.NewParser(cfg, flags.Default|flags.IgnoreUnknown).Parse()
		return err
	}

	zetaPaths := paths.New(cmd.ZetaHome)

	configWatcher, err := config.NewWatcher(context.Background(), log, zetaPaths, config.Use(parseFlagOpt))
	if err != nil {
		return err
	}

	// setup max memory usage
	memFactor, err := configWatcher.Get().GetMaxMemoryFactor()
	if err != nil {
		return err
	}

	// only set max memory if user didn't require 100%
	if memFactor != 1 {
		totalMem, err := memory.TotalMemory()
		if err != nil {
			return fmt.Errorf("failed to get total memory: %w", err)
		}
		debug.SetMemoryLimit(int64(float64(totalMem) * memFactor))
	}

	return (&start.NodeCommand{
		Log:         log,
		Version:     version.Get(),
		VersionHash: version.GetCommitHash(),
	}).Run(
		ctx,
		configWatcher,
		zetaPaths,
		args,
	)
}

func Node(ctx context.Context, parser *flags.Parser) error {
	startCmd = StartCmd{
		Config: config.NewDefaultConfig(),
	}
	cmd, err := parser.AddCommand("node", "deprecated, see data-node start instead", "deprecated, see data-node start instead", &startCmd)
	if err != nil {
		return err
	}

	// Print nested groups under parent's name using `::` as the separator.
	for _, parent := range cmd.Groups() {
		for _, grp := range parent.Groups() {
			grp.ShortDescription = parent.ShortDescription + "::" + grp.ShortDescription
		}
	}
	return nil
}

func Start(_ context.Context, parser *flags.Parser) error {
	startCmd = StartCmd{
		Config: config.NewDefaultConfig(),
	}
	cmd, err := parser.AddCommand("start", "Start a zeta data node", "Start a zeta data node as defined by the config files", &startCmd)
	if err != nil {
		return err
	}

	// Print nested groups under parent's name using `::` as the separator.
	for _, parent := range cmd.Groups() {
		for _, grp := range parent.Groups() {
			grp.ShortDescription = parent.ShortDescription + "::" + grp.ShortDescription
		}
	}
	return nil
}
