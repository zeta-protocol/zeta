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

package core_test

import (
	"context"
	"fmt"
	"os"

	"github.com/zeta-protocol/zeta/core/banking"
	"github.com/zeta-protocol/zeta/core/collateral"
	"github.com/zeta-protocol/zeta/core/delegation"
	"github.com/zeta-protocol/zeta/core/epochtime"
	"github.com/zeta-protocol/zeta/core/evtforward"
	"github.com/zeta-protocol/zeta/core/execution"
	"github.com/zeta-protocol/zeta/core/notary"
	"github.com/zeta-protocol/zeta/core/rewards"
	"github.com/zeta-protocol/zeta/core/validators"
	"github.com/zeta-protocol/zeta/libs/num"

	"github.com/zeta-protocol/zeta/core/integration/helpers"
	"github.com/zeta-protocol/zeta/core/integration/steps/market"
	"github.com/zeta-protocol/zeta/core/integration/stubs"
	"github.com/zeta-protocol/zeta/core/netparams"
	"github.com/zeta-protocol/zeta/core/oracles"
	"github.com/zeta-protocol/zeta/core/plugins"
	"github.com/zeta-protocol/zeta/core/types"
	"github.com/zeta-protocol/zeta/logging"
	protos "github.com/zeta-protocol/zeta/protos/zeta"

	"github.com/golang/mock/gomock"
)

var (
	execsetup *executionTestSetup
	reporter  tstReporter
)

type tstReporter struct {
	scenario string
}

func (t tstReporter) Errorf(format string, args ...interface{}) {
	fmt.Printf("%s ERROR: %s", t.scenario, fmt.Sprintf(format, args...))
}

func (t tstReporter) Fatalf(format string, args ...interface{}) {
	fmt.Printf("%s FATAL: %s", t.scenario, fmt.Sprintf(format, args...))
	os.Exit(1)
}

var marketConfig = market.NewMarketConfig()

type executionTestSetup struct {
	cfg              execution.Config
	log              *logging.Logger
	ctrl             *gomock.Controller
	timeService      *stubs.TimeStub
	broker           *stubs.BrokerStub
	executionEngine  *exEng
	collateralEngine *collateral.Engine
	oracleEngine     *oracles.Engine
	builtinOracle    *oracles.Builtin
	epochEngine      *epochtime.Svc
	delegationEngine *delegation.Engine
	positionPlugin   *plugins.Positions
	topology         *stubs.TopologyStub
	stakingAccount   *stubs.StakingAccountStub
	rewardsEngine    *rewards.Engine
	assetsEngine     *stubs.AssetStub
	banking          *banking.Engine

	// save party accounts state
	markets []types.Market

	block *helpers.Block

	netParams *netparams.Store

	// keep track of net deposits/withdrawals (ignores asset type)
	netDeposits *num.Uint

	// record parts of state before each step
	accountsBefore                []protos.Account
	ledgerMovementsBefore         int
	insurancePoolDepositsOverStep map[string]*num.Int
	eventsBefore                  int

	ntry           *notary.SnapshotNotary
	stateVarEngine *stubs.StateVarStub
	witness        *validators.Witness
}

func newExecutionTestSetup() *executionTestSetup {
	if execsetup != nil && execsetup.ctrl != nil {
		execsetup.ctrl.Finish()
	} else if execsetup == nil {
		execsetup = &executionTestSetup{}
	}

	ctrl := gomock.NewController(&reporter)
	execsetup.ctrl = ctrl
	execsetup.cfg = execution.NewDefaultConfig()
	execsetup.cfg.Position.StreamPositionVerbose = true
	execsetup.cfg.Risk.StreamMarginLevelsVerbose = true
	execsetup.log = logging.NewTestLogger()
	execsetup.timeService = stubs.NewTimeStub()
	execsetup.broker = stubs.NewBrokerStub()
	execsetup.collateralEngine = collateral.New(
		execsetup.log, collateral.NewDefaultConfig(), execsetup.timeService, execsetup.broker,
	)

	zetaAsset := types.Asset{
		ID: "ZETA",
		Details: &types.AssetDetails{
			Name:   "ZETA",
			Symbol: "ZETA",
		},
	}
	execsetup.collateralEngine.EnableAsset(context.Background(), zetaAsset)

	usdt := types.Asset{
		ID: "USDT",
		Details: &types.AssetDetails{
			Name:   "USDT",
			Symbol: "USDT",
		},
	}
	execsetup.collateralEngine.EnableAsset(context.Background(), usdt)

	usdc := types.Asset{
		ID: "USDC",
		Details: &types.AssetDetails{
			Name:   "USDC",
			Symbol: "USDC",
		},
	}

	execsetup.collateralEngine.EnableAsset(context.Background(), usdc)

	execsetup.epochEngine = epochtime.NewService(execsetup.log, epochtime.NewDefaultConfig(), execsetup.broker)

	execsetup.topology = stubs.NewTopologyStub("nodeID", execsetup.broker)

	execsetup.stakingAccount = stubs.NewStakingAccountStub()
	execsetup.epochEngine.NotifyOnEpoch(execsetup.stakingAccount.OnEpochEvent, execsetup.stakingAccount.OnEpochRestore)

	marketActivityTracker := execution.NewMarketActivityTracker(execsetup.log, execsetup.epochEngine)
	commander := stubs.NewCommanderStub()
	execsetup.netDeposits = num.UintZero()
	execsetup.witness = validators.NewWitness(execsetup.log, validators.NewDefaultConfig(), execsetup.topology, commander, execsetup.timeService)

	execsetup.ntry = notary.NewWithSnapshot(execsetup.log, notary.NewDefaultConfig(), execsetup.topology, execsetup.broker, commander)
	execsetup.assetsEngine = stubs.NewAssetStub()
	ethSourceNoop := evtforward.NewNoopEngine(execsetup.log, evtforward.NewDefaultConfig())
	execsetup.banking = banking.New(execsetup.log, banking.NewDefaultConfig(), execsetup.collateralEngine, execsetup.witness, execsetup.timeService, execsetup.assetsEngine, execsetup.ntry, execsetup.broker, execsetup.topology, execsetup.epochEngine, marketActivityTracker, stubs.NewBridgeViewStub(), ethSourceNoop)

	execsetup.delegationEngine = delegation.New(execsetup.log, delegation.NewDefaultConfig(), execsetup.broker, execsetup.topology, execsetup.stakingAccount, execsetup.epochEngine, execsetup.timeService)
	execsetup.rewardsEngine = rewards.New(execsetup.log, rewards.NewDefaultConfig(), execsetup.broker, execsetup.delegationEngine, execsetup.epochEngine, execsetup.collateralEngine, execsetup.timeService, marketActivityTracker, execsetup.topology)

	execsetup.oracleEngine = oracles.NewEngine(
		execsetup.log, oracles.NewDefaultConfig(), execsetup.timeService, execsetup.broker,
	)

	execsetup.builtinOracle = oracles.NewBuiltinOracle(execsetup.oracleEngine, execsetup.timeService)

	execsetup.stateVarEngine = stubs.NewStateVar()
	// @TODO stub assets engine and pass it in

	execsetup.executionEngine = newExEng(
		execution.NewEngine(
			execsetup.log,
			execsetup.cfg,
			execsetup.timeService,
			execsetup.collateralEngine,
			execsetup.oracleEngine,
			execsetup.broker,
			execsetup.stateVarEngine,
			marketActivityTracker,
			execsetup.assetsEngine, // assets
		),
		execsetup.broker,
	)

	execsetup.positionPlugin = plugins.NewPositions(context.Background())
	execsetup.broker.Subscribe(execsetup.positionPlugin)

	execsetup.block = helpers.NewBlock()

	execsetup.registerTimeServiceCallbacks()

	execsetup.netParams = netparams.New(execsetup.log, netparams.NewDefaultConfig(), execsetup.broker)
	if err := execsetup.registerNetParamsCallbacks(); err != nil {
		panic(err)
	}

	execsetup.netParams.Watch(
		netparams.WatchParam{
			Param:   netparams.FloatingPointUpdatesDuration,
			Watcher: execsetup.stateVarEngine.OnFloatingPointUpdatesDurationUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.TransferFeeFactor,
			Watcher: execsetup.banking.OnTransferFeeFactorUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.TransferMinTransferQuantumMultiple,
			Watcher: execsetup.banking.OnMinTransferQuantumMultiple,
		},
		netparams.WatchParam{
			Param:   netparams.MaxPeggedOrders,
			Watcher: execsetup.executionEngine.OnMaxPeggedOrderUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarkPriceUpdateMaximumFrequency,
			Watcher: execsetup.executionEngine.OnMarkPriceUpdateMaximumFrequency,
		},
	)
	return execsetup
}

func (e *executionTestSetup) registerTimeServiceCallbacks() {
	e.timeService.NotifyOnTick(
		e.epochEngine.OnTick,
		e.witness.OnTick,
		e.ntry.OnTick,
		e.banking.OnTick,
		e.delegationEngine.OnTick,
		e.builtinOracle.OnTick,
		e.stateVarEngine.OnTick,
		e.executionEngine.OnTick,
	)
}

func (e *executionTestSetup) registerNetParamsCallbacks() error {
	return e.netParams.Watch(
		netparams.WatchParam{
			Param:   netparams.StakingAndDelegationRewardMinimumValidatorStake,
			Watcher: e.topology.OnMinDelegationUpdated,
		},
		netparams.WatchParam{
			Param:   netparams.MarketMarginScalingFactors,
			Watcher: e.executionEngine.OnMarketMarginScalingFactorsUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketFeeFactorsMakerFee,
			Watcher: e.executionEngine.OnMarketFeeFactorsMakerFeeUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketFeeFactorsInfrastructureFee,
			Watcher: e.executionEngine.OnMarketFeeFactorsInfrastructureFeeUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketLiquidityStakeToCCYVolume,
			Watcher: e.executionEngine.OnSuppliedStakeToObligationFactorUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketValueWindowLength,
			Watcher: e.executionEngine.OnMarketValueWindowLengthUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketLiquidityProvidersFeeDistribitionTimeStep,
			Watcher: e.executionEngine.OnMarketLiquidityProvidersFeeDistributionTimeStep,
		},
		netparams.WatchParam{
			Param:   netparams.MarketLiquidityProvisionShapesMaxSize,
			Watcher: e.executionEngine.OnMarketLiquidityProvisionShapesMaxSizeUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketLiquidityMaximumLiquidityFeeFactorLevel,
			Watcher: e.executionEngine.OnMarketLiquidityMaximumLiquidityFeeFactorLevelUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketLiquidityBondPenaltyParameter,
			Watcher: e.executionEngine.OnMarketLiquidityBondPenaltyUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketAuctionMinimumDuration,
			Watcher: e.executionEngine.OnMarketAuctionMinimumDurationUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketProbabilityOfTradingTauScaling,
			Watcher: e.executionEngine.OnMarketProbabilityOfTradingTauScalingUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.DelegationMinAmount,
			Watcher: e.delegationEngine.OnMinAmountChanged,
		},
		netparams.WatchParam{
			Param:   netparams.StakingAndDelegationRewardMaxPayoutPerParticipant,
			Watcher: e.rewardsEngine.UpdateMaxPayoutPerParticipantForStakingRewardScheme,
		},
		netparams.WatchParam{
			Param:   netparams.StakingAndDelegationRewardDelegatorShare,
			Watcher: e.rewardsEngine.UpdateDelegatorShareForStakingRewardScheme,
		},
		netparams.WatchParam{
			Param:   netparams.StakingAndDelegationRewardMinimumValidatorStake,
			Watcher: e.rewardsEngine.UpdateMinimumValidatorStakeForStakingRewardScheme,
		},
		netparams.WatchParam{
			Param:   netparams.RewardAsset,
			Watcher: e.rewardsEngine.UpdateAssetForStakingAndDelegation,
		},
		netparams.WatchParam{
			Param:   netparams.StakingAndDelegationRewardCompetitionLevel,
			Watcher: e.rewardsEngine.UpdateCompetitionLevelForStakingRewardScheme,
		},
		netparams.WatchParam{
			Param:   netparams.StakingAndDelegationRewardsMinValidators,
			Watcher: e.rewardsEngine.UpdateMinValidatorsStakingRewardScheme,
		},
		netparams.WatchParam{
			Param:   netparams.StakingAndDelegationRewardOptimalStakeMultiplier,
			Watcher: e.rewardsEngine.UpdateOptimalStakeMultiplierStakingRewardScheme,
		},
		netparams.WatchParam{
			Param:   netparams.ValidatorsEpochLength,
			Watcher: e.epochEngine.OnEpochLengthUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketMinLpStakeQuantumMultiple,
			Watcher: e.executionEngine.OnMinLpStakeQuantumMultipleUpdate,
		},
		netparams.WatchParam{
			Param:   netparams.MarketMinProbabilityOfTradingForLPOrders,
			Watcher: e.executionEngine.OnMarketMinProbabilityOfTradingForLPOrdersUpdate,
		},
	)
}
