package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"

	"go.pantheon.tech/stonework/plugins/puntmgr"
)

// TODO:
// - decide when to run agentctl inside docker container or outside of docker container

type ConfigCmdOptions struct {
	ShowInternal bool
	Args         []string
}

func (opts *ConfigCmdOptions) InstallFlags(flagset *pflag.FlagSet) {
	flagset.BoolVar(&opts.ShowInternal, "show-internal", false, "Add Stonework internal configuration to output if possible")
}

func NewConfigCmd(cli Cli) *cobra.Command {
	var (
		opts ConfigCmdOptions
	)
	cmd := &cobra.Command{
		Use:   "config [flags] ACTION",
		Short: "Manage config of StoneWork components",
		Args:  cobra.ArbitraryArgs,
		// DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return runConfigCmd(cli, opts)
		},
	}
	opts.InstallFlags(cmd.PersistentFlags())
	return cmd
}

func runConfigCmd(cli Cli, opts ConfigCmdOptions) error {
	args := opts.Args

	if slices.Contains(args, "get") && !opts.ShowInternal {
		hideInternalFlag := "--labels=\"!" + puntmgr.InternalConfigLabelKey + "=" + puntmgr.InternalConfigLabelValue + "\""
		args = append(args, hideInternalFlag)
	}

	stdout, stderr, err := cli.Exec("agentctl config", args, false)
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.Err(), stderr)
	fmt.Fprintln(cli.Out(), stdout)
	return nil
}
