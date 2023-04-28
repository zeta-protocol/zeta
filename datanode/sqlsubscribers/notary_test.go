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

package sqlsubscribers_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-protocol/zeta/core/events"
	"github.com/zeta-protocol/zeta/datanode/sqlsubscribers"
	"github.com/zeta-protocol/zeta/datanode/sqlsubscribers/mocks"
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
	v1 "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
)

func TestNotary_Push(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockNotaryStore(ctrl)

	store.EXPECT().Add(context.Background(), gomock.Any()).Times(1)
	subscriber := sqlsubscribers.NewNotary(store)
	err := subscriber.Push(context.Background(),
		events.NewNodeSignatureEvent(context.Background(),
			v1.NodeSignature{
				Id:   "someid",
				Sig:  []byte("somesig"),
				Kind: v1.NodeSignatureKind_NODE_SIGNATURE_KIND_ASSET_WITHDRAWAL,
			},
		),
	)
	require.NoError(t, err)
}

func TestNotary_PushWrongEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	store := mocks.NewMockNotaryStore(ctrl)
	subscriber := sqlsubscribers.NewNotary(store)
	subscriber.Push(context.Background(), events.NewOracleDataEvent(context.Background(), zetapb.OracleData{}))
}
