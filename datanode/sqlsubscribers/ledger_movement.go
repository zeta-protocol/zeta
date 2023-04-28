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
	"fmt"
	"strings"
	"time"

	"github.com/zeta-protocol/zeta/core/events"
	"github.com/zeta-protocol/zeta/datanode/entities"
	"github.com/zeta-protocol/zeta/protos/zeta"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Ledger interface {
	AddLedgerEntry(entities.LedgerEntry) error
	AddTransferResponse(*zeta.LedgerMovement)
	Flush(ctx context.Context) error
}

type TransferResponseEvent interface {
	events.Event
	LedgerMovements() []*zeta.LedgerMovement
}

type TransferResponse struct {
	subscriber
	ledger   Ledger
	accounts AccountService
}

func NewTransferResponse(ledger Ledger, accounts AccountService) *TransferResponse {
	return &TransferResponse{
		ledger:   ledger,
		accounts: accounts,
	}
}

func (t *TransferResponse) Types() []events.Type {
	return []events.Type{events.LedgerMovementsEvent}
}

func (t *TransferResponse) Flush(ctx context.Context) error {
	err := t.ledger.Flush(ctx)
	return errors.Wrap(err, "flushing ledger")
}

func (t *TransferResponse) Push(ctx context.Context, evt events.Event) error {
	return t.consume(ctx, evt.(TransferResponseEvent))
}

func (t *TransferResponse) consume(ctx context.Context, e TransferResponseEvent) error {
	var errs strings.Builder
	for _, tr := range e.LedgerMovements() {
		t.ledger.AddTransferResponse(tr)
		for _, vle := range tr.Entries {
			if err := t.addLedgerEntry(ctx, vle, e.TxHash(), t.zetaTime); err != nil {
				errs.WriteString(fmt.Sprintf("couldn't add ledger entry: %v, error:%s\n", vle, err))
			}
		}
	}

	if errs.Len() != 0 {
		return errors.Errorf("processing transfer response:%s", errs.String())
	}

	return nil
}

func (t *TransferResponse) addLedgerEntry(ctx context.Context, vle *zeta.LedgerEntry, txHash string, zetaTime time.Time) error {
	fromAcc, err := t.obtainAccountWithAccountDetails(ctx, vle.FromAccount, txHash, zetaTime)
	if err != nil {
		return errors.Wrap(err, "obtaining 'from' account")
	}

	toAcc, err := t.obtainAccountWithAccountDetails(ctx, vle.ToAccount, txHash, zetaTime)
	if err != nil {
		return errors.Wrap(err, "obtaining 'to' account")
	}

	quantity, err := decimal.NewFromString(vle.Amount)
	if err != nil {
		return errors.Wrap(err, "parsing amount string")
	}

	fromAccountBalance, err := decimal.NewFromString(vle.FromAccountBalance)
	if err != nil {
		return errors.Wrap(err, "parsing FromAccountBalance string")
	}

	toAccountBalance, err := decimal.NewFromString(vle.ToAccountBalance)
	if err != nil {
		return errors.Wrap(err, "parsing ToAccountBalance string")
	}

	le := entities.LedgerEntry{
		FromAccountID:      fromAcc.ID,
		ToAccountID:        toAcc.ID,
		Quantity:           quantity,
		TxHash:             entities.TxHash(txHash),
		ZetaTime:           zetaTime,
		TransferTime:       time.Unix(0, vle.Timestamp),
		Type:               entities.LedgerMovementType(vle.Type),
		FromAccountBalance: fromAccountBalance,
		ToAccountBalance:   toAccountBalance,
	}

	err = t.ledger.AddLedgerEntry(le)
	if err != nil {
		return errors.Wrap(err, "adding to store")
	}
	return nil
}

// Parse the zeta account ID; if that account already exists in the db, fetch it; else create it.
func (t *TransferResponse) obtainAccountWithAccountDetails(ctx context.Context, ad *zeta.AccountDetails, txHash string, zetaTime time.Time) (entities.Account, error) {
	a, err := entities.AccountProtoFromDetails(ad, entities.TxHash(txHash))
	if err != nil {
		return entities.Account{}, errors.Wrapf(err, "parsing account id: %s", ad.String())
	}
	a.ZetaTime = zetaTime
	err = t.accounts.Obtain(ctx, &a)
	if err != nil {
		return entities.Account{}, errors.Wrapf(err, "obtaining account for id: %s", ad.String())
	}
	return a, nil
}
