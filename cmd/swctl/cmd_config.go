package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: add support for showing only non-internal user config - which excludes stonework-CNF wiring (punts)

type ConfigCmdOptions struct {
	Args []string
}

func NewConfigCmd(cli Cli) *cobra.Command {
	var (
		opts ConfigCmdOptions
	)
	cmd := &cobra.Command{
		Use:                "config [flags] ACTION",
		Short:              "Manage config of StoneWork components",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return runConfigCmd(cli, opts)
		},
	}
	return cmd
}

func runConfigCmd(cli Cli, opts ConfigCmdOptions) error {
	args := opts.Args

	out, err := cli.Exec("agentctl config", args)
	if err != nil {
		return err
	}
	
	fmt.Fprintln(cli.Out(), out)
	return nil
}
