package entities

import (
	"encoding/hex"
	"time"

	"github.com/zeta-protocol/zeta/core/events"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"
	"github.com/pkg/errors"
)

type BeginBlockEvent interface {
	events.Event
	BeginBlock() eventspb.BeginBlock
}

func BlockFromBeginBlock(b BeginBlockEvent) (*Block, error) {
	hash, err := hex.DecodeString(b.TraceID())
	if err != nil {
		return nil, errors.Wrapf(err, "Trace ID is not valid hex string, trace ID:%s", b.TraceID())
	}

	zetaTime := time.Unix(0, b.BeginBlock().Timestamp)

	// Postgres only stores timestamps in microsecond resolution
	block := Block{
		ZetaTime: zetaTime.Truncate(time.Microsecond),
		Hash:     hash,
		Height:   b.BlockNr(),
	}
	return &block, err
}
