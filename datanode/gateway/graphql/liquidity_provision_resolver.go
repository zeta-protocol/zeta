package gql

import (
	"context"
	"strconv"

	"github.com/zeta-protocol/zeta/protos/zeta"
)

//type LiquidityProvisionUpdateResolver interface {
//	CreatedAt(ctx context.Context, obj *zeta.LiquidityProvision) (string, error)
//	UpdatedAt(ctx context.Context, obj *zeta.LiquidityProvision) (*string, error)
//
//	Version(ctx context.Context, obj *zeta.LiquidityProvision) (string, error)
//}

type liquidityProvisionUpdateResolver ZetaResolverRoot

func (r *liquidityProvisionUpdateResolver) CreatedAt(ctx context.Context, obj *zeta.LiquidityProvision) (string, error) {
	return strconv.FormatInt(obj.CreatedAt, 10), nil
}

func (r *liquidityProvisionUpdateResolver) UpdatedAt(ctx context.Context, obj *zeta.LiquidityProvision) (*string, error) {
	if obj.UpdatedAt == 0 {
		return nil, nil
	}

	updatedAt := strconv.FormatInt(obj.UpdatedAt, 10)

	return &updatedAt, nil
}

func (r *liquidityProvisionUpdateResolver) Version(ctx context.Context, obj *zeta.LiquidityProvision) (string, error) {
	return strconv.FormatUint(obj.Version, 10), nil
}
