// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package gql

import (
	"context"
	"fmt"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
	"github.com/zeta-protocol/zeta/protos/zeta"
)

type rewardSummaryResolver ZetaResolverRoot

func (r *rewardSummaryResolver) Asset(ctx context.Context, obj *zeta.RewardSummary) (*zeta.Asset, error) {
	return r.r.getAssetByID(ctx, obj.AssetId)
}

func (r *rewardSummaryResolver) RewardsConnection(ctx context.Context, summary *zeta.RewardSummary, assetID *string, pagination *v2.Pagination) (*v2.RewardsConnection, error) {
	req := v2.ListRewardsRequest{
		PartyId:    summary.PartyId,
		AssetId:    assetID,
		Pagination: pagination,
	}
	resp, err := r.tradingDataClientV2.ListRewards(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rewards information: %w", err)
	}

	return resp.Rewards, nil
}
