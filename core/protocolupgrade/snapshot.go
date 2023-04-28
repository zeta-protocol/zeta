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

package protocolupgrade

import (
	"context"
	"sort"

	"github.com/zeta-protocol/zeta/core/types"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"
	snappb "github.com/zeta-protocol/zeta/protos/zeta/snapshot/v1"
	"github.com/golang/protobuf/proto"
)

func (e *Engine) Namespace() types.SnapshotNamespace {
	return types.ProtocolUpgradeSnapshot
}

func (e *Engine) Keys() []string {
	return e.hashKeys
}

func (e *Engine) Stopped() bool {
	return false
}

// get the serialised form and hash of the given key.
func (e *Engine) serialise() ([]byte, error) {
	events := make([]*eventspb.ProtocolUpgradeEvent, 0, len(e.activeProposals))
	for _, evt := range e.events {
		events = append(events, evt)
	}

	sort.SliceStable(events, func(i, j int) bool {
		if events[i].ZetaReleaseTag == events[j].ZetaReleaseTag {
			return events[i].UpgradeBlockHeight < events[j].UpgradeBlockHeight
		}
		return events[i].ZetaReleaseTag < events[j].ZetaReleaseTag
	})

	payloadProtocolUpgradeProposals := &types.PayloadProtocolUpgradeProposals{
		Proposals: &snappb.ProtocolUpgradeProposals{
			ActiveProposals: events,
		},
	}

	if e.upgradeStatus.AcceptedReleaseInfo != nil {
		payloadProtocolUpgradeProposals.Proposals.AcceptedProposal = &snappb.AcceptedProtocolUpgradeProposal{
			UpgradeBlockHeight: e.upgradeStatus.AcceptedReleaseInfo.UpgradeBlockHeight,
			ZetaReleaseTag:     e.upgradeStatus.AcceptedReleaseInfo.ZetaReleaseTag,
		}
	}

	payload := types.Payload{
		Data: payloadProtocolUpgradeProposals,
	}

	data, err := proto.Marshal(payload.IntoProto())
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *Engine) GetState(k string) ([]byte, []types.StateProvider, error) {
	state, err := e.serialise()
	return state, nil, err
}

func (e *Engine) LoadState(ctx context.Context, p *types.Payload) ([]types.StateProvider, error) {
	if e.Namespace() != p.Data.Namespace() {
		return nil, types.ErrInvalidSnapshotNamespace
	}
	pl := p.Data.(*types.PayloadProtocolUpgradeProposals)
	e.activeProposals = make(map[string]*protocolUpgradeProposal, len(pl.Proposals.ActiveProposals))
	e.events = make(map[string]*eventspb.ProtocolUpgradeEvent, len(pl.Proposals.ActiveProposals))

	for _, pue := range pl.Proposals.ActiveProposals {
		ID := protocolUpgradeProposalID(pue.UpgradeBlockHeight, pue.ZetaReleaseTag)
		e.events[ID] = pue
	}

	for ID, evt := range e.events {
		e.activeProposals[ID] = &protocolUpgradeProposal{
			zetaReleaseTag: evt.ZetaReleaseTag,
			blockHeight:    evt.UpgradeBlockHeight,
			accepted:       make(map[string]struct{}, len(evt.Approvers)),
		}
		for _, approver := range evt.Approvers {
			e.activeProposals[ID].accepted[approver] = struct{}{}
		}
	}

	if pl.Proposals.AcceptedProposal != nil {
		e.upgradeStatus.AcceptedReleaseInfo = &types.ReleaseInfo{
			UpgradeBlockHeight: pl.Proposals.AcceptedProposal.UpgradeBlockHeight,
			ZetaReleaseTag:     pl.Proposals.AcceptedProposal.ZetaReleaseTag,
		}
	}

	return nil, nil
}
