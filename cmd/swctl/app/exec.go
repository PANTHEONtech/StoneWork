package app

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

const (
	MissingExternalToolMessageKey = "MissingExternalTool"
)

type externalExe string

// known external programs that swctl can execute
const (
	cmdDocker   externalExe = "docker"
	cmdVppProbe externalExe = "vpp-probe"
	cmdAgentCtl externalExe = "agentctl"
)

func (ee externalExe) installPath() string {
	return filepath.Join(binaryToolsInstallDir, string(ee))
}

func (ee externalExe) validateUsability(cli *CLI) error {
	// check for file existence
	fileInfo, err := os.Stat(ee.installPath())
	if err != nil {
		if os.IsNotExist(err) {
			if messageOverride, found := cli.customizations[MissingExternalToolMessageKey]; found {
				return fmt.Errorf(messageOverride.(string), ee.installPath())
			}
			return fmt.Errorf("%s is not installed, try running "+
				"`swctl dependency install-tools`", ee.installPath())
		}
		return fmt.Errorf("existence of %s could not be determined due to %w", ee.installPath(), err)
	}

	// check for execution permissions
	if uint32(fileInfo.Mode()&0111) != 0111 {
		// someone can execute the file (file owner or owner's group or someone else)
		// -> need to check if this process can execute it with more expensive OS call
		if unix.Access(ee.installPath(), unix.X_OK) != nil {
			return fmt.Errorf("%s can't be executed by this process/user", ee.installPath())
		}
	}

	return nil
}

type externalCmd struct {
	cli   *CLI
	exe   externalExe
	name  string
	args  []string
	env   []string
	color bool
}

func newExternalCmd(cmd externalExe, args []string, cli *CLI) *externalCmd {
	ecmd := &externalCmd{
		exe:   cmd,
		name:  string(cmd),
		cli:   cli,
		color: cli.GlobalOptions().Color != "never" && cli.Out().IsTerminal(),
	}
	ecmd.args = append(ecmd.args, args...)
	ecmd.setDebugArg(cli.GlobalOptions().Debug)
	ecmd.setLogLevelArg(logrus.GetLevel())
	ecmd.setColorEnv()
	ecmd.setMiscFlags()
	return ecmd
}

func (ec *externalCmd) setDebugArg(debug bool) {
	switch ec.exe {
	case cmdAgentCtl, cmdDocker, cmdVppProbe:
		if debug {
			ec.prependUniqueArg("--debug", "", "-D")
		}
	}
}

func (ec *externalCmd) setLogLevelArg(loglvl logrus.Level) {
	if loglvl < logrus.FatalLevel {
		loglvl = logrus.FatalLevel
	}
	if loglvl > logrus.DebugLevel {
		loglvl = logrus.DebugLevel
	}
	switch ec.exe {
	case cmdAgentCtl, cmdDocker:
		ec.prependUniqueArg("--log-level", loglvl.String(), "-l")
	case cmdVppProbe:
		ec.prependUniqueArg("--loglevel", loglvl.String(), "-L")
	}
}

// https://no-color.org/
func (ec *externalCmd) setColorEnv() {
	if !ec.color {
		ec.env = append(ec.env, "NO_COLOR=1")
	}
}

func (ec *externalCmd) setMiscFlags() {
	switch ec.exe {
	case cmdVppProbe:
		if ec.color {
			ec.prependUniqueArg("--color", "always")
		}
	case cmdAgentCtl:
		ec.prependUniqueArg("--host", ec.cli.client.GetHost(), "-H")
	case cmdDocker:
		if i := slices.Index(ec.args, "compose"); i >= 0 {
			var argVals []string
			for _, cf := range ec.cli.GlobalOptions().ComposeFiles {
				argVals = append(argVals, fmt.Sprintf("--file=%s", cf))
			}
			ec.args = tryInsertArgVals(ec.args, i+1, argVals...)
		}
	}
}

type ExecResult struct {
	Took   time.Duration
	Status int
	Stdout string
	Stderr string
}

func (ec *externalCmd) exec(liveOutput bool) (*ExecResult, error) {
	var stdout, stderr bytes.Buffer

	// compute executable name that will define used filepath to executable file
	executable := ec.name // just file name -> using PATH to resolve to absolute path
	if ec.exe == cmdVppProbe || ec.exe == cmdAgentCtl {
		// this is override of executed command to take the properly installed version of external command
		// instead of whetever is in the linux PATH
		executable = filepath.Join(binaryToolsInstallDir, string(ec.exe)) // absolute path to external tool
		if err := ec.exe.validateUsability(ec.cli); err != nil {
			return nil, fmt.Errorf("can't use external tool %s due to: %w", string(ec.exe), err)
		}
	}

	cmd := exec.Command(executable, ec.args...)
	if liveOutput {
		stdoutMultiWriter := io.MultiWriter(&stdout, ec.cli.out)
		stderrMultiWriter := io.MultiWriter(&stderr, ec.cli.err)
		cmd.Stdout = stdoutMultiWriter
		cmd.Stderr = stderrMultiWriter
	} else {
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	cmd.Env = ec.env

	now := time.Now()
	logrus.Tracef("[%s] %q", color.Gray.Sprint("EXEC"), cmd.String())
	err := cmd.Run()
	took := time.Since(now)
	l := logrus.WithField("took", took.Round(1*time.Millisecond).String())

	var execRes ExecResult
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			l.Tracef("[%s] %q (%v)\n%s\n", color.Red.Sprint("ERR"), cmd.Args, color.Red.Sprint(err), color.LightRed.Sprint(stderr.String()))
			execRes.Status = ee.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			return nil, err
		}
	}

	execRes.Took = took
	execRes.Stdout = strings.TrimRight(stdout.String(), "\n")
	execRes.Stderr = strings.TrimRight(stderr.String(), "\n")
	l.Tracef("[%s] %q\n%s\n", color.Green.Sprint("OK"), cmd.Args, color.FgGray.Sprint(execRes.Stdout))
	return &execRes, nil
}

func (ec *externalCmd) prependUniqueArg(arg string, val string, aliases ...string) {
	ec.args = tryInsertUniqueArg(ec.args, 0, arg, val, aliases...)
}

func tryInsertUniqueArg(args []string, pos int, arg string, val string, aliases ...string) []string {
	if hasAnyPrefix(args, arg) || hasAnyPrefix(args, aliases...) {
		return args
	}
	return tryInsertArg(args, pos, arg, val)
}

func tryInsertArg(args []string, pos int, arg string, val string) []string {
	var argVal string
	if val == "" {
		argVal = arg
	} else {
		argVal = fmt.Sprintf("%s=%s", arg, val)
	}
	return tryInsertArgVals(args, pos, argVal)
}

func tryInsertArgVals(args []string, pos int, argVals ...string) []string {
	if pos < 0 || pos > len(args) {
		return args
	}
	return slices.Insert(args, pos, argVals...)
}

func anyPrefixIndex(elems []string, prefixes ...string) int {
	for i, elem := range elems {
		for _, pref := range prefixes {
			if strings.HasPrefix(elem, pref) {
				return i
			}
		}
	}
	return -1
}

func hasAnyPrefix(elems []string, prefixes ...string) bool {
	return anyPrefixIndex(elems, prefixes...) >= 0
}
