package app

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func deploymentExample(appName string) string {
	return `
  <white># Create and start a deployment</>
  $ <yellow>` + appName + ` deployment up</>

  <white># Stop and remove a deployment</>
  $ <yellow>` + appName + ` deployment down</>

  <white># Show info about deployment</>
  $ <yellow>` + appName + ` deployment info</>

  <white># Print configuration of deployment</>
  $ <yellow>` + appName + ` deployment config</>

  <white># List images in deployment</>
  $ <yellow>` + appName + ` deployment images</>

  <white># List services in deployment</>
  $ <yellow>` + appName + ` deployment services</>
`
}

func NewDeploymentCmd(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployment [flags] COMMAND",
		Short: "Manage deployments of StoneWork",
		Long: `Command deploy manages deployment of StoneWork.
It hides the specifics of a deployment and provides unified interface to user. 
Currently only Docker Compose is supported, for which the deploy command offers a wrapper.`,
		Example: color.Sprint(deploymentExample(cli.AppName())),
		Aliases: []string{"deploy"},
		Args:    cobra.ArbitraryArgs,
		//DisableFlagParsing: true,
	}
	cmd.AddCommand(
		newDeploymentUp(cli),
		newDeploymentDown(cli),
		newDeploymentConfig(cli),
		newDeploymentInfo(cli),
		newDeploymentImages(cli),
		newDeploymentServices(cli),
	)
	return cmd
}

func newDeploymentUp(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "up [flags]",
		Short:              "Create and start deployment",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout, stderr, err := cli.Exec("docker compose up", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), stdout)
			fmt.Fprintln(cli.Err(), stderr)
			return nil
		},
	}
	return cmd
}

func newDeploymentDown(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "down [flags]",
		Short:              "Stop and remove deployment",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout, stderr, err := cli.Exec("docker compose down", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), stdout)
			fmt.Fprintln(cli.Err(), stderr)
			return nil
		},
	}
	return cmd
}

func newDeploymentConfig(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "config [flags]",
		Short:              "Show deployment config",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout, stderr, err := cli.Exec("docker compose convert", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), stdout)
			fmt.Fprintln(cli.Err(), stderr)
			return nil
		},
	}
	return cmd
}

func newDeploymentInfo(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "info [flags]",
		Short:              "Show info about deployment",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout, stderr, err := cli.Exec("docker compose ps", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), stdout)
			fmt.Fprintln(cli.Err(), stderr)
			return nil
		},
	}
	return cmd
}

func newDeploymentImages(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "images [flags]",
		Short:              "List images in deployment",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout, stderr, err := cli.Exec("docker compose images", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), stdout)
			fmt.Fprintln(cli.Err(), stderr)
			return nil
		},
	}
	return cmd
}

func newDeploymentServices(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "services [flags]",
		Short:              "List services in deployment",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout, stderr, err := cli.Exec("docker compose ps --services", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), stdout)
			fmt.Fprintln(cli.Err(), stderr)
			return nil
		},
	}
	return cmd
}
