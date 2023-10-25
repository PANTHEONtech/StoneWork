package app

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:embed docker/agentctl.Dockerfile
var embeddedAgentctlDockerfile []byte

const (
	dockerVersion         = "default"
	vppProbeTagVersion    = "v0.2.0"
	agentctlCommitVersion = "723f8db0bf7a67908e2dda1d860444a4747a99d8"
)

var binaryToolsInstallDir = filepath.Join(os.Getenv("HOME"), ".cache", "stonework", "bin")

func exampleDependencyCmd(appName string) string {
	return `
  <white># Status of all dependencies</>
  $ <yellow>` + appName + ` dependency status</>

  <white># Install external tools (docker, docker compose, vpp-probe, agentctl)</>
  $ <yellow>` + appName + ` dependency install-tools</>

  <white># Set quantity of runtime 2MB HugePages manually</>
  $ <yellow>` + appName + ` dependency set-hugepages <value></>

  <white># Assign(up) or Unassign(down) interfaces to/from kernel</>
  $ <yellow>` + appName + ` dependency link <pci ...> up | down</>

  <white># Print out VPP startup config with dpdk interfaces</>
  $ <yellow>` + appName + ` dependency get-vpp-startup [<interfacePci:StoneworkInterfaceName ...>]</>

  <white># Print out VPP startup config with dpdk plugin disable</>
  $ <yellow>` + appName + ` dependency get-vpp-startup</>
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

	cmd.AddCommand(installExternalToolsCmd(cli),
		dependencyStatusCmd(cli),
		installHugePagesCmd(cli),
		linkSetUpDownCmd(cli),
		vppStartupConfCmd(cli))

	return cmd
}

func installExternalToolsCmd(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-tools",
		Short: "Install external tools",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			color.Fprintln(cli.Out(), "checking docker availability...")
			docker, err := IsDockerAvailable(cli)
			if err != nil {
				return errors.New(fmt.Sprintf("Unable to check docker availability: %v", err))
			}
			if docker {
				color.Fprintln(cli.Out(), "Docker is already installed")
			} else {
				color.Fprintln(cli.Out(), "Installing docker...")
				err = InstallDocker(cli, dockerVersion)
				if err != nil {
					return err
				}
				color.Fprintln(cli.Out(), "Installation of docker was successful")
			}

			vppProbeAvailable, err := IsVPPProbeAvailable(vppProbeTagVersion, cli.Out())
			if err != nil {
				return errors.New(fmt.Sprintf("Unable to check vpp-probe availability: %v", err))
			}
			if vppProbeAvailable {
				color.Fprintln(cli.Out(), "VPP-probe is already installed")
			} else {
				color.Fprintln(cli.Out(), "Installing vpp-probe...")
				err = InstallVPPProbe(cli, vppProbeTagVersion)
				if err != nil {
					return err
				}
				color.Fprintln(cli.Out(), "Installation of the vpp-probe tool was successful")
			}

			agentctlAvailable, err := IsAgentctlAvailable(agentctlCommitVersion, cli.Out())
			if err != nil {
				return errors.New(fmt.Sprintf("Unable to check agentctl availability: %v", err))
			}
			if agentctlAvailable {
				color.Fprintln(cli.Out(), "Agentctl is already installed")
			} else {
				color.Fprintln(cli.Out(), "Installing agentctl...")
				err = InstallAgentCtl(cli, agentctlCommitVersion)
				if err != nil {
					return err
				}
				color.Fprintln(cli.Out(), "Installation of the agentctl tool was successful")
			}
			if !docker {
				color.Fprintln(cli.Out(), color.Red.Sprintf("Please restart for the docker group membership changes to take effect"))
			}
			return nil
		},
	}
	return cmd
}

func dependencyStatusCmd(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "status",
		Short:         "status",
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// docker
			dockerAvailable, err := IsDockerAvailable(cli)
			status := "Not installed"
			if err != nil {
				status = fmt.Sprintf("<unable to check: %s>", err)
			} else if dockerAvailable {
				status = "OK"
			}
			color.Fprintf(cli.Out(), "Docker: %s\n", status)

			// vpp-probe
			vppProbeAvailable, err := IsVPPProbeAvailable(vppProbeTagVersion, nil)
			status = fmt.Sprintf("Not installed/installed incorrect version "+
				"(needed version is %s)", vppProbeTagVersion)
			if err != nil {
				status = fmt.Sprintf("<unable to check: %s>", err)
			} else if vppProbeAvailable {
				status = "OK"
			}
			color.Fprintf(cli.Out(), "VPP-Probe: %s\n", status)

			// agentctl
			agentctlAvailable, err := IsAgentctlAvailable(agentctlCommitVersion, nil)
			status = fmt.Sprintf("Not installed/installed incorrect version "+
				"(needed version from commit %s in ligato/vpp-agent repository)", agentctlCommitVersion)
			if err != nil {
				status = fmt.Sprintf("<unable to check: %s>", err)
			} else if agentctlAvailable {
				status = "OK"
			}
			color.Fprintf(cli.Out(), "Agentctl: %s\n", status)

			// hugepages
			hugePagesCount, err := AllocatedHugePages(cli)
			status = "Disabled"
			if err != nil {
				status = fmt.Sprintf("<unable to check: %s>", err)
			} else if hugePagesCount != 0 {
				status = strconv.Itoa(hugePagesCount)
			}
			color.Fprintf(cli.Out(), "Hugepages: %s\n", status)

			// physical interfaces
			physicalInterfaces, err := DumpDevices(cli)
			if err != nil {
				color.Fprintf(cli.Out(), "Physical interfaces: <unable to check: %s>\n", err)
				return err
			}
			if physicalInterfaces == nil {
				color.Fprint(cli.Out(), "No available interfaces\n")
			} else {
				color.Fprint(cli.Out(), "Physical interfaces:\n")
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
			// errors are already logged to console, so returning only last error to indicate
			// partial/full failure in cmd status/return code
			return err
		},
	}
	return cmd
}

func installHugePagesCmd(cli Cli) *cobra.Command {
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

// linkSetUpDownCmd changes the link state and binds/unbinds the pci driver
// (up=kernel driver usable for kernel network stack, down=no kernel driver).
// The kernel driver unbinding is helpfull in case of DPDK interfaces that
// can't have kernel network stack usable driver when VPP should use them
// as DPDK interfaces.
func linkSetUpDownCmd(cli Cli) *cobra.Command {
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
						color.Fprintln(cli.Out(), fmt.Sprintf("don't need to unbind the already unbinded pci %s", physicalInterfaces[matchId].Name))
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
					stdout, stderr, err := cli.Exec("ls", []string{"/sys/bus/pci/devices/" + physicalInterfaces[matchId].Pci + "/net"}, false)
					if stderr != "" {
						return errors.New(stderr)
					}
					if err != err {
						return err
					}
					if stdout != "" {
						_, _, err = cli.Exec(GetSudoPrefix(cli)+"ip link set "+stdout+" down", nil, false)
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

func vppStartupConfCmd(cli Cli) *cobra.Command {
	const vppStartupconfig = `unix {
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
		Use:   "get-vpp-startup",
		Short: "Print out VPP startup config",
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

			t := template.Must(template.New("vppStartupConf").Parse(vppStartupconfig))
			err := t.Execute(cli.Out(), desiredInterfaces)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func IsDockerAvailable(cli Cli) (bool, error) {
	out, _, err := cli.Exec("whereis docker", nil, false)
	if err != nil {
		return false, err
	}
	if strings.Contains(out, "/docker") {
		return true, nil
	}
	return false, nil
}

func AllocatedHugePages(cli Cli) (int, error) {
	out, _, err := cli.Exec("sysctl vm.nr_hugepages -n", nil, false)
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
		color.Fprintln(cli.Out(), "Skipping hugepages")
		return nil
	}
	_, _, err := cli.Exec(fmt.Sprintf(GetSudoPrefix(cli)+"sysctl -w vm.nr_hugepages=%d", size), nil, false)
	if err != nil {
		return err
	}
	allocatedHP, err := AllocatedHugePages(cli)
	if err != nil {
		return err
	}

	if size != uint(allocatedHP) {
		return errors.New(fmt.Sprintf("failed to allocate enough hugepages (%d),successfully allocated %d hugepages, totally continuous memory %d MB ",
			size, allocatedHP, (allocatedHP*hugePageSize)/1000))
	}

	return nil
}

func InstallDocker(cli Cli, dockerVersion string) error {
	commands := []string{"apt-get update -y",
		"apt-get install ca-certificates curl gnupg -y",
		"install -m 0755 -d /etc/apt/keyrings",
		"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | " + GetSudoPrefix(cli) + "gpg --dearmor -o /etc/apt/keyrings/docker.gpg --yes",
		"chmod a+r /etc/apt/keyrings/docker.gpg",
		`echo \
		"deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
		"$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
		` + GetSudoPrefix(cli) + `tee /etc/apt/sources.list.d/docker.list > /dev/null`,
		"apt-get update -y",
		"echo \"uio_pci_generic\" | " + GetSudoPrefix(cli) + "tee -a /etc/modules",
	}
	if dockerVersion == "default" {
		cmd := `apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin`
		commands = append(commands, cmd)
	} else {
		cmd := `apt-get install -y docker-ce=` + dockerVersion + ` docker-ce-cli=` + dockerVersion + ` containerd.io docker-buildx-plugin docker-compose-plugin`
		commands = append(commands, cmd)

	}

	for _, command := range commands {
		out, stderr, err := cli.Exec(GetSudoPrefix(cli)+"bash -c", []string{command}, false)
		if stderr != "" {
			return errors.New(command + ": " + stderr)
		}
		if err != nil {
			return errors.New(err.Error() + "(" + command + ")")
		}
		color.Fprintln(cli.Out(), out)

	}
	err := dockerPostInstall(cli)
	if err != nil {
		return errors.New("dockerPostInstall:" + err.Error())
	}

	return nil
}

func InstallVPPProbe(cli Cli, vppProbeTagVersion string) error {
	const (
		repoOwner = "ligato"
		repoName  = "vpp-probe"
	)

	// Construct the path to the installation file
	installPath := cmdVppProbe.installPath() // absolute path to vpp-probe executable
	versionPath := installPath + ".version"

	// Get release info from vpp-probe github repo
	assetUrl, err := retrieveReleaseAssetUrl(repoOwner, repoName, vppProbeTagVersion)
	if err != nil {
		return err
	}

	// Create the installation directory if it doesn't exist
	if err := os.MkdirAll(binaryToolsInstallDir, 0755); err != nil {
		return err
	}

	// Get the vpp-probe binary and copy it to install
	err = downloadAndExtractSubAsset(assetUrl, "vpp-probe", installPath)
	if err != nil {
		return err
	}

	// Store the release version info
	if err := os.WriteFile(versionPath, []byte(vppProbeTagVersion), 0755); err != nil {
		return fmt.Errorf("writing version to file failed: %w", err)
	}

	return nil
}
func IsVPPProbeAvailable(vppProbeTagVersion string, logger io.Writer) (bool, error) {
	return isExternalToolAvailable(cmdVppProbe, vppProbeTagVersion, logger)
}

func IsAgentctlAvailable(agentctlCommitVersion string, logger io.Writer) (bool, error) {
	return isExternalToolAvailable(cmdAgentCtl, agentctlCommitVersion, logger)
}

func isExternalToolAvailable(tool externalExe, targetVersion string, logger io.Writer) (bool, error) {
	// Construct the path to the installation file
	installPath := tool.installPath() // absolute path to tool executable
	versionPath := installPath + ".version"

	// Check current state of tool installation and tool version file in the system
	if logger != nil {
		color.Fprintf(logger, "checking availability of external tool %s\n", string(tool))
	}
	var installedVersion string
	if _, err := os.Stat(installPath); err == nil {
		version, err := os.ReadFile(versionPath)
		if err == nil {
			installedVersion = string(version)
		} else if os.IsNotExist(err) {
			if logger != nil {
				color.Fprintf(logger, "%s version file not found, proceed to download\n", string(tool))
			}
			return false, nil
		} else if err != nil {
			return false, err
		}
	} else if os.IsNotExist(err) {
		if logger != nil {
			color.Fprintf(logger, "%s installation not found, proceed to download\n", string(tool))
		}
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Check whether desired version is already in place
	if installedVersion != "" {
		if logger != nil {
			color.Fprintf(logger, "installed version of %s: %v\n", string(tool), installedVersion)
		}
		if installedVersion == targetVersion {
			if logger != nil {
				color.Fprintf(logger, "installed version of %s is the correct version to be used\n", string(tool))
			}
			return true, nil
		}
		if logger != nil {
			color.Fprintf(logger, "required version of %s is %s, proceed to download\n",
				string(tool), targetVersion)
		}
	}
	return false, nil
}

func InstallAgentCtl(cli Cli, agentctlCommitVersion string) error {
	const builderImage = "agentctl.builder:latest"

	// write embedded agentctl builder dockerfile to tmp folder
	dir, err := os.MkdirTemp("", "agentctl-builder-*")
	if err != nil {
		return fmt.Errorf("can't create tmp folder for agentctl building due to %w", err)
	}
	defer os.RemoveAll(dir) // cleanup of files needed for docker build of agentctl
	dockerFile := filepath.Join(dir, "agentctl.Dockerfile")
	if err = os.WriteFile(dockerFile, embeddedAgentctlDockerfile, 0644); err != nil {
		return fmt.Errorf("can't write agenctl builder dockerfile to tmp folder %s due to %w", dir, err)
	}

	// run agentctl build in docker
	color.Fprintln(cli.Out(), "building agentctl in docker container...")
	_, _, err = cli.Exec(GetSudoPrefix(cli)+"docker build", []string{
		"-f", dockerFile,
		"--build-arg", fmt.Sprintf("COMMIT=%s", agentctlCommitVersion),
		"-t", builderImage,
		"--rm=true",
		dir},
		true)
	if err != nil {
		return fmt.Errorf("can't build agentctl due to builder docker build failure: %w", err)
	}
	defer func() { // cleanup of docker image
		_, _, err = cli.Exec("docker rmi", []string{"-f", builderImage}, false)
		if err != nil {
			color.Fprintf(cli.Out(), "clean up of agentctl builder image failed (%v), continuing... ", err)
		}
	}()

	// extract agentctl into external tools binary folder
	stdout, _, err := cli.Exec(GetSudoPrefix(cli)+"docker create", []string{builderImage}, false)
	if err != nil {
		return fmt.Errorf("can't extract agentctl from builder docker image due "+
			"to container creation failure: %w", err)
	}
	containerId := fmt.Sprint(stdout)
	defer func() { // cleanup of docker container
		_, _, err = cli.Exec(GetSudoPrefix(cli)+"docker rm", []string{containerId}, false)
		if err != nil {
			color.Fprintf(cli.Out(), "clean up of agentctl builder container failed (%v), continuing... ", err)
		}
	}()
	_, _, err = cli.Exec(GetSudoPrefix(cli)+"docker cp", []string{
		fmt.Sprintf("%s:/go/bin/agentctl", containerId), cmdAgentCtl.installPath(),
	}, false)
	if err != nil {
		return fmt.Errorf("can't extract agentctl from builder docker image due "+
			"to docker cp failure: %w", err)
	}

	_, stderr, err := cli.Exec(GetSudoPrefix(cli)+"chmod", []string{"755", cmdAgentCtl.installPath()}, false)
	if stderr != "" {
		return fmt.Errorf("InstallAgentCtl: chmod error: %s", stderr)
	}
	if err != nil {
		return fmt.Errorf("InstallAgentCtl: %s", err)
	}

	// store the release version info
	versionPath := cmdAgentCtl.installPath() + ".version"
	if err := os.WriteFile(versionPath, []byte(agentctlCommitVersion), 0755); err != nil {
		return fmt.Errorf("writing version to file failed: %w", err)
	}

	return nil
}

func unbindDevice(cli Cli, pci string, driver string) error {
	// dpdk drivers like uio_pci_generic, vfio-pci etc..
	// kernel drivers like e1000, ...
	//Mostly
	path := fmt.Sprintf("/sys/bus/pci/drivers/%s/unbind", driver)

	_, stderr, err := cli.Exec(GetSudoPrefix(cli)+"bash -c", []string{"echo \"" + pci + "\" > " + path}, false)
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

	_, stderr, err := cli.Exec(GetSudoPrefix(cli)+"bash -c", []string{"echo \"" + pci + "\" > " + path}, false)
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

	stdout, _, err := cli.Exec("lspci", []string{"-Dvmmnnk"}, false)
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

func GetSudoPrefix(cli Cli) string {
	isRoot, err := isUserRoot()
	if err != nil {
		fmt.Fprintf(cli.Err(), "isUserRoot: %s", err)
		return "sudo "
	}
	if isRoot {
		return ""
	}
	return "sudo "

}

func isUserRoot() (bool, error) {
	user, err := user.Current()
	if err != nil {
		return false, err
	}
	if user.Uid != "0" {
		return false, nil
	}
	return true, nil

}
func dockerPostInstall(cli Cli) error {
	sudoName, err := logname(cli)
	logrus.Tracef("detected login name: %s\n", sudoName)
	if err != nil {
		// handling case when linux has no login name (e.g. container)
		if strings.Contains(err.Error(), "no login") {
			return nil
		}
		return err
	}

	command := fmt.Sprintf("usermod -aG docker %s", sudoName) // add user do docker group

	out, stderr, err := cli.Exec(GetSudoPrefix(cli)+"bash -c", []string{command}, false)
	logrus.Tracef("dockerPostInstall %s: %s\n", command, out)

	if stderr != "" {
		return errors.New(command + ": " + stderr)
	}
	if err != nil {
		return errors.New(err.Error() + "(" + command + ")")
	}
	return err
}

func logname(cli Cli) (string, error) {
	stdout, stderr, err := cli.Exec("logname", nil, false)
	if stderr != "" {
		return "", errors.New(stderr)
	}
	if err != nil {
		return "", err
	}
	return stdout, err
}
