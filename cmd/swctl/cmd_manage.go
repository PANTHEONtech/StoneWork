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
	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"gopkg.in/yaml.v3"
)

type ManageOptions struct {
	Format     string
	Target     string
	Count      uint
	Force      bool
	Offset     int
	DryRun     bool
	Vars       map[string]string
	ShowConfig bool
}

func NewManageCmd(cli Cli) *cobra.Command {
	var (
		opts ManageOptions
	)
	cmd := &cobra.Command{
		Use:              "manage ENTITY [add|del]",
		Short:            "Manage config changes with entities",
		Args:             cobra.ArbitraryArgs,
		TraverseChildren: true,
		//DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(cli.Entities()) == 0 {
				return fmt.Errorf("no entities defined, use --entityfile=FILE to specify file to load from or add 'entities.yaml' to work directory")
			}
			return runManageCmd(cli, opts, args)
		},
	}

	cmd.PersistentFlags().UintVarP(&opts.Count, "count", "c", 1, "Number of instances to add")
	cmd.PersistentFlags().IntVar(&opts.Offset, "offset", 0, "Offset for the starting index")
	cmd.PersistentFlags().BoolVar(&opts.Force, "force", false, "Force the action")
	cmd.PersistentFlags().StringVar(&opts.Target, "target", "", "Target config file used as base")
	cmd.PersistentFlags().StringVar(&opts.Format, "format", "", "Format for the output (yaml, json, proto, go template..)")
	cmd.PersistentFlags().BoolVar(&opts.DryRun, "dryrun", false, "Run without actually modifying anything")
	cmd.PersistentFlags().BoolVar(&opts.ShowConfig, "show-config", false, "Print config for entity detail")
	cmd.PersistentFlags().StringToStringVar(&opts.Vars, "var", nil, "Override values for variables (--var VAR_NAME=value)")

	return cmd
}

func runManageCmd(cli Cli, opts ManageOptions, args []string) error {
	if len(args) > 2 {
		return fmt.Errorf("expected at most 2 arguments (ENTITY ACTION), got %d arguments: %q", len(args), args)
	}

	// list all entities
	if len(args) == 0 {
		color.Fprintf(cli.Out(), "Listing <bold>%d</> loaded entities:\n\n", len(cli.Entities()))
		for _, e := range cli.Entities() {
			name := e.Name
			desc := ""
			if e.Description != "" {
				desc = prefixTmpl(e.Description, "   ") + "\n"
			}
			color.Fprintf(cli.Out(), " - <cyan>%s</> <gray>(%v vars)</>\n<gray>%s</>", name, len(e.GetVars()), desc)
		}
		return nil
	}

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

	logrus.Debugf("managing entity: %v (%d vars), config has %d bytes", entity.GetName(), len(entity.GetVars()), len(entity.Config))

	// print entity detail
	if len(args) == 1 {
		fmt.Fprintf(cli.Out(), "ENTITY: %s\n", color.Yellow.Sprint(entity.Name))
		if len(entity.Description) > 0 {
			desc := strings.TrimSpace(entity.Description)
			fmt.Fprintf(cli.Out(), "DESCRIPTION: %s\n", color.LightWhite.Sprint(desc))
		}
		fmt.Fprintf(cli.Out(), "ORIGIN: %s\n", color.LightWhite.Sprint(entity.Origin))
		fmt.Fprintf(cli.Out(), "VARS:\n")
		maxLen := getEntityVarMaxLen(entity)
		for _, v := range entity.Vars {
			name := color.Style{color.Bold, color.HiWhite}.Sprintf("%-"+fmt.Sprint(maxLen)+"s", v.Name)
			valclr := color.Style{color.OpBold, color.FgLightCyan}
			val := valclr.Sprint(v.Value)
			if v.Type != "" {
				val = fmt.Sprintf("%v (%v)", val, color.FgGray.Sprint(v.Type))
			}
			if v.Description != "" {
				fmt.Fprintf(cli.Out(), "  %v: %v\n%v\n", color.Notice.Sprint(name), val, prefixTmpl(v.Description, "   "))
			} else {
				fmt.Fprintf(cli.Out(), "  %v: %v\n", name, val)
			}
		}
		if opts.ShowConfig {
			fmt.Fprintf(cli.Out(), "CONFIG:\n%s\n", color.LightWhite.Sprint(prefixTmpl(entity.Config, "   ")))
		}
		return nil
	}

	// determine action
	action := args[1]
	switch strings.ToLower(action) {
	case "create", "new", "gen", "generate", "make":
		fallthrough
	case "add", "append":
		action = "ADD"
	case "del", "delete", "rem", "remove":
		action = "DEL"
		if opts.Target == "" {
			return fmt.Errorf("target config file must be specified for delete action")
		}
	default:
		return fmt.Errorf("unknown action %q, supported actions are: 'add' and 'del'", action)
	}

	var finalConf protoreflect.ProtoMessage

	// prepare config
	allConf, err := client.NewDynamicConfig(allModels())
	if err != nil {
		return fmt.Errorf("failed to create dynamic config for all models")
	}

	logrus.Tracef("generating entity config (count: %d, offset: %d)", opts.Count, opts.Offset)

	// repeat for given count
	for i := 0; i < int(opts.Count); i++ {
		idx := i + opts.Offset
		id := idx + 1

		// prepare entity vars
		evars := map[string]string{}
		if !entity.Single {
			evars["IDX"] = fmt.Sprint(idx)
			evars["ID"] = fmt.Sprint(id)
		}

		// apply overrides
		for k, v := range opts.Vars {
			if isIndexVar(k) {
				if entity.Single {
					return fmt.Errorf("single instance entity does not have internal index var: %v", k)
				}
				return fmt.Errorf("cannot override internal var: %q", k)
			}
			var ok bool
			for _, evar := range entity.Vars {
				if k == evar.Name {
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("found override for undefined variable %q", k)
			}
			evars[k] = v
		}

		// render vars
		vars, err := renderEntityVars(entity, evars)
		if err != nil {
			return fmt.Errorf("failed to render vars (idx: %v): %w", i, err)
		}

		// generate config
		config, err := renderEntityConfig(entity, vars)
		if err != nil {
			return fmt.Errorf("failed to render config (idx: %v): %w", i, err)
		}

		logrus.Tracef(" - [#%d] vars: %+v\n%v", i, vars, config)

		conf := allConf.New().Interface()

		playWithYaml(config)

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

	ckeys := extractModelKeysFromConfig(allConf)

	if opts.Target != "" {
		// load target file config
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

		tkeys := extractModelKeysFromConfig(tconf)

		if action == "ADD" {
			// check for conflicts
			conflictKeys := findConflictingKeys(ckeys, tkeys)
			if len(conflictKeys) > 0 {
				logrus.Debugf("listing %d conflicting keys:", len(conflictKeys))
				for _, c := range conflictKeys {
					logrus.Debugf(" - conflict key %v", c)
				}
				return fmt.Errorf("found %d conflicting keys in target config", len(conflictKeys))
			}

			proto.Merge(tconf, allConf)

			logrus.Tracef("MERGED CONFIG:\n%s", yamlTmpl(tconf))

			// extract items
			items, err := client.DynamicConfigExport(tconf.(*dynamicpb.Message))
			if err != nil {
				return err
			}
			logrus.Debugf("extracted %d items", len(items))

			// list extracted items
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
		} else {
			return fmt.Errorf("not yet supported")
		}

		finalConf = tconf
	} else {
		finalConf = allConf
	}

	// format final output
	if opts.Format == "" {
		opts.Format = "yaml"
	}
	if err := formatAsTemplate(cli.Out(), opts.Format, finalConf); err != nil {
		return err
	}

	return nil
}

func playWithYaml(config string) {
	dec := yaml.NewDecoder(strings.NewReader(config))

	var node yaml.Node
	if err := dec.Decode(&node); err != nil {
		logrus.Tracef("ERROR: yaml decode: %v", err)
	}

	logrus.Tracef("NODE:\n%s\n", yamlTmpl(node))
}

func getEntityVarMaxLen(entity Entity) int {
	var max = 3
	for _, v := range entity.Vars {
		if len(v.Name) > max {
			max = len(v.Name)
		}
	}
	return max
}

func extractModelKeysFromConfig(configMsg proto.Message) []string {
	items, err := client.DynamicConfigExport(configMsg.(*dynamicpb.Message))
	if err != nil {
		logrus.Debugf("DynamicConfigExport error: %v", err)
		return nil
	}
	logrus.Debugf("extracted %d items", len(items))
	var keys []string
	for i, item := range items {
		l := logrus.WithFields(map[string]interface{}{
			"item":      fmt.Sprintf("%d/%d", i+1, len(items)),
			"protoName": item.ProtoReflect().Descriptor().FullName(),
		})
		model, err := models.GetModelFor(item)
		if err != nil {
			l.Debugf("no model found for item: %v", item)
			continue
		}
		l = l.WithField("model", model.Name())
		name, err := model.InstanceName(item)
		if err != nil {
			l.Debugf("instance name error: %v", err)
			continue
		}
		if name == "" {
			l.Debugf("intance has empty name, skipping item")
			continue
		}
		key := path.Join(model.KeyPrefix(), name)
		l.Tracef("extracted key: %q", key)
		keys = append(keys, key)
	}
	return keys
}

func findConflictingKeys(slice1 []string, slice2 []string) []string {
	var conflicts []string
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
	"inc": func(a int) int {
		return a + 1
	},
	"dec": func(a int) int {
		return a - 1
	},
	"previp": func(ip string, dec int) (string, error) {
		x, err := netip.ParseAddr(ip)
		if err != nil {
			return "", err
		}
		if dec <= 0 {
			return x.String(), nil
		}
		for i := 1; i <= dec; i++ {
			x = x.Prev()
			if !x.IsValid() {
				return "", fmt.Errorf("no previous IP: %w", err)
			}
		}
		return x.String(), nil
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
				return "", fmt.Errorf("no next IP: %w", err)
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
		template.New("entity").Funcs(funcMap).Option("missingkey=error").Parse(t),
	)
	var b strings.Builder
	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func renderEntityConfig(e Entity, evars map[string]string) (string, error) {
	tmpl, err := interpolateStr(e.Config, evars)
	if err != nil {
		return "", err
	}
	config, err := renderTmpl(tmpl, evars)
	if err != nil {
		return "", err
	}
	return config, nil
}

func renderEntityVars(e Entity, evars map[string]string) (map[string]string, error) {
	for _, v := range e.Vars {
		vv := v.Value
		if ov, ok := evars[v.Name]; ok {
			vv = ov
		}
		tmpl, err := interpolateStr(vv, evars)
		if err != nil {
			return nil, err
		}
		val, err := renderTmpl(tmpl, evars)
		if err != nil {
			return nil, err
		}
		evars[v.Name] = val
	}
	return evars, nil
}

func interpolateStr(str string, vars map[string]string) (string, error) {
	env := interpolate.NewMapEnv(vars)
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
