package app

import (
	"errors"
	"fmt"
	"regexp"
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

  <white># Set quantity of runtime HugePages manually</>
  $ <yellow>` + appName + ` dependency set-hugepages <value></>

  <white># Assign(up) or Unassign(down) interfaces to/from kernel</>
  $ <yellow>` + appName + ` dependency link <interfaces ...> up | down</>

  <white># Print out startup config with dpdk interfaces</>
  $ <yellow>` + appName + ` dependency get-startup [<interfaceName:StoneworkInterfaceName ...>]</>

  <white># Print out startup config with dpdk plugin disable</>
  $ <yellow>` + appName + ` dependency get-startup</>
`
}

type NetworkInterface struct {
	Name        string
	Pci         string
	Description string
	LinkUp      bool
	SwName      string
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
			var err error
			docker, err := IsDockerAvailable(cli)
			if err != nil {
				return err
			}
			hugePages, err := AllocatedHugePages(cli)
			if err != nil {
				return err

			}
			physicalInterfaces, err := DumpNetworkInterfaces(cli)
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
				table.SetHeader([]string{"Name", "Pci", "Mode", "State"})

				for _, n := range physicalInterfaces {
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

			docker, err := IsDockerAvailable(cli)
			if err != nil {
				return err
			}

			if !docker {
				err = InstallDocker(cli, "default")
				if err != nil {
					return err
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

/* DPDK interface cannot be used by kernel otherwise it won't connect to VPP*/
func linkSetUpDown(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link ",
		Short: "link <interfaces ...> up | down",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 2 {
				physicalInterfaces, err := DumpNetworkInterfaces(cli)
				if err != nil {
					return err
				}
				if strings.Compare(args[(len(args)-1)], "up") == 0 {
					for _, arg := range args[:(len(args) - 1)] {
						matchId := -1
						for i, physicalInterface := range physicalInterfaces {
							if physicalInterface.Name == arg {
								matchId = i
								break
							}
						}
						if matchId == -1 {
							return errors.New("Interface: " + arg + "does not exist.")
						}
						// TODO get PCI from name func
						out, _, err := cli.Exec("sudo dpdk-devbind.py -u "+physicalInterfaces[matchId].Pci, nil)
						if err != nil {
							return err
						}
						fmt.Fprintln(cli.Out(), out)
						// TODO get PCI from name func
						out, _, err = cli.Exec("sudo ip link set "+physicalInterfaces[matchId].Name+" up", nil)
						if err != nil {
							return err
						}

					}
				} else if strings.Compare(args[(len(args)-1)], "down") == 0 {
					for _, arg := range args[:(len(args) - 1)] {
						matchId := -1
						for i, physicalInterface := range physicalInterfaces {
							if physicalInterface.Name == arg {
								matchId = i
								break
							}
						}
						if matchId == -1 {
							return errors.New("Interface: " + arg + "does not exist.")
						}
						out, _, err := cli.Exec("sudo ip link set "+physicalInterfaces[matchId].Name+" down", nil)
						if err != nil {
							return err
						}
						fmt.Fprintln(cli.Out(), out)
						// TODO get PCI from name func
						out, _, err = cli.Exec("dpdk-devbind.py -b uio_pci_generic "+physicalInterfaces[matchId].Pci, nil)
						if err != nil {
							return err
						}

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
	if hugePgSize == 0 {
		return 0, err
	}

	return hugePgSize, err
}

func ResizeHugePages(cli Cli, size uint) error {
	//TODO: Make persistent hugepages
	//TODO: Handle numa case, Big hugepages(are immutable and can be setted only during booting)
	if size == 0 {
		fmt.Fprintln(cli.Out(), "Skipping hugepages")
		return nil
	}
	_, _, err := cli.Exec(fmt.Sprintf("sudo sysctl -w vm.nr_hugepages=%d", size), nil)
	if err != nil {
		return err
	}
	allocatedHP, _ := AllocatedHugePages(cli)
	if size != uint(allocatedHP) {
		return errors.New("Failed to allocate enough continuous memory")
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
		"sudo apt install dpdk -y",
		"echo \"uio_pci_generic\" | sudo tee -a /etc/modules",
		//"sudo modprobe uio_pci_generic",
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
		fmt.Println(out)

	}

	return nil
}

// Dump only physical interfaces
func DumpNetworkInterfaces(cli Cli) ([]NetworkInterface, error) {
	//path leads to networking devices in the OS
	const path = "/sys/class/net"
	var allDevices []NetworkInterface
	var physicalDevices []NetworkInterface

	cmd := "ls -b"
	out, _, err := cli.Exec(cmd, []string{path})
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Command: " + cmd + " finished with error " + err.Error())
	}

	for _, name := range strings.Fields(out) {
		allDevices = append(allDevices, NetworkInterface{Name: name})
	}

	for _, nic := range allDevices {
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

			physicalDevices = append(physicalDevices, newNic)
		}

	}

	return physicalDevices, nil
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
			physicalInterfaces, err := DumpNetworkInterfaces(cli)
			if err != nil {
				return err
			}

			for i, desiredInterface := range desiredInterfaces {
				for _, dumpedInterface := range physicalInterfaces {
					if desiredInterface.Name == dumpedInterface.Name {
						desiredInterfaces[i].Pci = dumpedInterface.Pci
						break
					}
					return errors.New("Requested interface " + desiredInterface.Name + " does not exist")
				}
			}

			t := template.Must(template.New("startupConf").Parse(startupconfig))
			err = t.Execute(cli.Out(), desiredInterfaces)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
