package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type TraceCmdOptions struct {
	Args []string
}

func NewTraceCmd(cli Cli) *cobra.Command {
	var opts TraceCmdOptions
	cmd := &cobra.Command{
		Use:                "trace [flags]",
		Short:              "Trace packets across data path",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return runTraceCmd(cli, opts)
		},
	}
	return cmd
}

func runTraceCmd(cli Cli, opts TraceCmdOptions) error {

	// TODO: improve tracing usage

	args := append([]string{"--print"}, opts.Args...)

	out, err := cli.Exec("vpp-probe --env=docker trace", args)
	if err != nil {
		return err
	}
	fmt.Fprintln(cli.Out(), out)

	return nil
}
