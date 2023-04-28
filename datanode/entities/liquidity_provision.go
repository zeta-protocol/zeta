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
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/encoding/protojson"
)

type _LiquidityProvision struct{}

type LiquidityProvisionID = ID[_LiquidityProvision]

type LiquidityOrderReference struct {
	*zeta.LiquidityOrderReference
}

func (l LiquidityOrderReference) EncodeBinary(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	protoBytes, err := protojson.Marshal(l.LiquidityOrderReference)
	if err != nil {
		return buf, fmt.Errorf("failed to marshal LiquidityOrderReference: %w", err)
	}
	return append(buf, protoBytes...), nil
}

func (l *LiquidityOrderReference) DecodeBinary(_ *pgtype.ConnInfo, src []byte) error {
	return protojson.Unmarshal(src, l)
}

type LiquidityProvision struct {
	ID               LiquidityProvisionID
	PartyID          PartyID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	MarketID         MarketID
	CommitmentAmount decimal.Decimal
	Fee              decimal.Decimal
	Sells            []LiquidityOrderReference
	Buys             []LiquidityOrderReference
	Version          int64
	Status           LiquidityProvisionStatus
	Reference        string
	TxHash           TxHash
	ZetaTime         time.Time
}

func LiquidityProvisionFromProto(lpProto *zeta.LiquidityProvision, txHash TxHash, zetaTime time.Time) (LiquidityProvision, error) {
	lpID := LiquidityProvisionID(lpProto.Id)
	partyID := PartyID(lpProto.PartyId)
	marketID := MarketID(lpProto.MarketId)

	commitmentAmount, err := decimal.NewFromString(lpProto.CommitmentAmount)
	if err != nil {
		return LiquidityProvision{}, fmt.Errorf("liquidity provision has invalid commitement amount: %w", err)
	}

	fee, err := decimal.NewFromString(lpProto.Fee)
	if err != nil {
		return LiquidityProvision{}, fmt.Errorf("liquidity provision has invalid fee amount: %w", err)
	}

	sells := make([]LiquidityOrderReference, 0, len(lpProto.Sells))
	buys := make([]LiquidityOrderReference, 0, len(lpProto.Buys))

	for _, sell := range lpProto.Sells {
		sells = append(sells, LiquidityOrderReference{sell})
	}

	for _, buy := range lpProto.Buys {
		buys = append(buys, LiquidityOrderReference{buy})
	}

	return LiquidityProvision{
		ID:               lpID,
		PartyID:          partyID,
		CreatedAt:        NanosToPostgresTimestamp(lpProto.CreatedAt),
		UpdatedAt:        NanosToPostgresTimestamp(lpProto.UpdatedAt),
		MarketID:         marketID,
		CommitmentAmount: commitmentAmount,
		Fee:              fee,
		Sells:            sells,
		Buys:             buys,
		Version:          int64(lpProto.Version),
		Status:           LiquidityProvisionStatus(lpProto.Status),
		Reference:        lpProto.Reference,
		TxHash:           txHash,
		ZetaTime:         zetaTime,
	}, nil
}

func (lp LiquidityProvision) ToProto() *zeta.LiquidityProvision {
	sells := make([]*zeta.LiquidityOrderReference, 0, len(lp.Sells))
	buys := make([]*zeta.LiquidityOrderReference, 0, len(lp.Buys))

	for _, sell := range lp.Sells {
		sells = append(sells, sell.LiquidityOrderReference)
	}
	for _, buy := range lp.Buys {
		buys = append(buys, buy.LiquidityOrderReference)
	}

	return &zeta.LiquidityProvision{
		Id:               lp.ID.String(),
		PartyId:          lp.PartyID.String(),
		CreatedAt:        lp.CreatedAt.UnixNano(),
		UpdatedAt:        lp.UpdatedAt.UnixNano(),
		MarketId:         lp.MarketID.String(),
		CommitmentAmount: lp.CommitmentAmount.String(),
		Fee:              lp.Fee.String(),
		Sells:            sells,
		Buys:             buys,
		Version:          uint64(lp.Version),
		Status:           zeta.LiquidityProvision_Status(lp.Status),
		Reference:        lp.Reference,
	}
}

type LiquidityProvisionKey struct {
	ID       LiquidityProvisionID
	ZetaTime time.Time
}

func (lp LiquidityProvision) Key() LiquidityProvisionKey {
	return LiquidityProvisionKey{lp.ID, lp.ZetaTime}
}

var LiquidityProvisionColumns = []string{
	"id", "party_id", "created_at", "updated_at", "market_id",
	"commitment_amount", "fee", "sells", "buys", "version",
	"status", "reference", "tx_hash", "zeta_time",
}

func (lp LiquidityProvision) ToRow() []interface{} {
	return []interface{}{
		lp.ID, lp.PartyID, lp.CreatedAt, lp.UpdatedAt, lp.MarketID,
		lp.CommitmentAmount, lp.Fee, lp.Sells, lp.Buys, lp.Version,
		lp.Status, lp.Reference, lp.TxHash, lp.ZetaTime,
	}
}

func (lp LiquidityProvision) Cursor() *Cursor {
	lc := LiquidityProvisionCursor{
		ZetaTime: lp.ZetaTime,
		ID:       lp.ID,
	}
	return NewCursor(lc.String())
}

func (lp LiquidityProvision) ToProtoEdge(_ ...any) (*v2.LiquidityProvisionsEdge, error) {
	return &v2.LiquidityProvisionsEdge{
		Node:   lp.ToProto(),
		Cursor: lp.Cursor().Encode(),
	}, nil
}

type LiquidityProvisionCursor struct {
	ZetaTime time.Time            `json:"zetaTime"`
	ID       LiquidityProvisionID `json:"id"`
}

func (lc LiquidityProvisionCursor) String() string {
	bs, err := json.Marshal(lc)
	if err != nil {
		panic(fmt.Errorf("could not marshal liquidity provision cursor: %w", err))
	}
	return string(bs)
}

func (lc *LiquidityProvisionCursor) Parse(cursorString string) error {
	if cursorString == "" {
		return nil
	}
	return json.Unmarshal([]byte(cursorString), lc)
}
