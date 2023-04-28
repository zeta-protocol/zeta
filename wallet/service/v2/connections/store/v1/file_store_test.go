package v1_test

import (
	"testing"
	"time"

	"github.com/zeta-protocol/zeta/libs/ptr"
	vgrand "github.com/zeta-protocol/zeta/libs/rand"
	"github.com/zeta-protocol/zeta/paths"
	"github.com/zeta-protocol/zeta/wallet/service/v2/connections"
	v1 "github.com/zeta-protocol/zeta/wallet/service/v2/connections/store/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStore(t *testing.T) {
	t.Run("List tokens succeeds", testFileStoreListTokensSucceeds)
	t.Run("Verifying an existing token succeeds", testFileStoreVerifyingExistingTokenSucceeds)
	t.Run("Verifying an unknown token fails", testFileStoreVerifyingUnknownTokenFails)
	t.Run("Deleting an existing token succeeds", testFileStoreDeletingExistingTokenSucceeds)
}

func testFileStoreListTokensSucceeds(t *testing.T) {
	zetaPaths := testHome(t)
	store := newTestFileStore(t, zetaPaths)

	// given
	description1 := connections.TokenDescription{
		Description:    vgrand.RandomStr(5),
		CreationDate:   time.Now().Add(-2 * time.Hour),
		ExpirationDate: ptr.From(time.Now().Add(-1 * time.Hour)),
		Token:          connections.GenerateToken(),
		Wallet: connections.WalletCredentials{
			Name:       vgrand.RandomStr(5),
			Passphrase: vgrand.RandomStr(5),
		},
	}

	// when
	err := store.SaveToken(description1)

	// then
	require.NoError(t, err)
	// given
	description2 := connections.TokenDescription{
		Description:  vgrand.RandomStr(5),
		CreationDate: time.Now().Add(-4 * time.Hour),
		Token:        connections.GenerateToken(),
		Wallet: connections.WalletCredentials{
			Name:       vgrand.RandomStr(5),
			Passphrase: vgrand.RandomStr(5),
		},
	}

	// when
	err = store.SaveToken(description2)

	// then
	require.NoError(t, err)
	// when
	tokenSummaries, err := store.ListTokens()

	// then
	require.NoError(t, err)
	assert.Equal(t, description1.Description, tokenSummaries[0].Description)
	assert.Equal(t, description1.Token, tokenSummaries[0].Token)
	assert.WithinDuration(t, description1.CreationDate, tokenSummaries[0].CreationDate, 0)
	assert.WithinDuration(t, *description1.ExpirationDate, *tokenSummaries[0].ExpirationDate, 0)
	assert.Equal(t, description2.Description, tokenSummaries[1].Description)
	assert.Equal(t, description2.Token, tokenSummaries[1].Token)
	assert.WithinDuration(t, description2.CreationDate, tokenSummaries[1].CreationDate, 0)
	assert.Equal(t, description2.ExpirationDate, tokenSummaries[1].ExpirationDate, 0)
}

func testFileStoreVerifyingExistingTokenSucceeds(t *testing.T) {
	zetaPaths := testHome(t)
	store := newTestFileStore(t, zetaPaths)

	// given
	description1 := connections.TokenDescription{
		Description:    vgrand.RandomStr(5),
		CreationDate:   time.Now().Add(-2 * time.Hour),
		ExpirationDate: ptr.From(time.Now().Add(-1 * time.Hour)),
		Token:          connections.GenerateToken(),
		Wallet: connections.WalletCredentials{
			Name:       vgrand.RandomStr(5),
			Passphrase: vgrand.RandomStr(5),
		},
	}

	// when
	err := store.SaveToken(description1)

	// then
	require.NoError(t, err)

	// when
	exists, err := store.TokenExists(description1.Token)

	// then
	require.NoError(t, err)
	assert.True(t, exists)
}

func testFileStoreVerifyingUnknownTokenFails(t *testing.T) {
	zetaPaths := testHome(t)
	store := newTestFileStore(t, zetaPaths)

	// given
	unknownToken := connections.GenerateToken()

	// when
	exists, err := store.TokenExists(unknownToken)

	// then
	require.NoError(t, err)
	assert.False(t, exists)
}

func testFileStoreDeletingExistingTokenSucceeds(t *testing.T) {
	zetaPaths := testHome(t)
	store := newTestFileStore(t, zetaPaths)

	// given
	description1 := connections.TokenDescription{
		Description:    vgrand.RandomStr(5),
		CreationDate:   time.Now().Add(-2 * time.Hour),
		ExpirationDate: ptr.From(time.Now().Add(-1 * time.Hour)),
		Token:          connections.GenerateToken(),
		Wallet: connections.WalletCredentials{
			Name:       vgrand.RandomStr(5),
			Passphrase: vgrand.RandomStr(5),
		},
	}

	// when
	err := store.SaveToken(description1)

	// then
	require.NoError(t, err)

	// given
	description2 := connections.TokenDescription{
		Description:  vgrand.RandomStr(5),
		CreationDate: time.Now().Add(-4 * time.Hour),
		Token:        connections.GenerateToken(),
		Wallet: connections.WalletCredentials{
			Name:       vgrand.RandomStr(5),
			Passphrase: vgrand.RandomStr(5),
		},
	}

	// when
	err = store.SaveToken(description2)

	// then
	require.NoError(t, err)

	// when
	err = store.DeleteToken(description1.Token)

	// then
	require.NoError(t, err)

	// when
	exists, err := store.TokenExists(description1.Token)

	// then
	require.NoError(t, err)
	assert.False(t, exists)

	// when
	exists, err = store.TokenExists(description2.Token)

	// then
	require.NoError(t, err)
	assert.True(t, exists)
}

type testFileStore struct {
	*v1.FileStore
}

func testHome(t *testing.T) paths.Paths {
	t.Helper()
	return paths.New(t.TempDir())
}

func newTestFileStore(t *testing.T, zetaPaths paths.Paths) *testFileStore {
	t.Helper()

	tokenStore, err := v1.ReinitialiseStore(zetaPaths, vgrand.RandomStr(5))
	if err != nil {
		t.Fatalf("could not initialise the file store for tests: %v", err)
	}
	t.Cleanup(func() {
		tokenStore.Close()
	})

	return &testFileStore{
		FileStore: tokenStore,
	}
}
