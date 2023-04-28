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
	"github.com/zeta-protocol/zeta/protos/zeta"
)

type DepositEvent interface {
	events.Event
	Deposit() zeta.Deposit
}

type DepositStore interface {
	Upsert(context.Context, *entities.Deposit) error
}

type Deposit struct {
	subscriber
	store DepositStore
}

func NewDeposit(store DepositStore) *Deposit {
	return &Deposit{
		store: store,
	}
}

func (d *Deposit) Types() []events.Type {
	return []events.Type{events.DepositEvent}
}

func (d *Deposit) Push(ctx context.Context, evt events.Event) error {
	return d.consume(ctx, evt.(DepositEvent))
}

func (d *Deposit) consume(ctx context.Context, event DepositEvent) error {
	deposit := event.Deposit()
	record, err := entities.DepositFromProto(&deposit, entities.TxHash(event.TxHash()), d.zetaTime)
	if err != nil {
		return errors.Wrap(err, "converting deposit proto to database entity failed")
	}

	return errors.Wrap(d.store.Upsert(ctx, record), "inserting deposit to SQL store failed")
}
