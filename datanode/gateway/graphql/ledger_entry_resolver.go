package gql

import (
	"context"

	zeta "github.com/zeta-protocol/zeta/protos/zeta"
)

type ledgerEntryResolver ZetaResolverRoot

func (le ledgerEntryResolver) FromAccountID(ctx context.Context, obj *zeta.LedgerEntry) (*zeta.AccountDetails, error) {
	return obj.FromAccount, nil
}

func (le ledgerEntryResolver) ToAccountID(ctx context.Context, obj *zeta.LedgerEntry) (*zeta.AccountDetails, error) {
	return obj.ToAccount, nil
}
