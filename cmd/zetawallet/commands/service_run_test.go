package cmd_test

import (
	"testing"

	cmd "github.com/zeta-protocol/zeta/cmd/zetawallet/commands"
	"github.com/zeta-protocol/zeta/cmd/zetawallet/commands/flags"
	vgrand "github.com/zeta-protocol/zeta/libs/rand"
	"github.com/stretchr/testify/assert"
)

func TestRunServiceFlags(t *testing.T) {
	t.Run("Missing loads-token flag with tokens passphrase flag fails", testRunServiceFlagsTokenPassphraseWithoutWithLOngLivingTokenFails)
	t.Run("Missing network fails", testRunServiceFlagsMissingNetworkFails)
}

func testRunServiceFlagsTokenPassphraseWithoutWithLOngLivingTokenFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	networkName := vgrand.RandomStr(10)
	f := &cmd.RunServiceFlags{
		Network:              networkName,
		TokensPassphraseFile: passphraseFilePath,
	}

	// when
	err := f.Validate(&cmd.RootFlags{
		Home: testDir,
	})

	// then
	assert.ErrorIs(t, err, flags.OneOfParentsFlagMustBeSpecifiedError("tokens-passphrase-file", "load-tokens"))
}

func testRunServiceFlagsMissingNetworkFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newRunServiceFlags(t)
	f.Network = ""

	// when
	err := f.Validate(&cmd.RootFlags{
		Home: testDir,
	})

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("network"))
}

func newRunServiceFlags(t *testing.T) *cmd.RunServiceFlags {
	t.Helper()

	networkName := vgrand.RandomStr(10)

	return &cmd.RunServiceFlags{
		Network: networkName,
	}
}
