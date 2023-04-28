package tools

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zeta-protocol/zeta/core/config"
	"github.com/zeta-protocol/zeta/paths"
	"github.com/zeta-protocol/zeta/zetatools/snapshotdb"
)

type snapshotCmd struct {
	config.OutputFlag
	DBPath               string `short:"d" long:"db-path" description:"path to snapshot state data"`
	SnapshotContentsPath string `short:"c" long:"snapshot-contents" description:"path to file where to write the content of a snapshot"`
	BlockHeight          uint64 `short:"b" long:"block-height" description:"block-height of requested snapshot"`
}

func (opts *snapshotCmd) Execute(_ []string) error {
	if opts.SnapshotContentsPath != "" && opts.BlockHeight == 0 {
		return errors.New("must specify --block-height when using --write-payload")
	}

	db := opts.DBPath
	if opts.DBPath == "" {
		zetaPaths := paths.New(rootCmd.ZetaHome)
		db = zetaPaths.StatePathFor(paths.SnapshotStateHome)
	}

	if opts.SnapshotContentsPath != "" {
		fmt.Printf("finding payloads for block-height %d...\n", opts.BlockHeight)
		err := snapshotdb.SavePayloadsToFile(db, opts.SnapshotContentsPath, opts.BlockHeight)
		if err != nil {
			return err
		}
		fmt.Printf("payloads saved to '%s'\n", opts.SnapshotContentsPath)
		return nil
	}

	snapshots, invalid, err := snapshotdb.SnapshotData(db, opts.SnapshotContentsPath, opts.BlockHeight)
	if err != nil {
		return err
	}
	if opts.Output.IsJSON() {
		o := struct {
			Snapshots []snapshotdb.Data `json:"snapshots"`
			Invalid   []snapshotdb.Data `json:"invalidSnapshots,omitempty"`
		}{
			Snapshots: snapshots,
			Invalid:   invalid,
		}
		b, err := json.Marshal(o)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}

	fmt.Println("Snapshots available:", len(snapshots))
	for _, snap := range snapshots {
		fmt.Printf("\tHeight: %d, Version: %d, Size %d, Hash: %s\n", snap.Height, snap.Version, snap.Size, snap.Hash)
	}

	if len(invalid) == 0 {
		return nil
	}

	fmt.Println("Invalid snapshots:", len(invalid))
	for _, snap := range invalid {
		fmt.Printf("\tVersion: %d, Size %d, Hash: %s\n", snap.Version, snap.Size, snap.Hash)
	}
	return nil
}
