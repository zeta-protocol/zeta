package networkhistory

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/grpc/status"

	"github.com/zeta-protocol/zeta/datanode/networkhistory"

	"github.com/zeta-protocol/zeta/datanode/service"
	"github.com/zeta-protocol/zeta/datanode/sqlstore"

	"github.com/zeta-protocol/zeta/datanode/config"
	"github.com/zeta-protocol/zeta/logging"
	"github.com/zeta-protocol/zeta/paths"
)

type fetchCmd struct {
	config.ZetaHomeFlag
	config.Config
}

func (cmd *fetchCmd) Execute(args []string) error {
	ctx, cfunc := context.WithCancel(context.Background())
	defer cfunc()
	cfg := logging.NewDefaultConfig()
	cfg.Custom.Zap.Level = logging.InfoLevel
	cfg.Environment = "custom"
	log := logging.NewLoggerFromConfig(
		cfg,
	)
	defer log.AtExit()

	if len(args) != 2 {
		return errors.New("expected <history-segment-id> <num-blocks-to-fetch>")
	}

	rootSegmentID := args[0]

	numBlocksToFetch, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse number of blocks to fetch: %w", err)
	}

	zetaPaths := paths.New(cmd.ZetaHome)
	err = fixConfig(&cmd.Config, zetaPaths)
	if err != nil {
		return fmt.Errorf("failed to fix config:%w", err)
	}

	err = verifyChainID(cmd.SQLStore.ConnectionConfig, cmd.ChainID)
	if err != nil {
		return fmt.Errorf("failed to verify chain id:%w", err)
	}

	if !datanodeLive(cmd.Config) {
		return fmt.Errorf("datanode must be running for this command to work")
	}

	client := getDatanodeAdminClient(log, cmd.Config)
	blocksFetched, err := networkhistory.FetchHistoryBlocks(ctx, func(s string, args ...interface{}) {
		fmt.Printf(s+"\n", args...)
	}, rootSegmentID,
		func(ctx context.Context, historySegmentId string) (networkhistory.FetchResult, error) {
			resp, err := client.FetchNetworkHistorySegment(ctx, historySegmentId)
			if err != nil {
				return networkhistory.FetchResult{},
					errorFromGrpcError("failed to fetch network history segment", err)
			}

			return networkhistory.FetchResult{
				HeightFrom:               resp.HeightFrom,
				HeightTo:                 resp.HeightTo,
				PreviousHistorySegmentID: resp.PreviousHistorySegmentID,
			}, nil
		}, numBlocksToFetch)
	if err != nil {
		return fmt.Errorf("failed to fetch history blocks:%w", err)
	}

	fmt.Printf("\nfinished fetching history, %d blocks retrieved \n", blocksFetched)

	return nil
}

func verifyChainID(connConfig sqlstore.ConnectionConfig, chainID string) error {
	connSource, err := sqlstore.NewTransactionalConnectionSource(logging.NewTestLogger(), connConfig)
	if err != nil {
		return fmt.Errorf("failed to create new transactional connection source: %w", err)
	}

	defer connSource.Close()

	store := sqlstore.NewChain(connSource)
	chainService := service.NewChain(store)

	err = networkhistory.VerifyChainID(chainID, chainService)
	if err != nil {
		return fmt.Errorf("failed to verify chain id:%w", err)
	}
	return nil
}

func errorFromGrpcError(msg string, err error) error {
	s, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%s:%s", msg, err)
	}

	return fmt.Errorf("%s:%s", msg, s.Details())
}
