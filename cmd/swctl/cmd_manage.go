package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"path"
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
	"google.golang.org/protobuf/types/dynamicpb"
)

type ManageOptions struct {
	Format string
	Target string
	Count  uint
	Force  bool
	Offset int
	DryRun bool
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
					fmt.Printf(" - %s (%d options)\n", e.GetName(), len(e.GetOptions()))
				}
				return nil
			}

			return runManageCmd(cli, opts, args)
		},
	}

	cmd.PersistentFlags().UintVarP(&opts.Count, "count", "c", 1, "Number of instances for an action")
	cmd.PersistentFlags().IntVar(&opts.Offset, "offset", 0, "Offset for the starting index")
	cmd.PersistentFlags().BoolVar(&opts.Force, "force", false, "Force the action")
	cmd.PersistentFlags().StringVar(&opts.Target, "target", "", "Target config file to update")
	cmd.PersistentFlags().StringVar(&opts.Format, "format", "", "Format for the output")
	cmd.PersistentFlags().BoolVar(&opts.DryRun, "dryrun", false, "Do not modify anything")

	return cmd
}

func runManageCmd(cli Cli, opts ManageOptions, args []string) error {
	logrus.Tracef("running manage with args: %q %+v", args, opts)

	// lookup entity
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

	// print entity detail
	if len(args) == 1 {
		fmt.Fprintf(cli.Out(), "Entity: %s\n", entity.Name)
		fmt.Fprintf(cli.Out(), "Options:\n")
		for _, o := range entity.Options {
			fmt.Fprintf(cli.Out(), " - %s: %q\n", o.Name, o.Value)
		}
		fmt.Fprintf(cli.Out(), "Config:\n%s\n", entity.Config)
		return nil
	} else if len(args) > 2 {
		return fmt.Errorf("expected 2 arguments (ENTITY ACTION), got %d arguments: %q", len(args), args)
	}

	action := args[1]
	switch action {
	case "add", "create", "new":
		action = "ADD"
	case "del", "delete", "remove":
		action = "DEL"
		if opts.Target == "" {
			return fmt.Errorf("target config file must be specified for delete action")
		}
	default:
		return fmt.Errorf("unknown action %q", action)
	}

	// build config
	allConf, err := client.NewDynamicConfig(allModels())
	if err != nil {
		return fmt.Errorf("failed to create dynamic config for all models")
	}

	logrus.Tracef("generating entity config (count: %d, offset: %d)", opts.Count, opts.Offset)

	for i := 0; i < int(opts.Count); i++ {
		params, err := renderEntityOptions(entity, i+opts.Offset)
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
			return fmt.Errorf("cannot unmarshall data into dynamic config due to: %w", err)
		}

		logrus.Tracef("CONFIG(%d):\n%v", i, yamlTmpl(conf))

		proto.Merge(allConf, conf)
	}

	ckeys := extractModelKeys(allConf)

	if opts.Target != "" {
		config, err := os.ReadFile(opts.Target)
		if err != nil {
			return fmt.Errorf("failed to read target config file: %w", err)
		}

		tconf := allConf.New().Interface()

		bj, err := yaml2.YAMLToJSON(config)
		if err != nil {
			return fmt.Errorf("cannot convert YAML to JSON: %w", err)
		}
		err = protojson.Unmarshal(bj, tconf)
		if err != nil {
			return fmt.Errorf("cannot unmarshall data into dynamic config due to: %w", err)
		}

		logrus.Tracef("TARGET CONFIG:\n%s", yamlTmpl(tconf))

		tkeys := extractModelKeys(tconf)

		if action == "ADD" {
			conflicts := findConflictingStrings(ckeys, tkeys)

			if len(conflicts) > 0 {
				for _, c := range conflicts {
					logrus.Debugf(" - %v", c)
				}
				return fmt.Errorf("%d conflicting config items", len(conflicts))
			}

			proto.Merge(tconf, allConf)

			logrus.Tracef("MERGED CONFIG:\n%s", yamlTmpl(tconf))

			items, err := client.DynamicConfigExport(tconf.(*dynamicpb.Message))
			if err != nil {
				return err
			}

			logrus.Debugf("extracted %d items", len(items))
			for _, item := range items {
				model, err := models.GetModelFor(item)
				if err != nil {
					logrus.Tracef("failed to get model for item: %+v", item)
					continue
				}
				name, err := model.InstanceName(item)
				if err != nil {
					logrus.Tracef("cannot compute instance name due to: %v", err)
					continue
				}
				key := path.Join(model.KeyPrefix(), name)

				data := protojson.MarshalOptions{EmitUnpopulated: true}.Format(item)
				logrus.Tracef(" - %s [%s] %v", model.ProtoName(), key, data)
			}

			if opts.Format == "" {
				opts.Format = "yaml"
			}
			if err := formatAsTemplate(cli.Out(), opts.Format, tconf); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not yet supported")
		}

		return nil
	}

	if opts.Format == "" {
		opts.Format = "yaml"
	}
	if err := formatAsTemplate(cli.Out(), opts.Format, allConf); err != nil {
		return err
	}

	return nil
}

func extractModelKeys(m proto.Message) []string {
	items, err := client.DynamicConfigExport(m.(*dynamicpb.Message))
	if err != nil {
		logrus.Debugf("DynamicConfigExport error: %v", err)
		return nil
	}
	logrus.Debugf("extracted %d items", len(items))

	var keys []string
	for _, item := range items {
		model, err := models.GetModelFor(item)
		if err != nil {
			logrus.Debugf("failed to get model for item: %+v", item)
			continue
		}
		name, err := model.InstanceName(item)
		if err != nil {
			logrus.Debugf("cannot compute instance name due to: %v", err)
			continue
		}
		if name == "" {
			logrus.Debugf("skipping model %v", model.ProtoName())
			continue
		}
		key := path.Join(model.KeyPrefix(), name)
		keys = append(keys, key)
	}
	return keys
}

func findConflictingStrings(slice1 []string, slice2 []string) []string {
	conflicts := []string{}

	for _, str1 := range slice1 {
		for _, str2 := range slice2 {
			if str1 == str2 {
				conflicts = append(conflicts, str1)
			}
		}
	}

	return conflicts
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
			if !x.IsValid() {
				return "", fmt.Errorf("no next ip: %w", err)
			}
		}
		return x.String(), nil
	},
	"subnet": func(addr string, inc int) (string, error) {
		_, ipnet, err := net.ParseCIDR(addr)
		if err != nil {
			return "", err
		}
		if inc <= 0 {
			return ipnet.String(), nil
		}
		ones, _ := ipnet.Mask.Size()
		ipnet.Mask = net.CIDRMask(ones-8, 32)
		ipnet, err = cidr.Subnet(ipnet, 8, inc)
		if err != nil {
			return "", err
		}
		return ipnet.String(), nil
	},
	"trimsuffix": func(s, suffix string) string {
		return strings.TrimSuffix(s, suffix)
	},
	"trimprefix": func(s, prefix string) string {
		return strings.TrimPrefix(s, prefix)
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
