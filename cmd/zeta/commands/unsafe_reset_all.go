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
	"os"
	"path/filepath"

	"github.com/zeta-protocol/zeta/core/config"
	"github.com/zeta-protocol/zeta/logging"
	"github.com/zeta-protocol/zeta/paths"

	"github.com/jessevdk/go-flags"
)

type UnsafeResetAllCmd struct {
	config.ZetaHomeFlag
}

//nolint:unparam
func (cmd *UnsafeResetAllCmd) Execute(_ []string) error {
	log := logging.NewLoggerFromConfig(
		logging.NewDefaultConfig(),
	)
	defer log.AtExit()

	zetaPaths := paths.New(cmd.ZetaHome)

	snapshotsPath := zetaPaths.StatePathFor(paths.SnapshotStateHome)
	if err := deleteAll(log, snapshotsPath); err != nil {
		log.Error("Unable to remove snapshot state", logging.Error(err))
	} else {
		log.Info("Removed snapshot state", logging.String("path", snapshotsPath))
	}

	checkpointsPath := zetaPaths.StatePathFor(paths.CheckpointStateHome)
	if err := deleteAll(log, checkpointsPath); err != nil {
		log.Error("Unable to remove checkpoint state", logging.Error(err))
	} else {
		log.Info("Removed checkpoint state", logging.String("path", checkpointsPath))
	}
	return nil
}

func deleteAll(log *logging.Logger, dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	files, err := d.Readdir(0)
	if err != nil {
		return err
	}

	for _, f := range files {
		filePath := filepath.Join(dir, f.Name())
		if err := os.RemoveAll(filePath); err != nil {
			return err
		}
		log.Info("Removed file", logging.String("path", filePath))
	}

	return nil
}

var unsafeResetCmd UnsafeResetAllCmd

func UnsafeResetAll(ctx context.Context, parser *flags.Parser) error {
	unsafeResetCmd = UnsafeResetAllCmd{}

	_, err := parser.AddCommand("unsafe_reset_all", "(unsafe) Remove all application state", "(unsafe) Remove all zeta application state (checkpoints and snapshots)", &unsafeResetCmd)
	return err
}
