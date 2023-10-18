package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	composecli "github.com/compose-spec/compose-go/cli"
	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	EnvVarDebug    = "SWCTL_DEBUG"
	EnvVarLogLevel = "SWCTL_LOGLEVEL"
)

type GlobalOptions struct {
	Debug    bool
	LogLevel string
	Color    string

	ComposeFiles       []string
	EntityFiles        []string
	EmbeddedEntityByte []byte

	// TODO: support config file
	// Config string
}

func InitGlobalOptions(cli Cli, opts *GlobalOptions) {
	// color mode
	if opts.Color == "" && os.Getenv("NO_COLOR") != "" {
		// https://no-color.org/
		opts.Color = "never"
	}
	switch strings.ToLower(opts.Color) {
	case "auto", "":
		if !cli.Out().IsTerminal() {
			color.Disable()
		}
	case "on", "enabled", "always", "1", "true":
		color.Enable = true
	case "off", "disabled", "never", "0", "false":
		color.Disable()
	default:
		logrus.Fatalf("invalid color mode: %q", opts.Color)
	}

	// debug mode
	if os.Getenv(EnvVarDebug) != "" {
		opts.Debug = true
	}

	// log level
	if loglvl := os.Getenv(EnvVarLogLevel); loglvl != "" {
		opts.LogLevel = loglvl
	}
	if opts.LogLevel != "" {
		if lvl, err := logrus.ParseLevel(opts.LogLevel); err == nil {
			logrus.SetLevel(lvl)
			if lvl >= logrus.TraceLevel {
				logrus.SetReportCaller(true)
				//infralogrus.DefaultLogger().SetLevel(logging.LogLevel(lvl))
			}
			logrus.Tracef("log level set to %v", lvl)
		} else {
			logrus.Fatalf("log level invalid: %v", err)
		}
	} else if opts.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
		//infralogrus.DefaultLogger().SetLevel(logging.ErrorLevel)
	}

	var err error
	if len(opts.ComposeFiles) == 0 {
		for _, f := range composecli.DefaultFileNames {
			absPath, err := filepath.Abs(f)
			if err != nil {
				continue
			}
			if _, err := os.Stat(absPath); err == nil {
				opts.ComposeFiles = append(opts.ComposeFiles, absPath)
				break
			}
		}
	}
	if opts.ComposeFiles, err = initFiles(opts.ComposeFiles); err != nil {
		logrus.Fatal(err)
	}
	if opts.EntityFiles, err = initFiles(opts.EntityFiles); err != nil {
		logrus.Fatal(err)
	}
	logrus.Debugf("Initialized compose files: %v", opts.ComposeFiles)
	logrus.Debugf("Initialized entity files: %v", opts.EntityFiles)
}

func (glob *GlobalOptions) InstallFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&glob.Debug, "debug", "D", false, "Enable debug mode")
	flags.StringVarP(&glob.LogLevel, "log-level", "L", "", "Set logging level")
	flags.StringVar(&glob.Color, "color", "", "Color mode; auto/always/never")

	flags.StringSliceVar(&glob.ComposeFiles, "composefile", nil, "Docker Compose configuration files")
	must(cobra.MarkFlagFilename(flags, "composefile", "yaml", "yml"))
	flags.StringSliceVar(&glob.EntityFiles, "entityfile", nil, "Entity configuration files")
	must(cobra.MarkFlagFilename(flags, "entityfile", "yaml", "yml"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func initFiles(files []string) ([]string, error) {
	var res []string
	for _, path := range files {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("can't obtain absolute path of path %s due to: %w", path, err)
		}
		if _, err := os.Stat(absPath); err != nil {
			return nil, fmt.Errorf("can't retrieve info about file %s due to: %w", absPath, err)
		}
		res = append(res, absPath)
	}
	return res, nil
}
