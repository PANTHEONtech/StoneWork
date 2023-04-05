package main

import (
	"fmt"
	"net"
	"net/netip"
	"strings"
	"text/template"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/buildkite/interpolate"
	yaml2 "github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ManageOptions struct {
	Count uint
	Force bool
}

func NewManageCmd(cli Cli) *cobra.Command {
	var (
		opts ManageOptions
	)
	cmd := &cobra.Command{
		Use:              "manage ENTITY",
		Short:            "Manage entities in StoneWork",
		Args:             cobra.ArbitraryArgs,
		TraverseChildren: true,
		//DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Tracef("running manage with %d args: %v", len(args), args)

			if len(args) == 0 {
				logrus.Tracef("listing %d entities", len(cli.Entities()))

				fmt.Printf("Entities:\n")
				for _, e := range cli.Entities() {
					fmt.Printf(" - %s\n", e.GetName())
				}
				return nil
			}

			return runManageCmd(cli, opts, args)
		},
	}

	cmd.PersistentFlags().UintVarP(&opts.Count, "count", "c", 1, "Number of instances for an action")

	return cmd
}

func runManageCmd(cli Cli, opts ManageOptions, args []string) error {
	logrus.Tracef("running manage with args: %q %+v", args, opts)

	entityName := strings.ToLower(args[0])
	var entity Entity
	for _, e := range cli.Entities() {
		if strings.ToLower(e.GetName()) == entityName {
			entity = e
			break
		}
	}
	if entity.Name == "" {
		return fmt.Errorf("entity %q not found", entityName)
	}

	logrus.Debugf("entity: %v (%d options) and config (%d bytes)", entity.GetName(), len(entity.GetOptions()), len(entity.Config))

	allConf, err := client.NewDynamicConfig(allModels())
	if err != nil {
		return fmt.Errorf("failed to create dynamic config for all models")
	}

	logrus.Tracef("generating entity with count %d", opts.Count)

	for i := 0; i < int(opts.Count); i++ {
		params, err := renderEntityOptions(entity, i)
		if err != nil {
			return fmt.Errorf("failed to render parameters (idx: %v): %w", i, err)
		}

		config, err := interpolateStr(entity.Config, params)
		if err != nil {
			return fmt.Errorf("failed to interpolate parameters (idx: %v): %w", i, err)
		}

		logrus.Tracef(" - [#%d] params: %+v\n%v", i, params, config)

		conf := allConf.New().Interface()

		bj, err := yaml2.YAMLToJSON([]byte(config))
		if err != nil {
			return fmt.Errorf("cannot convert YAML to JSON: %w", err)
		}
		err = protojson.Unmarshal(bj, conf)
		if err != nil {
			return fmt.Errorf("cannot unmarshall init file data into dynamic config due to: %w", err)
		}

		proto.Merge(allConf, conf)
	}

	logrus.Debugf("CONFIG:\n%v", yamlTmpl(allConf))

	return nil
}

var funcMap = map[string]any{
	"add": func(a, b int) int {
		return a + b
	},
	"nextip": func(ip string, inc int) (string, error) {
		x, err := netip.ParseAddr(ip)
		if err != nil {
			return "", err
		}
		if inc <= 0 {
			return x.String(), nil
		}
		for i := 1; i <= inc; i++ {
			x = x.Next()
		}
		return x.String(), nil
	},
	"nextsubnet": func(addr string) (string, error) {
		_, ipnet, err := net.ParseCIDR(addr)
		if err != nil {
			return "", err
		}
		_, bits := ipnet.Mask.Size()
		next, ok := cidr.NextSubnet(ipnet, bits)
		if !ok {
			return "", fmt.Errorf("failed")
		}
		return next.String(), nil
	},
}

func renderTmpl(t string, data any) (string, error) {
	tmpl := template.Must(
		template.New("params").Funcs(funcMap).Option("missingkey=error").Parse(t),
	)
	var b strings.Builder
	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func renderEntityOptions(e Entity, idx int) (map[string]string, error) {
	id := idx + 1
	opts := map[string]string{
		"IDX": fmt.Sprint(idx),
		"ID":  fmt.Sprint(id),
	}
	for _, o := range e.Options {
		tmpl, err := interpolateStr(o.Value, opts)
		if err != nil {
			return nil, err
		}
		val, err := renderTmpl(tmpl, opts)
		if err != nil {
			return nil, err
		}
		opts[o.Name] = val
	}
	return opts, nil
}

func listOptions(str string) []string {
	idents, err := interpolate.Identifiers(str)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	//logrus.Infof("IDENTS: %+v", idents)

	var opts []string
	uniq := map[string]struct{}{}

	for _, ident := range idents {
		if _, ok := uniq[ident]; ok {
			continue
		}
		uniq[ident] = struct{}{}
		opts = append(opts, ident)
	}

	return opts
}

func interpolateStr(str string, vars map[string]string) (string, error) {
	env := interpolate.NewMapEnv(vars)

	//opts := listOptions(str)
	//logrus.Infof("opts: %+v", opts)

	output, err := interpolate.Interpolate(env, str)
	if err != nil {
		return "", err
	}
	return output, nil
}

func allModels() []*models.ModelInfo {
	var knownModels []*models.ModelInfo
	for _, model := range models.RegisteredModels() {
		if model.Spec().Class == "config" {
			knownModels = append(knownModels, &models.ModelInfo{
				ModelDetail:       model.ModelDetail(),
				MessageDescriptor: model.NewInstance().ProtoReflect().Descriptor(),
			})
		}
	}
	return knownModels
}

func flagValuesFromActionOptions(set *pflag.FlagSet, entity Entity) map[string]string {
	opts := make(map[string]string)
	for _, opt := range entity.Options {
		if opt.Type == "" || opt.Type == "string" {
			var str string
			if set.Changed(opt.Name) {
				var err error
				str, err = set.GetString(opt.Name)
				if err != nil {
					str = err.Error()
				}
			} else {
				str = opt.Value
			}
			opts[opt.Name] = str
		}
	}
	return opts
}

func flagSetFromActionOptions(entity Entity) *pflag.FlagSet {
	set := pflag.NewFlagSet(entity.Name, pflag.ContinueOnError)
	for _, opt := range entity.Options {
		switch opt.Type {
		case "":
			fallthrough
		case "string":
			set.String(opt.Name, opt.Value, opt.Description)
		}
	}
	return set
}
