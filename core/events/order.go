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

	"github.com/zeta-protocol/zeta/core/types"
	"github.com/zeta-protocol/zeta/libs/num"
	ptypes "github.com/zeta-protocol/zeta/protos/zeta"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"
)

type Order struct {
	*Base
	o *ptypes.Order
}

func NewOrderEvent(ctx context.Context, o *types.Order) *Order {
	order := &Order{
		Base: newBase(ctx, OrderEvent),
		o:    o.IntoProto(),
	}
	// set to original order price
	order.o.Price = num.UintToString(o.OriginalPrice)
	return order
}

func (o Order) IsParty(id string) bool {
	return o.o.PartyId == id
}

func (o Order) PartyID() string {
	return o.o.PartyId
}

func (o Order) MarketID() string {
	return o.o.MarketId
}

func (o *Order) Order() *ptypes.Order {
	return o.o
}

func (o Order) Proto() ptypes.Order {
	return *o.o
}

func (o Order) StreamMessage() *eventspb.BusEvent {
	busEvent := newBusEventFromBase(o.Base)
	busEvent.Event = &eventspb.BusEvent_Order{
		Order: o.o,
	}

	return busEvent
}

func OrderEventFromStream(ctx context.Context, be *eventspb.BusEvent) *Order {
	order := &Order{
		Base: newBaseFromBusEvent(ctx, OrderEvent, be),
		o:    be.GetOrder(),
	}
	return order
}
