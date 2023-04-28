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

type proposalE interface {
	events.Event
	Proposal() zetapb.Proposal
}

type Proposals struct {
	*subscribers.Base
	ctx context.Context

	mu        sync.RWMutex
	proposals map[string]zetapb.Proposal
	// map of proposer -> set of proposal id
	proposalsPerProposer map[string]map[string]struct{}
	ch                   chan zetapb.Proposal
}

func NewProposals(ctx context.Context) (proposals *Proposals) {
	defer func() { go proposals.consume() }()
	return &Proposals{
		Base:                 subscribers.NewBase(ctx, 1000, true),
		ctx:                  ctx,
		proposals:            map[string]zetapb.Proposal{},
		proposalsPerProposer: map[string]map[string]struct{}{},
		ch:                   make(chan zetapb.Proposal, 100),
	}
}

func (p *Proposals) consume() {
	defer func() { close(p.ch) }()
	for {
		select {
		case <-p.Closed():
			return
		case prop, ok := <-p.ch:
			if !ok {
				// cleanup base
				p.Halt()
				// channel is closed
				return
			}
			p.mu.Lock()
			p.proposals[prop.Id] = prop
			proposals, ok := p.proposalsPerProposer[prop.PartyId]
			if !ok {
				proposals = map[string]struct{}{}
				p.proposalsPerProposer[prop.PartyId] = proposals
			}
			proposals[prop.Id] = struct{}{}
			p.mu.Unlock()
		}
	}
}

func (p *Proposals) Push(evts ...events.Event) {
	for _, e := range evts {
		if ae, ok := e.(proposalE); ok {
			p.ch <- ae.Proposal()
		}
	}
}

func (p *Proposals) List(proposal, party string) []*zetapb.Proposal {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(proposal) <= 0 && len(party) <= 0 {
		return p.getAllProposals()
	} else if len(party) > 0 {
		return p.getProposalsPerParty(proposal, party)
	} else if len(proposal) > 0 {
		return p.getProposalByID(proposal)
	}
	return p.getAllProposals()
}

func (p *Proposals) getProposalsPerParty(proposal, party string) []*zetapb.Proposal {
	out := []*zetapb.Proposal{}
	partyProposals, ok := p.proposalsPerProposer[party]
	if !ok {
		return out
	}

	if len(proposal) > 0 {
		_, ok := partyProposals[proposal]
		if ok {
			prop := p.proposals[proposal]
			out = append(out, &prop)
		}
		return out
	}

	for k := range partyProposals {
		prop := p.proposals[k]
		out = append(out, &prop)
	}
	return out
}

func (p *Proposals) getProposalByID(proposal string) []*zetapb.Proposal {
	out := []*zetapb.Proposal{}
	asset, ok := p.proposals[proposal]
	if ok {
		out = append(out, &asset)
	}
	return out
}

func (p *Proposals) getAllProposals() []*zetapb.Proposal {
	out := make([]*zetapb.Proposal, 0, len(p.proposals))
	for _, v := range p.proposals {
		v := v
		out = append(out, &v)
	}
	return out
}

func (p *Proposals) Types() []events.Type {
	return []events.Type{
		events.ProposalEvent,
	}
}
