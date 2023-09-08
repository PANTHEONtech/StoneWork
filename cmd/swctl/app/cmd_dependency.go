package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func exampleDependencyCmd(appName string) string {
	return `
  <white># Status of all dependencies</>
  $ <yellow>` + appName + ` dependency status</>

  <white># Install external tools (docker, docker compose)</>
  $ <yellow>` + appName + ` dependency install-tools</>

  <white># Set quantity of runtime 2MB HugePages manually</>
  $ <yellow>` + appName + ` dependency set-hugepages <value></>

  <white># Assign(up) or Unassign(down) interfaces to/from kernel</>
  $ <yellow>` + appName + ` dependency link <pci ...> up | down</>

  <white># Print out startup config with dpdk interfaces</>
  $ <yellow>` + appName + ` dependency get-startup [<interfacePci:StoneworkInterfaceName ...>]</>

  <white># Print out startup config with dpdk plugin disable</>
  $ <yellow>` + appName + ` dependency get-startup</>
`
}

type NetworkInterface struct {
	Name        string
	Pci         string
	Description string
	SwName      string
	Module      string
	// Nil Driver means that device is unbounded and can be used by vpp which choose driver
	Driver string
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

	cmd.AddCommand(installExternalTools(cli), dependencyStatus(cli), installHugePages(cli), linkSetUpDown(cli), startupConf(cli))

	return cmd
}

func dependencyStatus(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "status",
		Short:         "status",
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			docker, err := IsDockerAvailable(cli)
			if err != nil {
				return err
			}
			hugePages, err := AllocatedHugePages(cli)
			if err != nil {
				return err

			}
			physicalInterfaces, err := DumpDevices(cli)
			if err != nil {
				return err
			}
			var status string
			if docker {
				status = "OK"
			} else {
				status = "Not installed"
			}
			fmt.Fprintf(cli.Out(), "Docker: %s\n", status)

			if hugePages == 0 {
				status = "Disabled"
			} else {
				status = strconv.Itoa(hugePages)
			}
			fmt.Fprintf(cli.Out(), "Hugepages: %s\n", status)

			if physicalInterfaces == nil {
				status = "No available interfaces\n"
				fmt.Fprintf(cli.Out(), status)
			} else {
				table := tablewriter.NewWriter(cli.Out())
				table.SetHeader([]string{"Name", "Pci", "Mode", "Driver"})

				for _, n := range physicalInterfaces {
					row := []string{n.Name, n.Pci}
					if n.Driver == n.Module {
						row = append(row, "Kernel")
					} else {
						row = append(row, "DPDK")
					}
					row = append(row, n.Driver)

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

			docker, err := IsDockerAvailable(cli)
			if err != nil {
				return errors.New(fmt.Sprintf("Unable to check docker availability: %v", err))
			}

			if docker {
				fmt.Fprintln(cli.Out(), "Docker is already installed")
				return nil
			}

			err = InstallDocker(cli, "default")
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

func installHugePages(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-hugepages",
		Short: "set-hugepages <value>",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			size, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			if err = ResizeHugePages(cli, uint(size)); err != nil {
				return err
			}
			return nil

		},
	}
	return cmd
}

/*
linkSetUpDown changes the link state and binds/unbinds the pci driver
(up=kernel driver usable for kernel network stack, down=no kernel driver).
The kernel driver unbinding is helpfull in case of DPDK interfaces that
can't have kernel network stack usable driver when VPP should use them
as DPDK interfaces.
*/
func linkSetUpDown(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link ",
		Short: "link < pci ...> up | down",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("command must consist of two or more arguments")
			}

			if !(args[len(args)-1] == "up" || args[len(args)-1] == "down") {
				return errors.New("last argument must define operation up or down upon selected interfaces")
			}

			physicalInterfaces, err := DumpDevices(cli)
			if err != nil {
				return err
			}

			for _, arg := range args[:(len(args) - 1)] {
				matchId := -1
				for i, physicalInterface := range physicalInterfaces {
					if physicalInterface.Pci == arg {
						matchId = i
						break
					}
				}
				if matchId == -1 {
					return errors.New("Interface: " + arg + "does not exist.")
				}
				if args[len(args)-1] == "up" {
					//returning interface back to kernel driver
					if physicalInterfaces[matchId].Driver == "" {
						fmt.Fprintln(cli.Out(), fmt.Sprintf("don't need to unbind the already unbinded pci %s", physicalInterfaces[matchId].Name))
					} else {
						err = unbindDevice(cli, physicalInterfaces[matchId].Pci, physicalInterfaces[matchId].Driver)
						if err != nil {
							return err
						}
					}

					err = bindDevice(cli, physicalInterfaces[matchId].Pci, physicalInterfaces[matchId].Module)

					if err != nil {
						return err
					}
				} else if args[len(args)-1] == "down" {
					//link down interface, only assigned network devices have /net directory which is name of interface
					stdout, stderr, err := cli.Exec("ls", []string{"/sys/bus/pci/devices/" + physicalInterfaces[matchId].Pci + "/net"})
					if stderr != "" {
						return errors.New(stderr)
					}
					if err != err {
						return err
					}
					if stdout != "" {
						_, _, err = cli.Exec("sudo ip link set "+stdout+" down", nil)
						if err != err {
							return err
						}
					}

					err = unbindDevice(cli, physicalInterfaces[matchId].Pci, physicalInterfaces[matchId].Driver)
					if err != nil {
						return err
					}
					// binding is not needed as VPP will bind interface automatically
				}

			}

			return nil
		},
	}
	return cmd
}

func IsDockerAvailable(cli Cli) (bool, error) {
	out, _, err := cli.Exec("whereis docker", nil)
	if err != nil {
		return false, err
	}
	if strings.Contains(out, "/docker") {
		return true, nil
	}
	return false, nil
}

func AllocatedHugePages(cli Cli) (int, error) {
	out, _, err := cli.Exec("sysctl vm.nr_hugepages -n", nil)
	if err != nil {
		return 0, err
	}
	hugePgSize, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, err
	}

	return hugePgSize, nil
}

func ResizeHugePages(cli Cli, size uint) error {
	const hugePageSize = 2048
	//TODO: Make persistent hugepages
	//TODO: Handle numa case, Big (1GB)hugepages(are immutable and can be setted only during booting)
	if size == 0 {
		fmt.Fprintln(cli.Out(), "Skipping hugepages")
		return nil
	}
	_, _, err := cli.Exec(fmt.Sprintf("sudo sysctl -w vm.nr_hugepages=%d", size), nil)
	if err != nil {
		return err
	}
	allocatedHP, err := AllocatedHugePages(cli)
	if err != nil {
		return err
	}

	if size != uint(allocatedHP) {
		return errors.New(fmt.Sprintf("failed to allocate enough hugepages (%s),successfully allocated %s hugepages, totally continuous memory %d MB ",
			size, allocatedHP, (allocatedHP*hugePageSize)/1000))
	}

	return nil
}

func InstallDocker(cli Cli, dockerVersion string) error {

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
		"echo \"uio_pci_generic\" | sudo tee -a /etc/modules",
	}
	if dockerVersion == "default" {
		cmd := `sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin`
		commands = append(commands, cmd)
	} else {
		cmd := `sudo apt-get install -y docker-ce=` + dockerVersion + ` docker-ce-cli=` + dockerVersion + ` containerd.io docker-buildx-plugin docker-compose-plugin`
		commands = append(commands, cmd)

	}

	for _, command := range commands {
		out, stderr, err := cli.Exec("bash -c", []string{command})
		if stderr != "" {
			return errors.New(command + ": " + stderr)
		}
		if err != nil {
			return errors.New(err.Error() + "(" + command + ")")
		}
		fmt.Fprintln(cli.Out(), out)

	}

	return nil
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
    name {{.SwName}}
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
		Use:   "get-startup",
		Short: "Print out startup config",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			var desiredInterfaces []NetworkInterface
			for _, arg := range args {

				var netInterface NetworkInterface
				trimIndex := strings.LastIndex(arg, ":")
				names := []string{arg[:trimIndex], arg[trimIndex+1:]}
				if len(names) != 2 {
					return errors.New("bad format of argument. Every argument in this command" +
						" must have \"word:word\" pattern")
				}
				netInterface.Pci = names[0]
				netInterface.SwName = names[1]

				desiredInterfaces = append(desiredInterfaces, netInterface)

			}

			t := template.Must(template.New("startupConf").Parse(startupconfig))
			err := t.Execute(cli.Out(), desiredInterfaces)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func unbindDevice(cli Cli, pci string, driver string) error {
	// dpdk drivers like uio_pci_generic, vfio-pci etc..
	// kernel drivers like e1000, ...
	//Mostly
	path := fmt.Sprintf("/sys/bus/pci/drivers/%s/unbind", driver)

	_, stderr, err := cli.Exec("sudo bash -c", []string{"echo \"" + pci + "\" > " + path})
	if stderr != "" {
		return errors.New(stderr)
	}
	if err != nil {
		return err
	}
	return nil
}
func bindDevice(cli Cli, pci string, driver string) error {

	path := fmt.Sprintf("/sys/bus/pci/drivers/%s/bind", driver)

	_, stderr, err := cli.Exec("sudo bash -c", []string{"echo \"" + pci + "\" > " + path})
	if stderr != "" {
		return errors.New(stderr)
	}
	if err != nil {
		return err
	}
	return nil
}

func DumpDevices(cli Cli) ([]NetworkInterface, error) {
	var nics []NetworkInterface

	stdout, _, err := cli.Exec("lspci", []string{"-Dvmmnnk"})
	if err != nil {
		return nil, err
	}
	devicesStr := strings.Split(stdout, "\n\n")
	var networkDevices []map[string]string
	for _, deviceStr := range devicesStr {
		device := make(map[string]string)
		attributes := strings.Split(deviceStr, "\n")
		// parse Slot,Class, Module,Driver,Device
		for _, attribute := range attributes {
			tokenized := strings.Split(attribute, "\t")
			device[strings.Trim(tokenized[0], ":")] = tokenized[1]
		}
		// Class 0200 is code determined for ethernet, we are not interested in other devices
		if strings.Contains(device["Class"], "0200") {
			networkDevices = append(networkDevices, device)
		}

	}
	for _, networkDevice := range networkDevices {
		nic := NetworkInterface{
			// Full name of interface
			Name: networkDevice["Device"],
			Pci:  networkDevice["Slot"],
			// Default kernel driver
			Module: networkDevice["Module"],
			// Actually used driver
			Driver: networkDevice["Driver"],
		}
		nics = append(nics, nic)
	}
	return nics, nil

}
