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
	coreapipb "github.com/zeta-protocol/zeta/protos/zeta/api/v1"
)

type accountE interface {
	events.Event
	Account() zetapb.Account
}

type Accounts struct {
	*subscribers.Base

	mu sync.RWMutex
	// parties -> accounts id -> accounts
	parties map[string]map[string]zetapb.Account
	// markets id -> accounts id -> account
	markets map[string]map[string]zetapb.Account
	// global accounts id -> account
	globals map[string]zetapb.Account
}

func NewAccounts(ctx context.Context) *Accounts {
	return &Accounts{
		Base:    subscribers.NewBase(ctx, 1000, true),
		parties: map[string]map[string]zetapb.Account{},
		markets: map[string]map[string]zetapb.Account{},
		globals: map[string]zetapb.Account{},
	}
}

func (a *Accounts) Push(evts ...events.Event) {
	if len(evts) == 0 {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, e := range evts {
		switch acc := e.(type) {
		case accountE:
			a.addAccount(acc.Account())
		}
	}
}

func (a *Accounts) List(party, market string) []*coreapipb.Account {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(party) > 0 {
		return a.getPartyAccounts(party, market)
	}
	if len(market) > 0 {
		return a.getMarketAccounts(market)
	}
	return a.getGlobalAccounts()
}

func (a *Accounts) Types() []events.Type {
	return []events.Type{
		events.AccountEvent,
	}
}

func (a *Accounts) getPartyAccounts(party, market string) []*coreapipb.Account {
	accs, ok := a.parties[party]
	if !ok {
		return nil
	}

	// at least one
	out := make([]*coreapipb.Account, 0, 1)
	for _, v := range accs {
		if len(market) > 0 && v.MarketId != market {
			continue
		}
		out = append(out, toAccount(v))
	}

	return out
}

func (a *Accounts) getMarketAccounts(market string) []*coreapipb.Account {
	accs, ok := a.markets[market]
	if !ok {
		return nil
	}

	out := make([]*coreapipb.Account, 0, len(accs))
	for _, v := range accs {
		out = append(out, toAccount(v))
	}

	return out
}

func (a *Accounts) getGlobalAccounts() []*coreapipb.Account {
	out := make([]*coreapipb.Account, 0, len(a.globals))
	for _, v := range a.globals {
		out = append(out, toAccount(v))
	}

	return out
}

func (a *Accounts) addAccount(acc zetapb.Account) {
	if acc.MarketId == "!" && acc.Owner == "*" {
		a.globals[acc.Id] = acc
	}

	if acc.Owner != "*" {
		a.addPartyAccount(acc)
	}

	a.addMarketAccount(acc)
}

func (a *Accounts) addPartyAccount(acc zetapb.Account) {
	accs, ok := a.parties[acc.Owner]
	if !ok {
		accs = map[string]zetapb.Account{}
		a.parties[acc.Owner] = accs
	}
	accs[acc.Id] = acc
}

func (a *Accounts) addMarketAccount(acc zetapb.Account) {
	accs, ok := a.parties[acc.MarketId]
	if !ok {
		accs = map[string]zetapb.Account{}
		a.parties[acc.MarketId] = accs
	}
	accs[acc.Id] = acc
}

func toAccount(acc zetapb.Account) *coreapipb.Account {
	market := ""
	if acc.MarketId != "!" {
		market = acc.MarketId
	}
	owner := "0000000000000000000000000000000000000000000000000000000000000000"
	if acc.Owner != "*" {
		owner = acc.Owner
	}

	return &coreapipb.Account{
		Party:   owner,
		Market:  market,
		Balance: acc.Balance,
		Asset:   acc.Asset,
		Type:    acc.Type.String(),
	}
}
