package cmd_test

import (
	"testing"

	cmd "github.com/zeta-protocol/zeta/cmd/zetawallet/commands"
	"github.com/stretchr/testify/require"
)

func TestInitAPIToken(t *testing.T) {
	t.Run("Initialising software succeeds", testInitialisingAPITokenSucceeds)
	t.Run("Forcing software initialisation succeeds", testForcingAPITokenInitialisationSucceeds)
}

func testInitialisingAPITokenSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	f := &cmd.InitAPITokenFlags{
		Force:          false,
		PassphraseFile: passphraseFilePath,
	}

	// when
	err := cmd.InitAPIToken(testDir, f)

	// then
	require.NoError(t, err)
}

func testForcingAPITokenInitialisationSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	f := &cmd.InitAPITokenFlags{
		Force:          false,
		PassphraseFile: passphraseFilePath,
	}

	// when
	err := cmd.InitAPIToken(testDir, f)

	// then
	require.NoError(t, err)

	// given
	f = &cmd.InitAPITokenFlags{
		Force:          true,
		PassphraseFile: passphraseFilePath,
	}

	// when
	err = cmd.InitAPIToken(testDir, f)

	// then
	require.NoError(t, err)
}
