package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.pantheon.tech/stonework/pkg/version"
)

const logo = `
 ███████╗██╗    ██╗ ██████╗████████╗██╗     
 ██╔════╝██║    ██║██╔════╝╚══██╔══╝██║     
 ███████╗██║ █╗ ██║██║        ██║   ██║       %s
 ╚════██║██║███╗██║██║        ██║   ██║       %s
 ███████║╚███╔███╔╝╚██████╗   ██║   ███████╗  %s
 ╚══════╝ ╚══╝╚══╝  ╚═════╝   ╚═╝   ╚══════╝

`

// NewRootCmd returns new root command
func NewRootCmd(cli Cli) *cobra.Command {
	var (
		opts Options
	)
	cmd := &cobra.Command{
		Use:           "swctl [options] [command]",
		Short:         "swctl is CLI app to manage StoneWork and its components",
		Long:          fmt.Sprintf(logo, version.Short(), version.BuildTime(), version.BuiltBy()),
		Version:       version.String(),
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			InitGlobalOptions(cli, &glob)

			return cli.Initialize(opts)
		},
		TraverseChildren:  true,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	}

	cmd.SetIn(cli.In())
	cmd.SetOut(cli.Out())
	cmd.SetErr(cli.Err())

	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false

	opts.InstallFlags(cmd.PersistentFlags())
	glob.InstallFlags(cmd.PersistentFlags())

	cmd.InitDefaultVersionFlag()
	cmd.InitDefaultHelpFlag()
	cmd.Flags().Lookup("help").Hidden = true

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(
		NewConfigCmd(cli),
		NewDeploymentCmd(cli),
		NewStatusCmd(cli),
		NewTraceCmd(cli),
		NewSupportCmd(cli),
	)

	cmd.InitDefaultHelpCmd()
	for _, c := range cmd.Commands() {
		if c.Name() == "help" {
			c.Hidden = true
		}
	}

	return cmd
}

func newVersionCmd() *cobra.Command {
	var (
		short bool
	)
	cmd := cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Run: func(cmd *cobra.Command, args []string) {
			if short {
				fmt.Println(version.String())
			} else {
				fmt.Println(version.Verbose())
			}
		},
		Hidden: true,
	}
	cmd.Flags().BoolVarP(&short, "short", "s", false, "Prints version info in short format")
	return &cmd
}
