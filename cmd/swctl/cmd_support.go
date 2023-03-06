package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type SupportCmdOptions struct {
}

func NewSupportCmd(cli Cli) *cobra.Command {
	var opts SupportCmdOptions
	cmd := &cobra.Command{
		Use:                "support [flags]",
		Short:              "Export support data",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSupportCmd(cli, opts, args)
		},
	}
	return cmd
}

func runSupportCmd(cli Cli, opts SupportCmdOptions, args []string) error {

	// TODO: add stonework/CNF related support data to the export

	out, err := cli.Exec("agentctl report", args)
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.Out(), out)

	return nil
}
