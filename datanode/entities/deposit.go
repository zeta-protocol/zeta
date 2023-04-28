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

package entities

import (
	"encoding/json"
	"fmt"
	"time"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
	"github.com/zeta-protocol/zeta/protos/zeta"
	"github.com/shopspring/decimal"
)

type _Deposit struct{}

type DepositID = ID[_Deposit]

type Deposit struct {
	ID                DepositID
	Status            DepositStatus
	PartyID           PartyID
	Asset             AssetID
	Amount            decimal.Decimal
	ForeignTxHash     string
	CreditedTimestamp time.Time
	CreatedTimestamp  time.Time
	TxHash            TxHash
	ZetaTime          time.Time
}

func DepositFromProto(deposit *zeta.Deposit, txHash TxHash, zetaTime time.Time) (*Deposit, error) {
	var err error
	var amount decimal.Decimal

	if amount, err = decimal.NewFromString(deposit.Amount); err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	return &Deposit{
		ID:                DepositID(deposit.Id),
		Status:            DepositStatus(deposit.Status),
		PartyID:           PartyID(deposit.PartyId),
		Asset:             AssetID(deposit.Asset),
		Amount:            amount,
		ForeignTxHash:     deposit.TxHash,
		CreditedTimestamp: NanosToPostgresTimestamp(deposit.CreditedTimestamp),
		CreatedTimestamp:  NanosToPostgresTimestamp(deposit.CreatedTimestamp),
		TxHash:            txHash,
		ZetaTime:          zetaTime,
	}, nil
}

func (d Deposit) ToProto() *zeta.Deposit {
	return &zeta.Deposit{
		Id:                d.ID.String(),
		Status:            zeta.Deposit_Status(d.Status),
		PartyId:           d.PartyID.String(),
		Asset:             d.Asset.String(),
		Amount:            d.Amount.String(),
		TxHash:            d.ForeignTxHash,
		CreditedTimestamp: d.CreditedTimestamp.UnixNano(),
		CreatedTimestamp:  d.CreatedTimestamp.UnixNano(),
	}
}

func (d Deposit) Cursor() *Cursor {
	cursor := DepositCursor{
		ZetaTime: d.ZetaTime,
		ID:       d.ID,
	}
	return NewCursor(cursor.String())
}

func (d Deposit) ToProtoEdge(_ ...any) (*v2.DepositEdge, error) {
	return &v2.DepositEdge{
		Node:   d.ToProto(),
		Cursor: d.Cursor().Encode(),
	}, nil
}

type DepositCursor struct {
	ZetaTime time.Time `json:"zetaTime"`
	ID       DepositID `json:"id"`
}

func (dc DepositCursor) String() string {
	bs, err := json.Marshal(dc)
	if err != nil {
		// This should never happen.
		panic(fmt.Errorf("couldn't marshal deposit cursor: %w", err))
	}
	return string(bs)
}

func (dc *DepositCursor) Parse(cursorString string) error {
	if cursorString == "" {
		return nil
	}
	return json.Unmarshal([]byte(cursorString), dc)
}
