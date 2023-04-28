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

func TestRenameNetworkFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testRenameNetworkFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testRenameNetworkFlagsMissingNetworkFails)
	t.Run("Missing new name fails", testRenameNetworkFlagsMissingNewNameFails)
}

func testRenameNetworkFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	walletName := vgrand.RandomStr(10)
	newName := vgrand.RandomStr(10)
	f := &cmd.RenameNetworkFlags{
		Network: walletName,
		NewName: newName,
	}

	expectedReq := api.AdminRenameNetworkParams{
		Network: walletName,
		NewName: newName,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testRenameNetworkFlagsMissingNetworkFails(t *testing.T) {
	// given
	f := newRenameNetworkFlags(t)
	f.Network = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("network"))
	assert.Empty(t, req)
}

func testRenameNetworkFlagsMissingNewNameFails(t *testing.T) {
	// given
	f := newRenameNetworkFlags(t)
	f.NewName = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBeSpecifiedError("new-name"))
	assert.Empty(t, req)
}

func newRenameNetworkFlags(t *testing.T) *cmd.RenameNetworkFlags {
	t.Helper()
	return &cmd.RenameNetworkFlags{
		Network: vgrand.RandomStr(10),
		NewName: vgrand.RandomStr(10),
	}
}
