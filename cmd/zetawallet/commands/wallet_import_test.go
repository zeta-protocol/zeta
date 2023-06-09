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

const recoveryPhrase = "swing ceiling chaos green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render"

func TestImportWalletFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testImportWalletFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testImportWalletFlagsMissingWalletFails)
	t.Run("Missing recovery phrase file fails", testImportWalletFlagsMissingRecoveryPhraseFileFails)
}

func testImportWalletFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	recoveryPhraseFilePath := NewFile(t, testDir, "recovery-phrase.txt", recoveryPhrase)
	walletName := vgrand.RandomStr(10)

	f := &cmd.ImportWalletFlags{
		Wallet:             walletName,
		RecoveryPhraseFile: recoveryPhraseFilePath,
		PassphraseFile:     passphraseFilePath,
	}

	expectedReq := api.AdminImportWalletParams{
		Wallet:         walletName,
		RecoveryPhrase: recoveryPhrase,
		Passphrase:     passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedReq, req)
}

func testImportWalletFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newImportWalletFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("wallet"))
	assert.Empty(t, req)
}

func testImportWalletFlagsMissingRecoveryPhraseFileFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newImportWalletFlags(t, testDir)
	f.RecoveryPhraseFile = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("recovery-phrase-file"))
	assert.Empty(t, req)
}

func newImportWalletFlags(t *testing.T, testDir string) *cmd.ImportWalletFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	NewFile(t, testDir, "recovery-phrase.txt", recoveryPhrase)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	return &cmd.ImportWalletFlags{
		Wallet:             walletName,
		RecoveryPhraseFile: pubKey,
		PassphraseFile:     passphraseFilePath,
	}
}
