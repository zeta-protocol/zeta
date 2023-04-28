package gql

import (
	"context"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
)

type dateRangeResolver ZetaResolverRoot

func (r *dateRangeResolver) Start(ctx context.Context, obj *v2.DateRange, data *int64) error {
	obj.StartTimestamp = data
	return nil
}

func (r *dateRangeResolver) End(ctx context.Context, obj *v2.DateRange, data *int64) error {
	obj.EndTimestamp = data
	return nil
}
