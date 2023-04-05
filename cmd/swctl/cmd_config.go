package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: add support for showing only non-internal user config - which excludes stonework-CNF wiring (punts)

func NewConfigCmd(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "config [flags] ACTION",
		Short:              "Manage config of StoneWork components",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := cli.Exec("agentctl config", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), out)
			return nil
		},
	}
	cmd.AddCommand(
		NewConfigGetCmd(cli),
		NewConfigSetCmd(cli),
	)
	return cmd
}

func NewConfigGetCmd(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "get [flags]",
		Short:              "Retrieve configuration from StoneWork",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := cli.Exec("agentctl config get", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), out)
			return nil
		},
	}
	return cmd
}

func NewConfigSetCmd(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "set [flags]",
		Aliases:            []string{"s", "update", "u"},
		Short:              "Update configuration in StoneWork",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := cli.Exec("agentctl config get", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), out)
			return nil
		},
	}
	return cmd
}
