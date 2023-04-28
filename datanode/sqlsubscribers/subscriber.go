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

package sqlsubscribers

import (
	"context"
	"time"
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/mocks.go -package mocks github.com/zeta-protocol/zeta/datanode/sqlsubscribers RiskFactorStore,TransferStore,WithdrawalStore,LiquidityProvisionStore,KeyRotationStore,OracleSpecStore,DepositStore,StakeLinkingStore,MarketDataStore,PositionStore,OracleDataStore,MarginLevelsStore,NotaryStore,NodeStore,MarketsStore,MarketSvc

type subscriber struct {
	zetaTime time.Time
}

func (s *subscriber) SetZetaTime(zetaTime time.Time) {
	s.zetaTime = zetaTime
}

func (s *subscriber) Flush(ctx context.Context) error {
	return nil
}
