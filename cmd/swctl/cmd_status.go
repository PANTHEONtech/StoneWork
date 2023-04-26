package main

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// TODO: improve status overview, show status of components (CNFs)
//   - instead of using raw output from vpp-probe, retrieve the important info
//     about the running/deployed components of StoneWork and show those by default
//   - optionally allow user to set more details which shows the more detailed output
//     similar to vpp-probe discover

const statusExample = `
  <white># Show status for all components</>
  $ <yellow>swctl status</>
`

type StatusCmdOptions struct {
	Args []string
}

func NewStatusCmd(cli Cli) *cobra.Command {
	var opts StatusCmdOptions
	cmd := &cobra.Command{
		Use:     "status [flags]",
		Short:   "Show status of StoneWork components",
		Example: statusExample,
		Args:    cobra.ArbitraryArgs,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return runStatusCmd(cli, opts)
		},
	}
	return cmd
}

func runStatusCmd(cli Cli, opts StatusCmdOptions) error {
	cmd := fmt.Sprintf("vpp-probe --env=%q discover", defaultVppProbeEnv)
	out, err := cli.Exec(cmd, opts.Args)
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%v: %s", ee.String(), ee.Stderr)
		}
		return err
	}

	fmt.Fprintln(cli.Out(), out)
	return nil
}
