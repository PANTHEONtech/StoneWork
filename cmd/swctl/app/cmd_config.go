package app

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
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
		Use:                "config [flags] ACTION",
		Short:              "Manage config of StoneWork components",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return runConfigCmd(cli, opts, cmd)
		},
	}
	opts.InstallFlags(cmd.PersistentFlags())
	return cmd
}

func runConfigCmd(cli Cli, opts ConfigCmdOptions, cmd *cobra.Command) error {
	args := opts.Args

	if slices.Contains(args, "--help") || slices.Contains(args, "-h") {
		// add Local flags to stdout of agentctl
		stdout, stderr, shouldReturn, returnValue := mergeHelpers(cmd, cli, args)
		if shouldReturn {
			return returnValue
		}

		fmt.Fprintln(cli.Err(), stderr)
		fmt.Fprintln(cli.Out(), stdout)
		return nil
	}

	if slices.Contains(args, "--show-internal") {
		opts.ShowInternal = true
		idx := slices.Index[string](args, "--show-internal")
		args = append(args[:idx], args[idx+1:]...)
	}

	if slices.Contains(args, "get") && !opts.ShowInternal {
		hideInternalFlag := "--labels=\"!" + puntmgr.InternalConfigLabelKey + "=" + puntmgr.InternalConfigLabelValue + "\""
		args = append(args, hideInternalFlag)
	}

	stdout, stderr, err := cli.Exec("agentctl config", args, false)
	if err != nil {
		return err
	}

	color.Fprintln(cli.Err(), stderr)
	color.Fprintln(cli.Out(), stdout)
	return nil
}

func mergeHelpers(cmd *cobra.Command, cli Cli, args []string) (string, string, bool, error) {
	flags := cmd.LocalFlags()
	bufb := flags.FlagUsages()

	stdout, stderr, err := cli.Exec("agentctl config", args, false)
	if err != nil {
		return "", "", true, err
	}
	stdout = strings.ReplaceAll(stdout, "agentctl", "swctl")

	globalsIndex := strings.Index(stdout, "GLOBALS:")
	if globalsIndex != -1 {
		stdout = stdout[:globalsIndex] + fmt.Sprintf("Flags:\n%s", bufb) + stdout[globalsIndex:]
	}
	return stdout, stderr, false, nil
}
