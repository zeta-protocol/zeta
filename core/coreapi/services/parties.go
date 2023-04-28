// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.ZETA file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package services

import (
	"context"
	"sync"

	"github.com/zeta-protocol/zeta/core/events"
	"github.com/zeta-protocol/zeta/core/subscribers"
	"github.com/zeta-protocol/zeta/core/types"
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
)

type partyE interface {
	events.Event
	Party() zetapb.Party
}

type Parties struct {
	*subscribers.Base
	ctx context.Context

	mu      sync.RWMutex
	parties map[string]zetapb.Party
	ch      chan zetapb.Party
}

func NewParties(ctx context.Context) (parties *Parties) {
	defer func() {
		parties.parties[types.NetworkParty] = zetapb.Party{Id: types.NetworkParty}
		go parties.consume()
	}()
	return &Parties{
		Base:    subscribers.NewBase(ctx, 1000, true),
		ctx:     ctx,
		parties: map[string]zetapb.Party{},
		ch:      make(chan zetapb.Party, 100),
	}
}

func (a *Parties) consume() {
	defer func() { close(a.ch) }()
	for {
		select {
		case <-a.Closed():
			return
		case party, ok := <-a.ch:
			if !ok {
				// cleanup base
				a.Halt()
				// channel is closed
				return
			}
			a.mu.Lock()
			a.parties[party.Id] = party
			a.mu.Unlock()
		}
	}
}

func (a *Parties) Push(evts ...events.Event) {
	for _, e := range evts {
		if ae, ok := e.(partyE); ok {
			a.ch <- ae.Party()
		}
	}
}

func (a *Parties) List() []*zetapb.Party {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]*zetapb.Party, 0, len(a.parties))
	for _, v := range a.parties {
		v := v
		out = append(out, &v)
	}
	return out
}

func (a *Parties) Types() []events.Type {
	return []events.Type{
		events.PartyEvent,
	}
}
