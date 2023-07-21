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
	args := append([]string{"--print"}, opts.Args...)

	// TODO: improve usage of trace command
	// - automatically use ping as default command
	// - consider selecting IP and source network namespace automatically
	// - consider allowing users to simply select component names for ping src/dst

	stdout, stderr, err := cli.Exec(fmt.Sprintf("vpp-probe --env=%s trace", defaultVppProbeEnv), args)
	if err != nil {
		return err
	}
	fmt.Fprintln(cli.Out(), stdout)
	fmt.Fprintln(cli.Err(), stderr)

	return nil
}
