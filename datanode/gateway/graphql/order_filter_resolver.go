package gql

import (
	"context"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
	"github.com/zeta-protocol/zeta/protos/zeta"
)

type orderFilterResolver ZetaResolverRoot

func (o orderFilterResolver) Status(ctx context.Context, obj *v2.OrderFilter, data []zeta.Order_Status) error {
	obj.Statuses = data
	return nil
}

func (o orderFilterResolver) TimeInForce(ctx context.Context, obj *v2.OrderFilter, data []zeta.Order_TimeInForce) error {
	obj.TimeInForces = data
	return nil
}
