package cmd

import (
	"fmt"
	"io"

	"github.com/zeta-protocol/zeta/cmd/zetawallet/commands/printer"
	"github.com/zeta-protocol/zeta/wallet/wallet"
	"github.com/spf13/cobra"
)

func NewCmdKey(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Manage Zeta wallets' keys",
		Long:  "Manage Zeta wallets' keys",
	}

	cmd.AddCommand(NewCmdAnnotateKey(w, rf))
	cmd.AddCommand(NewCmdGenerateKey(w, rf))
	cmd.AddCommand(NewCmdIsolateKey(w, rf))
	cmd.AddCommand(NewCmdListKeys(w, rf))
	cmd.AddCommand(NewCmdDescribeKey(w, rf))
	cmd.AddCommand(NewCmdTaintKey(w, rf))
	cmd.AddCommand(NewCmdUntaintKey(w, rf))
	cmd.AddCommand(NewCmdRotateKey(w, rf))
	return cmd
}

func printMeta(str *printer.FormattedString, meta []wallet.Metadata) {
	padding := 0
	for _, m := range meta {
		keyLen := len(m.Key)
		if keyLen > padding {
			padding = keyLen
		}
	}

	for _, m := range meta {
		str.WarningText(fmt.Sprintf("%-*s", padding, m.Key)).Text(" | ").WarningText(m.Value).NextLine()
	}
}
