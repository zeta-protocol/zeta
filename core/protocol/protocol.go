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

package protocol

import (
	"context"

	"github.com/zeta-protocol/zeta/libs/subscribers"

	"github.com/zeta-protocol/zeta/core/spam"

	"github.com/blang/semver"

	"github.com/zeta-protocol/zeta/core/api"
	"github.com/zeta-protocol/zeta/core/blockchain"
	"github.com/zeta-protocol/zeta/core/broker"
	ethclient "github.com/zeta-protocol/zeta/core/client/eth"
	"github.com/zeta-protocol/zeta/core/config"
	"github.com/zeta-protocol/zeta/core/evtforward"
	"github.com/zeta-protocol/zeta/core/netparams"
	"github.com/zeta-protocol/zeta/core/nodewallets"
	"github.com/zeta-protocol/zeta/core/processor"
	"github.com/zeta-protocol/zeta/core/protocolupgrade"
	"github.com/zeta-protocol/zeta/core/stats"
	"github.com/zeta-protocol/zeta/core/zetatime"
	"github.com/zeta-protocol/zeta/logging"
	"github.com/zeta-protocol/zeta/paths"
)

var Version = semver.MustParse("0.1.0")

type Protocol struct {
	*processor.App

	log *logging.Logger

	confWatcher     *config.Watcher
	confListenerIDs []int

	services *allServices
}

const namedLogger = "protocol"

func New(
	ctx context.Context,
	confWatcher *config.Watcher,
	log *logging.Logger,
	cancel func(),
	stopBlockchain func() error,
	nodewallets *nodewallets.NodeWallets,
	ethClient *ethclient.Client,
	ethConfirmation *ethclient.EthereumConfirmations,
	blockchainClient *blockchain.Client,
	zetaPaths paths.Paths,
	stats *stats.Stats,
) (p *Protocol, err error) {
	log = log.Named(namedLogger)

	defer func() {
		if err != nil {
			log.Error("unable to start protocol", logging.Error(err))
			return
		}

		ids := p.confWatcher.OnConfigUpdateWithID(
			func(cfg config.Config) { p.ReloadConf(cfg.Processor) },
		)
		p.confListenerIDs = ids
	}()

	svcs, err := newServices(
		ctx, log, confWatcher, nodewallets, ethClient, ethConfirmation, blockchainClient, zetaPaths, stats,
	)
	if err != nil {
		return nil, err
	}

	proto := &Protocol{
		App: processor.NewApp(
			log,
			svcs.zetaPaths,
			confWatcher.Get().Processor,
			cancel,
			stopBlockchain,
			svcs.assets,
			svcs.banking,
			svcs.broker,
			svcs.witness,
			svcs.eventForwarder,
			svcs.executionEngine,
			svcs.genesisHandler,
			svcs.governance,
			svcs.notary,
			svcs.stats.Blockchain,
			svcs.timeService,
			svcs.epochService,
			svcs.topology,
			svcs.netParams,
			&processor.Oracle{
				Engine:   svcs.oracle,
				Adaptors: svcs.oracleAdaptors,
			},
			svcs.delegation,
			svcs.limits,
			svcs.stakeVerifier,
			svcs.checkpoint,
			svcs.spam,
			svcs.pow,
			svcs.stakingAccounts,
			svcs.snapshot,
			svcs.statevar,
			svcs.blockchainClient,
			svcs.erc20MultiSigTopology,
			stats.GetVersion(),
			svcs.protocolUpgradeEngine,
			svcs.codec,
			svcs.gastimator,
		),
		log:         log,
		confWatcher: confWatcher,
		services:    svcs,
	}

	proto.services.netParams.Watch(
		netparams.WatchParam{
			Param:   netparams.SpamProtectionMaxBatchSize,
			Watcher: proto.App.OnSpamProtectionMaxBatchSizeUpdate,
		},
	)

	return proto, nil
}

// Start will start the protocol, this means it's ready to process
// blocks from the blockchain.
func (n *Protocol) Start() error {
	return nil
}

// Stop will stop all services of the protocol.
func (n *Protocol) Stop() error {
	// unregister conf listeners
	n.log.Info("Stopping protocol services")
	n.confWatcher.Unregister(n.confListenerIDs)
	n.services.Stop()
	return nil
}

func (n *Protocol) Protocol() semver.Version {
	return Version
}

func (n *Protocol) GetEventForwarder() *evtforward.Forwarder {
	return n.services.eventForwarder
}

func (n *Protocol) GetTimeService() *zetatime.Svc {
	return n.services.timeService
}

func (n *Protocol) GetEventService() *subscribers.Service {
	return n.services.eventService
}

func (n *Protocol) GetBroker() *broker.Broker {
	return n.services.broker
}

func (n *Protocol) GetPoW() api.ProofOfWorkParams {
	return n.services.pow
}

func (n *Protocol) GetProtocolUpgradeService() *protocolupgrade.Engine {
	return n.services.protocolUpgradeEngine
}

func (n *Protocol) GetSpamEngine() *spam.Engine {
	return n.services.spam
}

func (n *Protocol) GetPowEngine() processor.PoWEngine {
	return n.services.pow
}
