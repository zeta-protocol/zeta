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

package evtforward

import (
	"time"

	"github.com/zeta-protocol/zeta/core/evtforward/ethereum"
	"github.com/zeta-protocol/zeta/libs/config/encoding"
	"github.com/zeta-protocol/zeta/logging"
)

const (
	forwarderLogger = "forwarder"
	// how often the Event Forwarder needs to select a node to send the event
	// if nothing was received.
	defaultRetryRate = 10 * time.Second
)

// Config represents governance specific configuration.
type Config struct {
	// Level specifies the logging level of the Event Forwarder engine.
	Level     encoding.LogLevel `long:"log-level"`
	RetryRate encoding.Duration `long:"retry-rate"`
	// a list of allowlisted blockchain queue public keys
	BlockchainQueueAllowlist []string `long:"blockchain-queue-allowlist" description:" "`
	// Ethereum groups the configuration related to Ethereum implementation of
	// the Event Forwarder.
	Ethereum ethereum.Config `group:"Ethereum" namespace:"ethereum"`
}

// NewDefaultConfig creates an instance of the package specific configuration.
func NewDefaultConfig() Config {
	return Config{
		Level:                    encoding.LogLevel{Level: logging.InfoLevel},
		RetryRate:                encoding.Duration{Duration: defaultRetryRate},
		BlockchainQueueAllowlist: []string{},
		Ethereum:                 ethereum.NewDefaultConfig(),
	}
}
