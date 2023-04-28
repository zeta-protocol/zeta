package cmd_test

import (
	"testing"

	cmd "github.com/zeta-protocol/zeta/cmd/zetawallet/commands"
	"github.com/zeta-protocol/zeta/cmd/zetawallet/commands/flags"
	vgrand "github.com/zeta-protocol/zeta/libs/rand"
	"github.com/zeta-protocol/zeta/wallet/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWalletFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testCreateWalletFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testCreateWalletFlagsMissingWalletFails)
}

func testCreateWalletFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	walletName := vgrand.RandomStr(10)
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	f := &cmd.CreateWalletFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := api.AdminCreateWalletParams{
		Wallet:     walletName,
		Passphrase: passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testCreateWalletFlagsMissingWalletFails(t *testing.T) {
	// given
	f := newCreateWalletFlags(t)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, req)
}

func newCreateWalletFlags(t *testing.T) *cmd.CreateWalletFlags {
	t.Helper()
	return &cmd.CreateWalletFlags{
		Wallet:         vgrand.RandomStr(10),
		PassphraseFile: "/some/fake/path",
	}
}
