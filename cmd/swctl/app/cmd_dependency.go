package app

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/exp/slices"

	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func exampleDependencyCmd(appName string) string {
	return `
  <white># Status of all dependencies</>
  $ <yellow>` + appName + ` dependency status</>

  <white># Install external tools (docker, docker compose)</>
  $ <yellow>` + appName + ` dependency install-ext </>

  <white># Set HugePages in kB manually to size 2048kB</>
  $ <yellow>` + appName + ` dependency hugepages <value></>

  <white># Assign(up) or Unassign(down) interfaces to/from kernel</>
  $ <yellow>` + appName + ` dependency link <interfaces ...> up | down</>

  <white># Print out startup config with dpdk interfaces</>
  $ <yellow>` + appName + ` dependency startup [<interfaceName:StoneworkInterfaceName ...>]</>

  <white># Print out startup config with dpdk plugin disable</>
  $ <yellow>` + appName + ` dependency startup</>
`
}

type NetworkInterface struct {
	Name        string
	Pci         string
	Description string
	LinkUp      bool
	SwName      string
}

type Dependencies struct {
	Docker     bool
	HugePages  int
	Interfaces []NetworkInterface
}

func NewDependencyCmd(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dependency COMMAND",
		Short:         "Manage external dependencies",
		Example:       color.Sprint(exampleDependencyCmd(cli.AppName())),
		Args:          nil,
		SilenceErrors: true,
		SilenceUsage:  true,
		// overriding Root's PersistentPreRunE because in any dependency
		// commands is not needed docker client connection
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	var glob = GlobalOptions{
		Debug:              false,
		LogLevel:           "",
		Color:              "",
		ComposeFiles:       nil,
		EntityFiles:        nil,
		EmbeddedEntityByte: nil,
	}
	cli.Initialize(&glob)

	cmd.AddCommand(installExternalTools(cli), dependecyStatus(cli), installHugePages(cli), linkSetUpDown(cli), startupConf(cli))

	return cmd
}

func dependecyStatus(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "status",
		Short:         "status",
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			dpdcs := &Dependencies{}
			dpdcs.Docker = dpdcs.IsDockerAvailable(cli)
			_, dpdcs.HugePages = dpdcs.IsHugePagesEnabled(cli)
			dpdcs.Interfaces = dpdcs.DumpNetworkInterfaces(cli)
			type statusInfo struct {
			}
			var status string
			if dpdcs.Docker {
				status = "OK"
			} else {
				status = "Not installed"
			}
			fmt.Fprintf(cli.Out(), "Docker: %s\n", status)

			if dpdcs.HugePages == 0 {
				status = "Disabled"
			} else {
				status = strconv.Itoa(dpdcs.HugePages)
			}
			fmt.Fprintf(cli.Out(), "Hugepages: %s\n", status)

			if dpdcs.Interfaces == nil {
				status = "No available interfaces\n"
				fmt.Fprintf(cli.Out(), status)
			} else {
				table := tablewriter.NewWriter(cli.Out())
				table.SetHeader([]string{"Name", "Pci", "Mode", "State"})

				for _, n := range dpdcs.Interfaces {
					row := []string{n.Name, n.Pci, n.Description}
					if n.LinkUp == true {
						row = append(row, "LinkUp\n")
					} else {
						row = append(row, "LinkDown\n")
					}
					table.Append(row)
				}
				table.Render()
			}
			return nil
		},
	}
	return cmd
}

func installExternalTools(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-tools",
		Short: "Install external tools",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dpdcs := &Dependencies{}

			dpdcs.Docker = dpdcs.IsDockerAvailable(cli)

			if !dpdcs.Docker {
				err := dpdcs.InstallDocker(cli, "default")
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("Docker is already installed")

			return nil
		},
	}
	return cmd
}

func installHugePages(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hugepages ",
		Short: "hugepages <value>",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var dep Dependencies
			size, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			err = dep.ResizeHugePages(cli, uint(size))
			if err != nil {
				panic(err)
			}
			return nil

		},
	}
	return cmd
}

/* DPDK interface cannot be used by kernel otherwise it won't connect to VPP*/
func linkSetUpDown(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link ",
		Short: "link <interfaces ...> up | down",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 2 {
				if strings.Compare(args[(len(args)-1)], "up") == 0 {
					for _, arg := range args[:(len(args) - 1)] {
						out, _, err := cli.Exec("sudo ip link set "+arg+" up", nil)
						if err != nil {
							return err
						}
						fmt.Fprintln(cli.Out(), out)
					}
				} else if strings.Compare(args[(len(args)-1)], "down") == 0 {
					for _, arg := range args[:(len(args) - 1)] {
						out, _, err := cli.Exec("sudo ip link set "+arg+" down", nil)
						if err != nil {
							return err
						}
						fmt.Fprintln(cli.Out(), out)
					}

				} else {
					return errors.New("Last argument must define operation up or down upon selected interfaces")
				}
			} else {
				return errors.New("Command must consist of two or more arguments")
			}

			return nil
		},
	}
	return cmd
}

func (*Dependencies) IsDockerAvailable(cli Cli) bool {
	out, _, err := cli.Exec("whereis docker", nil)
	if err != nil {
		panic(err)
	}
	if strings.Contains(out, "/docker") {
		return true
	}
	return false
}

func (*Dependencies) IsHugePagesEnabled(cli Cli) (bool, int) {
	out, _, err := cli.Exec("sysctl vm.nr_hugepages -n", nil)
	if err != nil {
		panic(err)
	}
	hugePgSize, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		fmt.Println(err)
	}
	if hugePgSize == 0 {
		return false, hugePgSize
	}

	return true, hugePgSize
}

func (*Dependencies) ResizeHugePages(cli Cli, size uint) error {
	//TODO: Handle numa case, Big hugepages(are immutable and can be setted only during booting)
	if size == 0 {
		fmt.Fprintln(cli.Out(), "Skipping hugepages")
		return nil
	}
	_, _, err := cli.Exec(fmt.Sprintf("sudo sysctl -w vm.nr_hugepages=%d", size), nil)
	if err != nil {
		return err
	}
	return nil
}

func (*Dependencies) InstallDocker(cli Cli, dockerVersion string) error {

	commands := []string{"sudo apt-get update -y",
		"sudo apt-get install ca-certificates curl gnupg -y",
		"sudo install -m 0755 -d /etc/apt/keyrings",
		"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg --yes",
		"sudo chmod a+r /etc/apt/keyrings/docker.gpg",
		`echo \
		"deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
		"$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
		sudo tee /etc/apt/sources.list.d/docker.list > /dev/null`,
		"sudo apt-get update -y",
	}
	if dockerVersion == "default" {
		cmd := `sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin`
		commands = append(commands, cmd)
	} else {
		cmd := `sudo apt-get install -y docker-ce=` + dockerVersion + ` docker-ce-cli=` + dockerVersion + ` containerd.io docker-buildx-plugin docker-compose-plugin`
		commands = append(commands, cmd)

	}

	for _, command := range commands {
		out, _, err := cli.Exec("bash -c", []string{command})
		if err != nil {
			return err
		}
		fmt.Println(out)

	}

	return nil
}

// Dump only physical interfaces
func (*Dependencies) DumpNetworkInterfaces(cli Cli) []NetworkInterface {
	const path = "/sys/class/net"
	var list []NetworkInterface
	var realDevices []NetworkInterface

	out, _, err := cli.Exec("ls -b", []string{path})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	buff := strings.Fields(out)

	for _, name := range buff {
		list = append(list, NetworkInterface{Name: name})
	}

	for _, nic := range list {
		_, _, err := cli.Exec("ls ", []string{path + "/" + nic.Name})
		if err == nil {
			newNic := NetworkInterface{Name: nic.Name}

			info, _, _ := cli.Exec("cat", []string{path + "/" + nic.Name + "/device/uevent"})

			pciRe := regexp.MustCompile(`PCI_SLOT_NAME=(\S+)`)
			match := pciRe.FindStringSubmatch(info)
			if len(match) == 0 {
				continue
			}
			newNic.Pci = match[1]

			driverRe := regexp.MustCompile(`DRIVER=(\S+)`)
			match = driverRe.FindStringSubmatch(info)
			newNic.Description = match[1]

			_, stderr, _ := cli.Exec("cat", []string{path + "/" + nic.Name + "/carrier"})
			if stderr != "" {
				newNic.LinkUp = false
			} else {
				newNic.LinkUp = true
			}

			realDevices = append(realDevices, newNic)
		}

	}

	return realDevices
}

func startupConf(cli Cli) *cobra.Command {
	const startupconfig = `unix {
cli-no-pager
cli-listen /run/vpp/cli.sock
log /tmp/vpp.log
coredump-size unlimited
full-coredump
poll-sleep-usec 50
}
{{if .}}
dpdk {
{{range .}}  dev {{.Pci}} {
    name: {{.SwName}}
}
{{end}} 
}
{{else}}
plugins {
     plugin dpdk_plugin.so { disable }
}
{{end}}
api-trace {
    on
}

socksvr {
	default
}

statseg {
	default
	per-node-counters on
}

punt {
    socket /run/stonework/vpp/punt-to-vpp.sock
}
`
	cmd := &cobra.Command{
		Use:   "startup",
		Short: "Print out startup config",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			var desiredInterfaces []NetworkInterface
			dpdcs := &Dependencies{}
			for _, arg := range args {

				var netInterface NetworkInterface
				names := strings.Split(arg, ":")
				if len(names) != 2 {
					return errors.New("Bad format of argument. Every argument in this command" +
						" must have \"word:word\" pattern")
				}
				netInterface.Name = names[0]
				netInterface.SwName = names[1]

				desiredInterfaces = append(desiredInterfaces, netInterface)

			}
			dpdcs.Interfaces = dpdcs.DumpNetworkInterfaces(cli)

			for i, desiredInterface := range desiredInterfaces {
				for _, dumpedInterface := range dpdcs.Interfaces {
					if desiredInterface.Name == dumpedInterface.Name {
						desiredInterfaces[i].Pci = dumpedInterface.Pci
						break
					}
					return errors.New("Requested interface " + desiredInterface.Name + " does not exist")
				}
			}

			t := template.Must(template.New("startupConf").Parse(startupconfig))
			err := t.Execute(cli.Out(), desiredInterfaces)
			if err != nil {
				fmt.Println("Could not execute template")
			}
			return nil
		},
	}
	return cmd
}
func StartupConfManualInterfaces(cli Cli, interfaces []string) {
	const startupconfig = `unix {
cli-no-pager
cli-listen /run/vpp/cli.sock
log /tmp/vpp.log
coredump-size unlimited
full-coredump
poll-sleep-usec 50
}
{{if .}}
dpdk {
{{range .}}	dev {{.}}
{{end}} 
}
{{else}}
plugins {
     plugin dpdk_plugin.so { disable }
}
{{end}}
api-trace {
    on
}

socksvr {
	default
}

statseg {
	default
	per-node-counters on
}

punt {
    socket /run/stonework/vpp/punt-to-vpp.sock
}
`

	dpdcs := &Dependencies{}
	dpdcs.Interfaces = dpdcs.DumpNetworkInterfaces(cli)

	pcis := []string{}
	for _, intfc := range dpdcs.Interfaces {
		if ok := slices.Contains[string](interfaces, intfc.Name); ok {
			pcis = append(pcis, intfc.Pci)

		}
	}

	t := template.Must(template.New("startupConf").Parse(startupconfig))
	err := t.Execute(cli.Out(), pcis)
	if err != nil {
		fmt.Println("Could not execute template")
	}

}
