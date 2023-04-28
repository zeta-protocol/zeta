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

package target_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	snapshot "github.com/zeta-protocol/zeta/protos/zeta/snapshot/v1"

	"github.com/zeta-protocol/zeta/core/liquidity/target"
	"github.com/zeta-protocol/zeta/core/types"
	"github.com/zeta-protocol/zeta/libs/num"
	"github.com/zeta-protocol/zeta/libs/proto"
	"github.com/stretchr/testify/assert"
)

func newSnapshotEngine(marketID string) *target.SnapshotEngine {
	params := types.TargetStakeParameters{
		TimeWindow:    5,
		ScalingFactor: num.NewDecimalFromFloat(2),
	}
	var oiCalc target.OpenInterestCalculator

	return target.NewSnapshotEngine(params, oiCalc, marketID, num.DecimalFromFloat(1))
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	a := assert.New(t)
	marketID := "market-1"
	key := fmt.Sprintf("target:%s", marketID)
	se := newSnapshotEngine(marketID)

	s, _, err := se.GetState("")
	a.Empty(s)
	a.EqualError(err, types.ErrSnapshotKeyDoesNotExist.Error())

	d := time.Date(2015, time.December, 24, 19, 0, 0, 0, time.UTC)
	se.RecordOpenInterest(40, d)
	se.RecordOpenInterest(40, d.Add(time.Hour*3))

	s, _, err = se.GetState(key)
	a.NotEmpty(s)
	a.NoError(err)

	se2 := newSnapshotEngine(marketID)

	pl := snapshot.Payload{}
	assert.NoError(t, proto.Unmarshal(s, &pl))

	_, err = se2.LoadState(context.TODO(), types.PayloadFromProto(&pl))
	a.NoError(err)

	s2, _, err := se2.GetState(key)
	a.NoError(err)
	a.True(bytes.Equal(s, s2))
}

func TestStopSnapshotTaking(t *testing.T) {
	marketID := "market-1"
	key := fmt.Sprintf("target:%s", marketID)
	se := newSnapshotEngine(marketID)

	// signal to kill the engine's snapshots
	se.StopSnapshots()

	s, _, err := se.GetState(key)
	assert.NoError(t, err)
	assert.Nil(t, s)
	assert.True(t, se.Stopped())
}
