package connections

import (
	"context"
	"time"

	"github.com/zeta-protocol/zeta/wallet/wallet"
)

// Generates mocks
//go:generate go run github.com/golang/mock/mockgen -destination mocks/mocks.go -package mocks github.com/zeta-protocol/zeta/wallet/service/v2/connections WalletStore,TimeService,TokenStore

type TimeService interface {
	Now() time.Time
}

type WalletStore interface {
	UnlockWallet(ctx context.Context, name, passphrase string) error
	GetWallet(ctx context.Context, name string) (wallet.Wallet, error)
	OnUpdate(callbackFn func(context.Context, wallet.Event))
}

// TokenStore is the component used to retrieve and update the API tokens from the
// computer.
type TokenStore interface {
	TokenExists(token Token) (bool, error)
	ListTokens() ([]TokenSummary, error)
	DescribeToken(token Token) (TokenDescription, error)
	SaveToken(tokenConfig TokenDescription) error
	DeleteToken(token Token) error
	OnUpdate(callbackFn func(ctx context.Context, tokens ...TokenDescription))
}
