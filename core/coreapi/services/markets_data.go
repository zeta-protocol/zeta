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

type marketDataE interface {
	events.Event
	MarketData() zetapb.MarketData
}

type MarketsData struct {
	*subscribers.Base
	ctx context.Context

	mu          sync.RWMutex
	marketsData map[string]zetapb.MarketData
	ch          chan zetapb.MarketData
}

func NewMarketsData(ctx context.Context) (marketsData *MarketsData) {
	defer func() { go marketsData.consume() }()
	return &MarketsData{
		Base:        subscribers.NewBase(ctx, 1000, true),
		ctx:         ctx,
		marketsData: map[string]zetapb.MarketData{},
		ch:          make(chan zetapb.MarketData, 100),
	}
}

func (m *MarketsData) consume() {
	defer func() { close(m.ch) }()
	for {
		select {
		case <-m.Closed():
			return
		case marketData, ok := <-m.ch:
			if !ok {
				// cleanup base
				m.Halt()
				// channel is closed
				return
			}
			m.mu.Lock()
			m.marketsData[marketData.Market] = marketData
			m.mu.Unlock()
		}
	}
}

func (m *MarketsData) Push(evts ...events.Event) {
	for _, e := range evts {
		if ae, ok := e.(marketDataE); ok {
			m.ch <- ae.MarketData()
		}
	}
}

func (m *MarketsData) List(marketID string) []*zetapb.MarketData {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(marketID) > 0 {
		return m.getMarketData(marketID)
	}
	return m.getAllMarketsData()
}

func (m *MarketsData) getMarketData(marketID string) []*zetapb.MarketData {
	out := []*zetapb.MarketData{}
	asset, ok := m.marketsData[marketID]
	if ok {
		out = append(out, &asset)
	}
	return out
}

func (m *MarketsData) getAllMarketsData() []*zetapb.MarketData {
	out := make([]*zetapb.MarketData, 0, len(m.marketsData))
	for _, v := range m.marketsData {
		v := v
		out = append(out, &v)
	}
	return out
}

func (m *MarketsData) Types() []events.Type {
	return []events.Type{
		events.MarketDataEvent,
	}
}
