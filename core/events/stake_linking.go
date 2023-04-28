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

package events

import (
	"context"

	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"

	"github.com/zeta-protocol/zeta/core/types"
)

type StakeLinking struct {
	*Base
	evt eventspb.StakeLinking
}

func NewStakeLinking(ctx context.Context, evt types.StakeLinking) *StakeLinking {
	return &StakeLinking{
		Base: newBase(ctx, StakeLinkingEvent),
		evt:  *(evt.IntoProto()),
	}
}

func (s StakeLinking) StakeLinking() eventspb.StakeLinking {
	return s.evt
}

func (s StakeLinking) Proto() eventspb.StakeLinking {
	return s.evt
}

func (s StakeLinking) StreamMessage() *eventspb.BusEvent {
	busEvent := newBusEventFromBase(s.Base)
	busEvent.Event = &eventspb.BusEvent_StakeLinking{
		StakeLinking: &s.evt,
	}

	return busEvent
}

func StakeLinkingFromStream(ctx context.Context, be *eventspb.BusEvent) *StakeLinking {
	return &StakeLinking{
		Base: newBaseFromBusEvent(ctx, StakeLinkingEvent, be),
		evt:  *be.GetStakeLinking(),
	}
}
