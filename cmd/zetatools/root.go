package tools

import (
	"context"

	"github.com/zeta-protocol/zeta/core/config"

	"github.com/jessevdk/go-flags"
)

type RootCmd struct {
	// Global options
	config.ZetaHomeFlag

	// Subcommands
	Snapshot   snapshotCmd   `command:"snapshot" description:"Display information about saved snapshots"`
	Checkpoint checkpointCmd `command:"checkpoint" description:"Make checkpoint human-readable, or generate checkpoint from human readable format"`
	Stream     streamCmd     `command:"stream" description:"Stream events from zeta node"`
}

var rootCmd RootCmd

func ZetaTools(ctx context.Context, parser *flags.Parser) error {
	rootCmd = RootCmd{
		Snapshot:   snapshotCmd{},
		Checkpoint: checkpointCmd{},
		Stream:     streamCmd{},
	}

	var (
		short = "useful tooling for probing a zeta node and its state"
		long  = `useful tooling for probing a zeta node and its state`
	)
	_, err := parser.AddCommand("tools", short, long, &rootCmd)
	return err
}
