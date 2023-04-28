// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.ZETA file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

//go:build !race
// +build !race

package nodewallets_test

import (
	"testing"

	"github.com/zeta-protocol/zeta/core/nodewallets"
	vgrand "github.com/zeta-protocol/zeta/libs/rand"
	vgtesting "github.com/zeta-protocol/zeta/libs/testing"
	"github.com/zeta-protocol/zeta/paths"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	t.Run("Getting node wallets succeeds", testHandlerGettingNodeWalletsSucceeds)
	t.Run("Getting node wallets with wrong registry passphrase fails", testHandlerGettingNodeWalletsWithWrongRegistryPassphraseFails)
	t.Run("Getting Ethereum wallet succeeds", testHandlerGettingEthereumWalletSucceeds)
	t.Run("Getting Ethereum wallet succeeds", testHandlerGettingEthereumWalletWithWrongRegistryPassphraseFails)
	t.Run("Getting Zeta wallet succeeds", testHandlerGettingZetaWalletSucceeds)
	t.Run("Getting Zeta wallet succeeds", testHandlerGettingZetaWalletWithWrongRegistryPassphraseFails)
	t.Run("Generating Ethereum wallet succeeds", testHandlerGeneratingEthereumWalletSucceeds)
	t.Run("Generating an already existing Ethereum wallet fails", testHandlerGeneratingAlreadyExistingEthereumWalletFails)
	t.Run("Generating Ethereum wallet with overwrite succeeds", testHandlerGeneratingEthereumWalletWithOverwriteSucceeds)
	t.Run("Generating Zeta wallet succeeds", testHandlerGeneratingZetaWalletSucceeds)
	t.Run("Generating an already existing Zeta wallet fails", testHandlerGeneratingAlreadyExistingZetaWalletFails)
	t.Run("Generating Zeta wallet with overwrite succeeds", testHandlerGeneratingZetaWalletWithOverwriteSucceeds)
	t.Run("Importing Ethereum wallet succeeds", testHandlerImportingEthereumWalletSucceeds)
	t.Run("Importing an already existing Ethereum wallet fails", testHandlerImportingAlreadyExistingEthereumWalletFails)
	t.Run("Importing Ethereum wallet with overwrite succeeds", testHandlerImportingEthereumWalletWithOverwriteSucceeds)
	t.Run("Importing Zeta wallet succeeds", testHandlerImportingZetaWalletSucceeds)
	t.Run("Importing an already existing Zeta wallet fails", testHandlerImportingAlreadyExistingZetaWalletFails)
	t.Run("Importing Zeta wallet with overwrite succeeds", testHandlerImportingZetaWalletWithOverwriteSucceeds)
}

func testHandlerGettingNodeWalletsSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)
	config := nodewallets.NewDefaultConfig()

	// setup
	createTestNodeWallets(zetaPaths, registryPass, walletsPass)

	// when
	nw, err := nodewallets.GetNodeWallets(config, zetaPaths, registryPass)

	// assert
	require.NoError(t, err)
	require.NotNil(t, nw)
	require.NotNil(t, nw.Ethereum)
	require.NotNil(t, nw.Zeta)
}

func testHandlerGettingNodeWalletsWithWrongRegistryPassphraseFails(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	wrongRegistryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)
	config := nodewallets.NewDefaultConfig()

	// setup
	createTestNodeWallets(zetaPaths, registryPass, walletsPass)

	// when
	nw, err := nodewallets.GetNodeWallets(config, zetaPaths, wrongRegistryPass)

	// assert
	require.Error(t, err)
	assert.Nil(t, nw)
}

func testHandlerGettingEthereumWalletSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	// setup
	createTestNodeWallets(zetaPaths, registryPass, walletsPass)

	// when
	wallet, err := nodewallets.GetEthereumWallet(zetaPaths, registryPass)

	// assert
	require.NoError(t, err)
	assert.NotNil(t, wallet)
}

func testHandlerGettingEthereumWalletWithWrongRegistryPassphraseFails(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	wrongRegistryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	// setup
	createTestNodeWallets(zetaPaths, registryPass, walletsPass)

	// when
	wallet, err := nodewallets.GetEthereumWallet(zetaPaths, wrongRegistryPass)

	// assert
	require.Error(t, err)
	assert.Nil(t, wallet)
}

func testHandlerGettingZetaWalletSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	// setup
	createTestNodeWallets(zetaPaths, registryPass, walletsPass)

	// when
	wallet, err := nodewallets.GetZetaWallet(zetaPaths, registryPass)

	// then
	require.NoError(t, err)
	assert.NotNil(t, wallet)
}

func testHandlerGettingZetaWalletWithWrongRegistryPassphraseFails(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	wrongRegistryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	// setup
	createTestNodeWallets(zetaPaths, registryPass, walletsPass)

	// when
	wallet, err := nodewallets.GetZetaWallet(zetaPaths, wrongRegistryPass)

	// assert
	require.Error(t, err)
	assert.Nil(t, wallet)
}

func testHandlerGeneratingEthereumWalletSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	data, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletPass, "", false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data["registryFilePath"])
	assert.NotEmpty(t, data["walletFilePath"])
}

func testHandlerGeneratingAlreadyExistingEthereumWalletFails(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass1 := vgrand.RandomStr(10)

	// when
	data1, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletPass1, "", false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data1["registryFilePath"])
	assert.NotEmpty(t, data1["walletFilePath"])

	// given
	walletPass2 := vgrand.RandomStr(10)

	// when
	data2, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletPass2, "", false)

	// then
	require.EqualError(t, err, nodewallets.ErrEthereumWalletAlreadyExists.Error())
	assert.Empty(t, data2)
}

func testHandlerGeneratingEthereumWalletWithOverwriteSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass1 := vgrand.RandomStr(10)

	// when
	data1, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletPass1, "", false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data1["registryFilePath"])
	assert.NotEmpty(t, data1["walletFilePath"])

	// given
	walletPass2 := vgrand.RandomStr(10)

	// when
	data2, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletPass2, "", true)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data2["registryFilePath"])
	assert.Equal(t, data1["registryFilePath"], data2["registryFilePath"])
	assert.NotEmpty(t, data2["walletFilePath"])
	assert.NotEqual(t, data1["walletFilePath"], data2["walletFilePath"])
}

func testHandlerGeneratingZetaWalletSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	data, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass, false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data["registryFilePath"])
	assert.NotEmpty(t, data["walletFilePath"])
	assert.NotEmpty(t, data["mnemonic"])
}

func testHandlerGeneratingAlreadyExistingZetaWalletFails(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass1 := vgrand.RandomStr(10)

	// when
	data1, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass1, false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data1["registryFilePath"])
	assert.NotEmpty(t, data1["walletFilePath"])
	assert.NotEmpty(t, data1["mnemonic"])

	// given
	walletPass2 := vgrand.RandomStr(10)

	// when
	data2, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass2, false)

	// then
	require.EqualError(t, err, nodewallets.ErrZetaWalletAlreadyExists.Error())
	assert.Empty(t, data2)
}

func testHandlerGeneratingZetaWalletWithOverwriteSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass1 := vgrand.RandomStr(10)

	// when
	data1, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass1, false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data1["registryFilePath"])
	assert.NotEmpty(t, data1["walletFilePath"])

	// given
	walletPass2 := vgrand.RandomStr(10)

	// when
	data2, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass2, true)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data2["registryFilePath"])
	assert.Equal(t, data1["registryFilePath"], data2["registryFilePath"])
	assert.NotEmpty(t, data2["walletFilePath"])
	assert.NotEqual(t, data1["walletFilePath"], data2["walletFilePath"])
	assert.NotEmpty(t, data2["mnemonic"])
	assert.NotEqual(t, data1["mnemonic"], data2["mnemonic"])
}

func testHandlerImportingEthereumWalletSucceeds(t *testing.T) {
	// given
	genZetaPaths, genCleanupFn := vgtesting.NewZetaPaths()
	defer genCleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateEthereumWallet(genZetaPaths, registryPass, walletPass, "", false)

	// then
	require.NoError(t, err)

	// given
	importZetaPaths, importCleanupFn := vgtesting.NewZetaPaths()
	defer importCleanupFn()

	// when
	importData, err := nodewallets.ImportEthereumWallet(importZetaPaths, registryPass, walletPass, "", "", genData["walletFilePath"], false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, importData["registryFilePath"])
	assert.NotEqual(t, genData["registryFilePath"], importData["registryFilePath"])
	assert.NotEmpty(t, importData["walletFilePath"])
	assert.NotEqual(t, genData["walletFilePath"], importData["walletFilePath"])
}

func testHandlerImportingAlreadyExistingEthereumWalletFails(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletPass, "", false)

	// then
	require.NoError(t, err)

	// when
	importData, err := nodewallets.ImportEthereumWallet(zetaPaths, registryPass, walletPass, "", genData["walletFilePath"], "", false)

	// then
	require.EqualError(t, err, nodewallets.ErrEthereumWalletAlreadyExists.Error())
	assert.Empty(t, importData)
}

func testHandlerImportingEthereumWalletWithOverwriteSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletPass, "", false)

	// then
	require.NoError(t, err)

	// when
	importData, err := nodewallets.ImportEthereumWallet(zetaPaths, registryPass, walletPass, "", "", genData["walletFilePath"], true)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, genData["registryFilePath"])
	assert.Equal(t, importData["registryFilePath"], genData["registryFilePath"])
	assert.NotEmpty(t, genData["walletFilePath"])
	assert.Equal(t, importData["walletFilePath"], genData["walletFilePath"])
}

func testHandlerImportingZetaWalletSucceeds(t *testing.T) {
	// given
	genZetaPaths, genCleanupFn := vgtesting.NewZetaPaths()
	defer genCleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateZetaWallet(genZetaPaths, registryPass, walletPass, false)

	// then
	require.NoError(t, err)

	// given
	importZetaPaths, importCleanupFn := vgtesting.NewZetaPaths()
	defer importCleanupFn()

	// when
	importData, err := nodewallets.ImportZetaWallet(importZetaPaths, registryPass, walletPass, genData["walletFilePath"], false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, importData["registryFilePath"])
	assert.NotEqual(t, genData["registryFilePath"], importData["registryFilePath"])
	assert.NotEmpty(t, importData["walletFilePath"])
	assert.NotEqual(t, genData["walletFilePath"], importData["walletFilePath"])
}

func testHandlerImportingAlreadyExistingZetaWalletFails(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass, false)

	// then
	require.NoError(t, err)

	// when
	importData, err := nodewallets.ImportZetaWallet(zetaPaths, registryPass, walletPass, genData["walletFilePath"], false)

	// then
	require.EqualError(t, err, nodewallets.ErrZetaWalletAlreadyExists.Error())
	assert.Empty(t, importData)
}

func testHandlerImportingZetaWalletWithOverwriteSucceeds(t *testing.T) {
	// given
	zetaPaths, cleanupFn := vgtesting.NewZetaPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass, false)

	// then
	require.NoError(t, err)

	// when
	importData, err := nodewallets.ImportZetaWallet(zetaPaths, registryPass, walletPass, genData["walletFilePath"], true)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, importData["registryFilePath"])
	assert.Equal(t, genData["registryFilePath"], importData["registryFilePath"])
	assert.NotEmpty(t, importData["walletFilePath"])
	assert.NotEqual(t, genData["walletFilePath"], importData["walletFilePath"])
}

func createTestNodeWallets(zetaPaths paths.Paths, registryPass, walletPass string) {
	if _, err := nodewallets.GenerateEthereumWallet(zetaPaths, registryPass, walletPass, "", false); err != nil {
		panic("couldn't generate Ethereum node wallet for tests")
	}

	if _, err := nodewallets.GenerateZetaWallet(zetaPaths, registryPass, walletPass, false); err != nil {
		panic("couldn't generate Zeta node wallet for tests")
	}
}
