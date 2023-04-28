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

package admin

import (
	"os"
	"path"

	"github.com/zeta-protocol/zeta/libs/config/encoding"
	"github.com/zeta-protocol/zeta/logging"
)

const (
	nvServerNamedLogger   = "admin.non-validator-server"
	vServerNamedLogger    = "admin.validator-server"
	clientNamedLogger     = "admin-client"
	nodeWalletNamedLogger = "admin.node-wallet"
)

// ServerConfig represent the configuration of the server.
type ServerConfig struct {
	SocketPath string `long:"socket-path" description:"Listen for connection on UNIX socket path <file-path>"`
	HTTPPath   string `long:"http-path" description:"Http path of the socket HTTP RPC server"`
}

// Config represents the configuration of the admin package.
type Config struct {
	Level  encoding.LogLevel `long:"log-level"`
	Server ServerConfig      `group:"Server" namespace:"server"`
}

// NewDefaultConfig creates an instance of the package specific configuration, given a
// pointer to a logger instance to be used for logging within the package.
func NewDefaultConfig() Config {
	return Config{
		Level: encoding.LogLevel{Level: logging.InfoLevel},
		Server: ServerConfig{
			SocketPath: path.Join(os.TempDir(), "zeta.sock"),
			HTTPPath:   "/rpc",
		},
	}
}
