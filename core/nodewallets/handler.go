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

package nodewallets

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/zeta-protocol/zeta/core/nodewallets/registry"
	"github.com/zeta-protocol/zeta/core/nodewallets/zeta"
	"github.com/zeta-protocol/zeta/paths"
)

var (
	ErrEthereumWalletAlreadyExists   = errors.New("the Ethereum node wallet already exists")
	ErrZetaWalletAlreadyExists       = errors.New("the Zeta node wallet already exists")
	ErrTendermintPubkeyAlreadyExists = errors.New("the Tendermint pubkey already exists")
)

func GetZetaWallet(zetaPaths paths.Paths, registryPassphrase string) (*zeta.Wallet, error) {
	registryLoader, err := registry.NewLoader(zetaPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	registry, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if registry.Zeta == nil {
		return nil, ErrZetaWalletIsMissing
	}

	walletLoader, err := zeta.InitialiseWalletLoader(zetaPaths)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise Zeta node wallet loader: %w", err)
	}

	wallet, err := walletLoader.Load(registry.Zeta.Name, registry.Zeta.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load Ethereum node wallet: %w", err)
	}

	return wallet, nil
}

func GetNodeWallets(config Config, zetaPaths paths.Paths, registryPassphrase string) (*NodeWallets, error) {
	nodeWallets := &NodeWallets{}

	registryLoader, err := registry.NewLoader(zetaPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	reg, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if reg.Ethereum != nil {
		w, err := GetEthereumWalletWithRegistry(zetaPaths, reg)
		if err != nil {
			return nil, err
		}

		nodeWallets.Ethereum = w
	}

	if reg.Zeta != nil {
		zetaWalletLoader, err := zeta.InitialiseWalletLoader(zetaPaths)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise Zeta node wallet loader: %w", err)
		}

		nodeWallets.Zeta, err = zetaWalletLoader.Load(reg.Zeta.Name, reg.Zeta.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("couldn't load Zeta node wallet: %w", err)
		}
	}

	if reg.Tendermint != nil {
		nodeWallets.Tendermint = &TendermintPubkey{
			Pubkey: reg.Tendermint.Pubkey,
		}
	}

	return nodeWallets, nil
}

func GenerateZetaWallet(zetaPaths paths.Paths, registryPassphrase, walletPassphrase string, overwrite bool) (map[string]string, error) {
	registryLoader, err := registry.NewLoader(zetaPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	reg, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if !overwrite && reg.Zeta != nil {
		return nil, ErrZetaWalletAlreadyExists
	}

	zetaWalletLoader, err := zeta.InitialiseWalletLoader(zetaPaths)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise Zeta node wallet loader: %w", err)
	}

	w, data, err := zetaWalletLoader.Generate(walletPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate Zeta node wallet: %w", err)
	}

	reg.Zeta = &registry.RegisteredZetaWallet{
		Name:       w.Name(),
		Passphrase: walletPassphrase,
	}

	if err := registryLoader.Save(reg, registryPassphrase); err != nil {
		return nil, fmt.Errorf("couldn't save registry: %w", err)
	}

	data["registryFilePath"] = registryLoader.RegistryFilePath()
	return data, nil
}

func ImportZetaWallet(zetaPaths paths.Paths, registryPassphrase, walletPassphrase, sourceFilePath string, overwrite bool) (map[string]string, error) {
	if !filepath.IsAbs(sourceFilePath) {
		return nil, fmt.Errorf("path to the wallet file need to be absolute")
	}

	registryLoader, err := registry.NewLoader(zetaPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	reg, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if !overwrite && reg.Zeta != nil {
		return nil, ErrZetaWalletAlreadyExists
	}

	zetaWalletLoader, err := zeta.InitialiseWalletLoader(zetaPaths)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise Zeta node wallet loader: %w", err)
	}

	w, data, err := zetaWalletLoader.Import(sourceFilePath, walletPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Zeta node wallet: %w", err)
	}

	reg.Zeta = &registry.RegisteredZetaWallet{
		Name:       w.Name(),
		Passphrase: walletPassphrase,
	}

	if err := registryLoader.Save(reg, registryPassphrase); err != nil {
		return nil, fmt.Errorf("couldn't save registry: %w", err)
	}

	data["registryFilePath"] = registryLoader.RegistryFilePath()
	return data, nil
}

func ImportTendermintPubkey(
	zetaPaths paths.Paths,
	registryPassphrase, pubkey string,
	overwrite bool,
) (map[string]string, error) {
	registryLoader, err := registry.NewLoader(zetaPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	reg, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if !overwrite && reg.Tendermint != nil {
		return nil, ErrTendermintPubkeyAlreadyExists
	}

	reg.Tendermint = &registry.RegisteredTendermintPubkey{
		Pubkey: pubkey,
	}

	if err := registryLoader.Save(reg, registryPassphrase); err != nil {
		return nil, fmt.Errorf("couldn't save registry: %w", err)
	}

	return map[string]string{
		"registryFilePath": registryLoader.RegistryFilePath(),
		"tendermintPubkey": pubkey,
	}, nil
}
