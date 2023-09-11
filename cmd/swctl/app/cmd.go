package app

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"go.pantheon.tech/stonework/pkg/version"
)

const logo = `
<lightblue>  ███████╗██╗    ██╗ ██████╗████████╗██╗      </>
<lightblue>  ██╔════╝██║    ██║██╔════╝╚══██╔══╝██║      </>
<lightblue>  ███████╗██║ █╗ ██║██║        ██║   ██║      </><lightyellow> %s </>
<lightblue>  ╚════██║██║███╗██║██║        ██║   ██║      </><lightyellow> %s </>
<lightblue>  ███████║╚███╔███╔╝╚██████╗   ██║   ███████╗ </><lightyellow> %s </>
<lightblue>  ╚══════╝ ╚══╝╚══╝  ╚═════╝   ╚═╝   ╚══════╝ </>

`

// NewRootCmd returns new root command
func NewRootCmd(cli Cli) *cobra.Command {
	var (
		glob GlobalOptions
	)

	cmd := &cobra.Command{
		Use:           "swctl [options] [command]",
		Short:         "swctl is CLI tool to manage StoneWork and its components",
		Long:          color.Sprintf(logo, version.Short(), version.BuildTime(), version.BuiltBy()),
		Version:       version.String(),
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cli.Initialize(&glob)
			if err != nil {
				return err
			}
			logrus.Tracef("initialized global options: %+v", glob)
			return nil
		},
		TraverseChildren:  true,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	}

	cmd.SetIn(cli.In())
	cmd.SetOut(cli.Out())
	cmd.SetErr(cli.Err())

	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false

	glob.InstallFlags(cmd.PersistentFlags())

	cmd.InitDefaultVersionFlag()
	cmd.InitDefaultHelpFlag()
	cmd.Flags().Lookup("help").Hidden = true

	cmd.AddCommand(NewVersionCmd())
	cmd.AddCommand(
		NewConfigCmd(cli),
		NewDeploymentCmd(cli),
		NewStatusCmd(cli),
		NewTraceCmd(cli),
		NewSupportCmd(cli),
		NewManageCmd(cli),
		NewDependencyCmd(cli),
	)

	cmd.InitDefaultHelpCmd()
	for _, c := range cmd.Commands() {
		if c.Name() == "help" {
			c.Hidden = true
		}
	}

	return cmd
}

func NewVersionCmd() *cobra.Command {
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
