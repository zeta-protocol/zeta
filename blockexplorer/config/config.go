// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package config

import (
	"fmt"

	"github.com/zeta-protocol/zeta/blockexplorer/api"
	"github.com/zeta-protocol/zeta/blockexplorer/store"
	vgfs "github.com/zeta-protocol/zeta/libs/fs"
	"github.com/zeta-protocol/zeta/logging"
	"github.com/zeta-protocol/zeta/paths"
)

type ZetaHomeFlag struct {
	ZetaHome string `long:"home" description:"Path to the custom home for zeta"`
}

type Config struct {
	API     api.Config
	Store   store.Config
	Logging logging.Config `namespace:"logging" group:"logging"`
}

func NewDefaultConfig() Config {
	return Config{
		API:     api.NewDefaultConfig(),
		Store:   store.NewDefaultConfig(),
		Logging: logging.NewDefaultConfig(),
	}
}

type Loader struct {
	configFilePath string
}

func NewLoader(zetaPaths paths.Paths) (*Loader, error) {
	configFilePath, err := zetaPaths.CreateConfigPathFor(paths.BlockExplorerDefaultConfigFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't get path for %s: %w", paths.NodeDefaultConfigFile, err)
	}

	return &Loader{
		configFilePath: configFilePath,
	}, nil
}

func (l *Loader) ConfigFilePath() string {
	return l.configFilePath
}

func (l *Loader) ConfigExists() (bool, error) {
	exists, err := vgfs.FileExists(l.configFilePath)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (l *Loader) Save(cfg *Config) error {
	if err := paths.WriteStructuredFile(l.configFilePath, cfg); err != nil {
		return fmt.Errorf("couldn't write configuration file at %s: %w", l.configFilePath, err)
	}
	return nil
}

func (l *Loader) Get() (*Config, error) {
	cfg := NewDefaultConfig()
	if err := paths.ReadStructuredFile(l.configFilePath, &cfg); err != nil {
		return nil, fmt.Errorf("couldn't read configuration file at %s: %w", l.configFilePath, err)
	}
	return &cfg, nil
}
