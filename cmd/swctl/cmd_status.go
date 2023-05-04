package main

import (
	"fmt"
	"os/exec"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

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
		Example: color.Sprint(statusExample),
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

	// TODO: improve status overview, show status of components (CNFs)

	out, err := cli.Exec("vpp-probe --env=docker discover", opts.Args)
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%v: %s", ee.String(), ee.Stderr)
		}
		return err
	}

	fmt.Fprintln(cli.Out(), out)

	return nil
}
