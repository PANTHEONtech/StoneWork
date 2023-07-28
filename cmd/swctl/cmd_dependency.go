package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

const exampleDependencyCmd = `
  <white># List of all dependencies</>
  $ <yellow>swctl dependency status</>

  <white># Install all dependencies</>
  $ <yellow>swctl dependency install</>

  <white># Set HugePages manually</>
  $ <yellow>swctl dependency hugepages <value></>

  <white># Create multiple entity instances</>
  $ <yellow>swctl dependency linkdown <interfaces ...></>

  <white># Merge with existing config file</>
  $ <yellow>swctl dependency startup </>

`
const (
	DefaultHugepageSize = 1024
)

type networkInterface struct {
	name        string
	pci         string
	description string
	linkUp      bool
}

type dependencies struct {
	docker     bool
	hugePages  int
	interfaces []networkInterface
}

func NewInstallCmd(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dependency  COMMAND",
		Short:   "install",
		Example: color.Sprint(exampleDependencyCmd),
		Args:    nil,
	}
	cmd.AddCommand(installAll(cli), dependecyStatus(cli), installHugePages(cli), linkSetDown(cli), startupConf(cli))

	return cmd
}

func dependecyStatus(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "status",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dpdcs := &dependencies{}
			dpdcs.docker = dpdcs.isDockerAvailable(cli)
			_, dpdcs.hugePages = dpdcs.isHugePagesEnabled(cli)
			dpdcs.interfaces = dpdcs.dumpNetworkInterfaces(cli)
			var status string
			if dpdcs.docker {
				status = "OK"
			} else {
				status = "Not installed"
			}
			fmt.Fprintf(cli.Out(), "Docker: %s\n", status)

			if dpdcs.hugePages == 0 {
				status = "Disabled"
			} else {
				status = strconv.Itoa(dpdcs.hugePages)
			}
			fmt.Fprintf(cli.Out(), "Hugepages: %s\n", status)

			if dpdcs.interfaces == nil {
				status = "No available interfaces"
			} else {
				status = ""
				for _, n := range dpdcs.interfaces {
					status = status + fmt.Sprintf("%s\t%s\t%s\t", n.name, n.pci, n.description)
					if n.linkUp == true {
						status = status + fmt.Sprintf("LinkUp\n")
					} else {
						status = status + fmt.Sprintf("LinkDown\n")
					}
				}
			}
			fmt.Fprintf(cli.Out(), "Interfaces:\n %s\n", status)

			return nil
		},
	}
	return cmd
}

func installAll(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dpdcs := &dependencies{}
			dpdcs.docker = dpdcs.isDockerAvailable(cli)
			_, dpdcs.hugePages = dpdcs.isHugePagesEnabled(cli)
			dpdcs.interfaces = dpdcs.dumpNetworkInterfaces(cli)

			if !dpdcs.docker {
				err := dpdcs.installDocker(cli)
				if err != nil {
					panic(err)
				}
			}
			err := dpdcs.resizeHugePages(cli, uint(DefaultHugepageSize))
			if err != nil {
				panic(err)
			}

			return nil
		},
	}
	return cmd
}

func installDocker(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-docker",
		Short: "Install dependencies",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, _, err := cli.Exec("whereis docker", args)
			if err != nil {
				return err
			}
			fmt.Fprintln(cli.Out(), out)
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
			var dep dependencies
			size, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			err = dep.resizeHugePages(cli, uint(size))
			if err != nil {
				panic(err)
			}
			return nil

		},
	}
	return cmd
}

/* DPDK interface cannot be used by kernel otherwise it won't connect to VPP*/
func linkSetDown(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "linkdown ",
		Short: "linkdown <interfaces ...>",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				out, _, err := cli.Exec("sudo ip link set "+arg+" down", nil)
				if err != nil {
					return err
				}
				fmt.Fprintln(cli.Out(), out)
			}
			return nil
		},
	}
	return cmd
}

func (*dependencies) isDockerAvailable(cli Cli) bool {
	out, _, err := cli.Exec("whereis docker", nil)
	if err != nil {
		panic(err)
	}
	if strings.Contains(out, "/docker") {
		return true
	}
	return false
}

func (*dependencies) isHugePagesEnabled(cli Cli) (bool, int) {
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

func (*dependencies) resizeHugePages(cli Cli, size uint) error {
	//TODO: Handle numa case, Big hugepages(are immutable and can be setted only during booting)
	_, _, err := cli.Exec(fmt.Sprintf("sudo sysctl -w vm.nr_hugepages=%d", size), nil)
	if err != nil {
		return err
	}
	return nil
}

func (*dependencies) installDocker(cli Cli) error {

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
		"sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
	}
	_ = commands

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
func (*dependencies) dumpNetworkInterfaces(cli Cli) []networkInterface {
	const path = "/sys/class/net"
	var list []networkInterface
	var realDevices []networkInterface

	out, _, err := cli.Exec("ls -b", []string{path})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	buff := strings.Fields(out)

	for _, name := range buff {
		list = append(list, networkInterface{name: name})
	}

	for _, nic := range list {
		_, _, err := cli.Exec("ls ", []string{path + "/" + nic.name})
		if err == nil {
			newNic := networkInterface{name: nic.name}

			info, _, _ := cli.Exec("cat", []string{path + "/" + nic.name + "/device/uevent"})

			pciRe := regexp.MustCompile(`PCI_SLOT_NAME=(\S+)`)
			match := pciRe.FindStringSubmatch(info)
			if len(match) == 0 {
				continue
			}
			newNic.pci = match[1]

			driverRe := regexp.MustCompile(`DRIVER=(\S+)`)
			match = driverRe.FindStringSubmatch(info)
			newNic.description = match[1]

			_, _, err = cli.Exec("cat", []string{path + "/" + nic.name + "/carrier"})
			if err != nil {
				newNic.linkUp = false
			} else {
				newNic.linkUp = true
			}

			realDevices = append(realDevices, newNic)
		}

	}

	return realDevices
}

func startupConf(cli Cli) *cobra.Command {
	const startupconfig = "unix {\n    cli-no-pager\n    cli-listen /run/vpp/cli.sock\n    log /tmp/vpp.log\n    coredump-size unlimited\n    full-coredump\n    poll-sleep-usec 50\n}\n\ndpdk {\n\tdev {{range .}} {{.}}\n {{end}} \n}\n\napi-trace {\n    on\n}\n\nsocksvr {\n\tdefault\n}\n\nstatseg {\n\tdefault\n\tper-node-counters on\n}\n\npunt {\n    socket /run/stonework/vpp/punt-to-vpp.sock\n}"
	cmd := &cobra.Command{
		Use:   "startup",
		Short: "startup",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dpdcs := &dependencies{}
			dpdcs.interfaces = dpdcs.dumpNetworkInterfaces(cli)
			pcis := []string{}
			for _, intfc := range dpdcs.interfaces {
				pcis = append(pcis, intfc.pci)
			}

			t := template.Must(template.New("startupConf").Parse(startupconfig))
			err := t.Execute(cli.Out(), pcis)
			if err != nil {
				fmt.Println("Could not execute template")
			}
			return nil
		},
	}
	return cmd
}
