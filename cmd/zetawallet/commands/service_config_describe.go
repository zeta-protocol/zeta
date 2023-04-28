package cmd

import (
	"fmt"
	"io"

	"github.com/zeta-protocol/zeta/cmd/zetawallet/commands/cli"
	"github.com/zeta-protocol/zeta/cmd/zetawallet/commands/flags"
	"github.com/zeta-protocol/zeta/cmd/zetawallet/commands/printer"
	"github.com/zeta-protocol/zeta/paths"
	"github.com/zeta-protocol/zeta/wallet/service"
	svcStoreV1 "github.com/zeta-protocol/zeta/wallet/service/store/v1"

	"github.com/spf13/cobra"
)

var (
	describeServiceConfigLong = cli.LongDesc(`
	    Describe the service configuration.
	`)

	describeServiceConfigExample = cli.Examples(`
		# Describe the service configuration
		{{.Software}} service config describe
	`)
)

type DescribeServiceConfigHandler func() (*service.Config, error)

func NewCmdDescribeServiceConfig(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func() (*service.Config, error) {
		zetaPaths := paths.New(rf.Home)

		svcStore, err := svcStoreV1.InitialiseStore(zetaPaths)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise service store: %w", err)
		}

		cfg, err := svcStore.GetConfig()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve the service configuration: %w", err)
		}

		return cfg, nil
	}

	return BuildCmdDescribeServiceConfig(w, h, rf)
}

func BuildCmdDescribeServiceConfig(w io.Writer, handler DescribeServiceConfigHandler, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Describe the service configuration",
		Long:    describeServiceConfigLong,
		Example: describeServiceConfigExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := handler()
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintDescribeServiceConfigResponse(w, cfg)
			case flags.JSONOutput:
				return printer.FprintJSON(w, cfg)
			}

			return nil
		},
	}

	return cmd
}

func PrintDescribeServiceConfigResponse(w io.Writer, cfg *service.Config) {
	p := printer.NewInteractivePrinter(w)

	str := p.String()
	defer p.Print(str)

	str.NextLine()
	str.Text("Service URL: ").WarningText(cfg.Server.String()).NextSection()
	str.Text("Log level: ").WarningText(cfg.LogLevel.String()).NextSection()
	str.Text("API V1").NextLine()
	str.Pad().Text("Maximum token duration: ").WarningText(cfg.APIV1.MaximumTokenDuration.String()).NextLine()
}
