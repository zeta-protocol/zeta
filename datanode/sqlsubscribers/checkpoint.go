// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package sqlsubscribers

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeta-protocol/zeta/core/events"
	"github.com/zeta-protocol/zeta/datanode/entities"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"
)

type CheckpointEvent interface {
	events.Event
	Proto() eventspb.CheckpointEvent
}

type CheckpointStore interface {
	Add(context.Context, entities.Checkpoint) error
}

type Checkpoint struct {
	subscriber
	store CheckpointStore
}

func NewCheckpoint(store CheckpointStore) *Checkpoint {
	np := &Checkpoint{
		store: store,
	}
	return np
}

func (n *Checkpoint) Types() []events.Type {
	return []events.Type{events.CheckpointEvent}
}

func (n *Checkpoint) Push(ctx context.Context, evt events.Event) error {
	return n.consume(ctx, evt.(CheckpointEvent))
}

func (n *Checkpoint) consume(ctx context.Context, event CheckpointEvent) error {
	pnp := event.Proto()
	np, err := entities.CheckpointFromProto(&pnp, entities.TxHash(event.TxHash()))
	if err != nil {
		return errors.Wrap(err, "unable to parse checkpoint")
	}
	np.ZetaTime = n.zetaTime
	np.SeqNum = event.Sequence()

	if err := n.store.Add(ctx, np); err != nil {
		return errors.Wrap(err, "error adding checkpoint")
	}

	return nil
}
