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

	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"
)

type OracleSpec struct {
	*Base
	o zetapb.OracleSpec
}

func NewOracleSpecEvent(ctx context.Context, spec zetapb.OracleSpec) *OracleSpec {
	cpy := spec.DeepClone()
	return &OracleSpec{
		Base: newBase(ctx, OracleSpecEvent),
		o:    *cpy,
	}
}

func (o *OracleSpec) OracleSpec() zetapb.OracleSpec {
	return o.o
}

func (o OracleSpec) Proto() zetapb.OracleSpec {
	return o.o
}

func (o OracleSpec) StreamMessage() *eventspb.BusEvent {
	spec := o.o

	busEvent := newBusEventFromBase(o.Base)
	busEvent.Event = &eventspb.BusEvent_OracleSpec{
		OracleSpec: &spec,
	}

	return busEvent
}

func OracleSpecEventFromStream(ctx context.Context, be *eventspb.BusEvent) *OracleSpec {
	return &OracleSpec{
		Base: newBaseFromBusEvent(ctx, OracleSpecEvent, be),
		o:    *be.GetOracleSpec(),
	}
}
