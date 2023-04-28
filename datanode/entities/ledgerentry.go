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
	"context"
	"fmt"
	"time"

	"github.com/zeta-protocol/zeta/protos/zeta"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type LedgerEntry struct {
	LedgerEntryTime    time.Time
	FromAccountID      AccountID `db:"account_from_id"`
	ToAccountID        AccountID `db:"account_to_id"`
	Quantity           decimal.Decimal
	TxHash             TxHash
	ZetaTime           time.Time
	TransferTime       time.Time
	Type               LedgerMovementType
	FromAccountBalance decimal.Decimal `db:"account_from_balance"`
	ToAccountBalance   decimal.Decimal `db:"account_to_balance"`
}

var LedgerEntryColumns = []string{
	"ledger_entry_time",
	"account_from_id", "account_to_id", "quantity",
	"tx_hash", "zeta_time", "transfer_time", "type",
	"account_from_balance",
	"account_to_balance",
}

func (le LedgerEntry) ToProto(ctx context.Context, accountSource AccountSource) (*zeta.LedgerEntry, error) {
	fromAcc, err := accountSource.GetByID(ctx, le.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("getting from account for transfer proto:%w", err)
	}

	toAcc, err := accountSource.GetByID(ctx, le.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("getting to account for transfer proto:%w", err)
	}

	return &zeta.LedgerEntry{
		FromAccount:        fromAcc.ToAccountDetailsProto(),
		ToAccount:          toAcc.ToAccountDetailsProto(),
		Amount:             le.Quantity.String(),
		Type:               zeta.TransferType(le.Type),
		FromAccountBalance: le.FromAccountBalance.String(),
		ToAccountBalance:   le.ToAccountBalance.String(),
	}, nil
}

func (le LedgerEntry) ToRow() []any {
	return []any{
		le.LedgerEntryTime,
		le.FromAccountID,
		le.ToAccountID,
		le.Quantity,
		le.TxHash,
		le.ZetaTime,
		le.TransferTime,
		le.Type,
		le.FromAccountBalance,
		le.ToAccountBalance,
	}
}

func CreateLedgerEntryTime(zetaTime time.Time, seqNum int) time.Time {
	return zetaTime.Add(time.Duration(seqNum) * time.Microsecond)
}

type LedgerMovementType zeta.TransferType

const (
	// Default value, always invalid.
	LedgerMovementTypeUnspecified = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_UNSPECIFIED)
	// Loss.
	LedgerMovementTypeLoss = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_LOSS)
	// Win.
	LedgerMovementTypeWin = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_WIN)
	// Mark to market loss.
	LedgerMovementTypeMTMLoss = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_MTM_LOSS)
	// Mark to market win.
	LedgerMovementTypeMTMWin = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_MTM_WIN)
	// Margin too low.
	LedgerMovementTypeMarginLow = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_MARGIN_LOW)
	// Margin too high.
	LedgerMovementTypeMarginHigh = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_MARGIN_HIGH)
	// Margin was confiscated.
	LedgerMovementTypeMarginConfiscated = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_MARGIN_CONFISCATED)
	// Pay maker fee.
	LedgerMovementTypeMakerFeePay = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_MAKER_FEE_PAY)
	// Receive maker fee.
	LedgerMovementTypeMakerFeeReceive = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_MAKER_FEE_RECEIVE)
	// Pay infrastructure fee.
	LedgerMovementTypeInfrastructureFeePay = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_INFRASTRUCTURE_FEE_PAY)
	// Receive infrastructure fee.
	LedgerMovementTypeInfrastructureFeeDistribute = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_INFRASTRUCTURE_FEE_DISTRIBUTE)
	// Pay liquidity fee.
	LedgerMovementTypeLiquidityFeePay = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_LIQUIDITY_FEE_PAY)
	// Receive liquidity fee.
	LedgerMovementTypeLiquidityFeeDistribute = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_LIQUIDITY_FEE_DISTRIBUTE)
	// Bond too low.
	LedgerMovementTypeBondLow = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_BOND_LOW)
	// Bond too high.
	LedgerMovementTypeBondHigh = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_BOND_HIGH)
	// Actual withdraw from system.
	LedgerMovementTypeWithdraw = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_WITHDRAW)
	// Deposit funds.
	LedgerMovementTypeDeposit = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_DEPOSIT)
	// Bond slashing.
	LedgerMovementTypeBondSlashing = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_BOND_SLASHING)
	// Reward payout.
	LedgerMovementTypeRewardPayout            = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_REWARD_PAYOUT)
	LedgerMovementTypeTransferFundsSend       = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_TRANSFER_FUNDS_SEND)
	LedgerMovementTypeTransferFundsDistribute = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_TRANSFER_FUNDS_DISTRIBUTE)
	LedgerMovementTypeClearAccount            = LedgerMovementType(zeta.TransferType_TRANSFER_TYPE_CLEAR_ACCOUNT)
)

func (l LedgerMovementType) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	ty, ok := zeta.TransferType_name[int32(l)]
	if !ok {
		return buf, fmt.Errorf("unknown transfer status: %s", ty)
	}
	return append(buf, []byte(ty)...), nil
}

func (l *LedgerMovementType) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.TransferType_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown transfer status: %s", src)
	}

	*l = LedgerMovementType(val)
	return nil
}
