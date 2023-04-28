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
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
)

type marketE interface {
	events.Event
	Market() zetapb.Market
}

type Markets struct {
	*subscribers.Base
	ctx context.Context

	mu      sync.RWMutex
	markets map[string]zetapb.Market
	ch      chan zetapb.Market
}

func NewMarkets(ctx context.Context) (markets *Markets) {
	defer func() { go markets.consume() }()
	return &Markets{
		Base:    subscribers.NewBase(ctx, 1000, true),
		ctx:     ctx,
		markets: map[string]zetapb.Market{},
		ch:      make(chan zetapb.Market, 100),
	}
}

func (m *Markets) consume() {
	defer func() { close(m.ch) }()
	for {
		select {
		case <-m.Closed():
			return
		case market, ok := <-m.ch:
			if !ok {
				// cleanup base
				m.Halt()
				// channel is closed
				return
			}
			m.mu.Lock()
			m.markets[market.Id] = market
			m.mu.Unlock()
		}
	}
}

func (m *Markets) Push(evts ...events.Event) {
	for _, e := range evts {
		if ae, ok := e.(marketE); ok {
			m.ch <- ae.Market()
		}
	}
}

func (m *Markets) List(marketID string) []*zetapb.Market {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(marketID) > 0 {
		return m.getMarket(marketID)
	}
	return m.getAllMarkets()
}

func (m *Markets) getMarket(marketID string) []*zetapb.Market {
	out := []*zetapb.Market{}
	asset, ok := m.markets[marketID]
	if ok {
		out = append(out, &asset)
	}
	return out
}

func (m *Markets) getAllMarkets() []*zetapb.Market {
	out := make([]*zetapb.Market, 0, len(m.markets))
	for _, v := range m.markets {
		v := v
		out = append(out, &v)
	}
	return out
}

func (m *Markets) Types() []events.Type {
	return []events.Type{
		events.MarketCreatedEvent,
		events.MarketUpdatedEvent,
	}
}
