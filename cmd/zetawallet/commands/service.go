package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdService(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the Zeta wallet's service",
		Long:  "Manage the Zeta wallet's service",
	}

	cmd.AddCommand(NewCmdRunService(w, rf))
	cmd.AddCommand(NewCmdListEndpoints(w, rf))
	cmd.AddCommand(NewCmdServiceConfig(w, rf))
	return cmd
}
