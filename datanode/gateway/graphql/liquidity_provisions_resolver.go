package gql

import (
	"context"
	"strconv"

	"github.com/zeta-protocol/zeta/datanode/zetatime"
	types "github.com/zeta-protocol/zeta/protos/zeta"
)

// LiquidityProvision resolver

type myLiquidityProvisionResolver ZetaResolverRoot

func (r *myLiquidityProvisionResolver) Version(_ context.Context, obj *types.LiquidityProvision) (string, error) {
	return strconv.FormatUint(obj.Version, 10), nil
}

func (r *myLiquidityProvisionResolver) Party(_ context.Context, obj *types.LiquidityProvision) (*types.Party, error) {
	return &types.Party{Id: obj.PartyId}, nil
}

func (r *myLiquidityProvisionResolver) CreatedAt(ctx context.Context, obj *types.LiquidityProvision) (string, error) {
	return zetatime.Format(zetatime.UnixNano(obj.CreatedAt)), nil
}

func (r *myLiquidityProvisionResolver) UpdatedAt(ctx context.Context, obj *types.LiquidityProvision) (*string, error) {
	var updatedAt *string
	if obj.UpdatedAt > 0 {
		t := zetatime.Format(zetatime.UnixNano(obj.UpdatedAt))
		updatedAt = &t
	}
	return updatedAt, nil
}

func (r *myLiquidityProvisionResolver) Market(ctx context.Context, obj *types.LiquidityProvision) (*types.Market, error) {
	return r.r.getMarketByID(ctx, obj.MarketId)
}

func (r *myLiquidityProvisionResolver) CommitmentAmount(ctx context.Context, obj *types.LiquidityProvision) (string, error) {
	return obj.CommitmentAmount, nil
}
