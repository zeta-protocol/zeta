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

func TestListKeysFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testListKeysFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testListKeysFlagsMissingWalletFails)
}

func testListKeysFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	expectedPassphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	f := &cmd.ListKeysFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := api.AdminListKeysParams{
		Wallet: walletName,
	}

	// when
	req, passphrase, err := f.Validate()

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedReq, req)
	assert.Equal(t, expectedPassphrase, passphrase)
}

func testListKeysFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newListKeysFlags(t, testDir)
	f.Wallet = ""

	// when
	req, _, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, req)
}

func newListKeysFlags(t *testing.T, testDir string) *cmd.ListKeysFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	return &cmd.ListKeysFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}
}
