package main

import (
	"os"

	cmd "github.com/zeta-protocol/zeta/cmd/zetawallet/commands"
)

func main() {
	writer := &cmd.Writer{
		Out: os.Stdout,
		Err: os.Stderr,
	}
	cmd.Execute(writer)
}
