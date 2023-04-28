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

	"github.com/zeta-protocol/zeta/protos/zeta"
)

type updateFutureProductResolver ZetaResolverRoot

func (r *updateFutureProductResolver) DataSourceSpecForSettlementData(_ context.Context, obj *zeta.UpdateFutureProduct) (*DataSourceDefinition, error) {
	if obj.DataSourceSpecForSettlementData == nil {
		return nil, nil
	}
	return resolveDataSourceDefinition(obj.DataSourceSpecForSettlementData), nil
}

func (r *updateFutureProductResolver) DataSourceSpecForTradingTermination(_ context.Context, obj *zeta.UpdateFutureProduct) (*DataSourceDefinition, error) {
	if obj.DataSourceSpecForTradingTermination == nil {
		return nil, nil
	}
	return resolveDataSourceDefinition(obj.DataSourceSpecForTradingTermination), nil
}
