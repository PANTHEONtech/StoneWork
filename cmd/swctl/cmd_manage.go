package main

import (
	"fmt"
	"io"
	"net"
	"net/netip"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/buildkite/interpolate"
	yaml2 "github.com/ghodss/yaml"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.ligato.io/vpp-agent/v3/client"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

const exampleManageCmd = `
  <white># List all available entities</>
  $ <yellow>swctl manage</>

  <white># Show specific entity info</>
  $ <yellow>swctl manage ENTITY</>

  <white># Create entity config with defaults</>
  $ <yellow>swctl manage ENTITY create</>

  <white># With an offset for ID/IDX vars</>
  $ <yellow>swctl manage ENTITY create --id=100</>

  <white># Create multiple entity instances</>
  $ <yellow>swctl manage ENTITY create --count=5</>

  <white># Merge with existing config file</>
  $ <yellow>swctl manage ENTITY add --target="config.yaml""</>

  <white># Override default value of entity variables</>
  $ <yellow>swctl manage ENTITY create --var "VAR=VAL" --var "VAR2=VAL2"</>

  <white># Use interactive mode</>
  $ <yellow>swctl manage ENTITY create --interactive</>
`

type ManageOptions struct {
	Format      string
	Target      string
	Count       uint
	Force       bool
	IdOffset    int
	DryRun      bool
	Vars        map[string]string
	ShowConfig  bool
	Interactive bool
}

func (opts *ManageOptions) InstallFlags(flagset *pflag.FlagSet) {
	flagset.UintVarP(&opts.Count, "count", "c", 1, "Number of instances to add")
	flagset.IntVar(&opts.IdOffset, "id", 0, "Offset for the starting ID")
	flagset.BoolVar(&opts.Force, "force", false, "Force the action")
	flagset.StringVar(&opts.Target, "target", "", "Target config file used as base")
	flagset.StringVar(&opts.Format, "format", "", "Format for the output (yaml, json, proto, go template..)")
	flagset.BoolVar(&opts.DryRun, "dryrun", false, "Run without actually modifying anything")
	flagset.BoolVar(&opts.ShowConfig, "show-config", false, "Print config for entity detail")
	flagset.BoolVarP(&opts.Interactive, "interactive", "i", false, "Enable interactive mode")
	flagset.StringToStringVar(&opts.Vars, "var", nil, "Override values for variables (--var VAR_NAME=value)")
}

func NewManageCmd(cli Cli) *cobra.Command {
	var opts ManageOptions
	cmd := &cobra.Command{
		Use:              "manage ENTITY [ACTION]",
		Short:            "Manage config with entities",
		Long:             "Manages initial config setup and updates to config using entities",
		Example:          color.Sprint(exampleManageCmd),
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
	opts.InstallFlags(cmd.PersistentFlags())
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
			color.Fprintf(cli.Out(), " - %s <gray>(%v vars)</>\n%s",
				color.Yellow.Sprint(name), len(e.GetVars()), desc)
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
		color.Fprintf(cli.Out(), "<lightWhite>ENTITY</>:      %s\n", color.Yellow.Sprint(entity.Name))
		if len(entity.Description) > 0 {
			desc := strings.TrimSpace(entity.Description)
			color.Fprintf(cli.Out(), "<lightWhite>DESCRIPTION</>: %s\n", color.FgDefault.Sprint(desc))
		}
		color.Fprintf(cli.Out(), "<lightWhite>ORIGIN</>: %s\n", color.Cyan.Sprint(entity.Origin))
		if entity.Single {
			color.Fprintf(cli.Out(), "<lightWhite>TYPE</>: %s\n", color.Blue.Sprint("single instance"))
		} else {
			color.Fprintf(cli.Out(), "<lightWhite>TYPE</>: %s\n", color.Blue.Sprint("multi instance"))
		}
		color.Fprintf(cli.Out(), "<lightWhite>VARS</>:\n")
		maxLen := getEntityVarMaxLen(entity) + 1
		for _, v := range entity.Vars {
			name := color.Style{color.Bold}.Sprintf("%-"+fmt.Sprint(maxLen)+"s", v.Name+":")
			valclr := color.Style{color.FgCyan}
			val := valclr.Sprint(v.Value)
			if v.Type != "" {
				val = fmt.Sprintf("%v (%v)", val, color.FgGray.Sprint(v.Type))
			}
			if v.Description != "" {
				desc := prefixTmpl(v.Description, "   ")
				color.Fprintf(cli.Out(), "  %v %v\n<gray>%v</>\n", color.Bold.Sprint(name), val, desc)
			} else {
				color.Fprintf(cli.Out(), "  %v %v\n", color.Bold.Sprint(name), val)
			}
		}
		if opts.ShowConfig {
			color.Fprintf(cli.Out(), "<lightWhite>CONFIG</>:\n%s\n", color.LightWhite.Sprint(prefixTmpl(entity.Config, "   ")))
		} else {
			configLen := len(entity.Config)
			configLines := strings.Count(entity.Config, "\n")
			color.Fprintf(cli.Out(), "<lightWhite>CONFIG</>: %s %s\n", color.White.Sprintf("%d lines", configLines), color.Gray.Sprintf("(%d bytes)", configLen))
		}
		color.Fprintf(cli.Out(), "<lightWhite>FILES</>:\n")
		for _, f := range entity.Files {
			contentLen := len(f.Content)
			contentLines := strings.Count(f.Content, "\n")
			color.Fprintf(cli.Out(), "  %s: %s %s\n", color.Bold.Sprint(f.Name), color.White.Sprintf("%d lines", contentLines), color.Gray.Sprintf("(%d bytes)", contentLen))
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

	// validate var overrides
	for k, v := range opts.Vars {
		if isIndexVar(k) {
			if entity.Single {
				return fmt.Errorf("single instance entity does not have internal var: %v", k)
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
			return fmt.Errorf("found override for undefined variable: %q", k)
		}
		logrus.Tracef("override for var %v: %v", k, v)
	}

	if opts.Count > 1 && entity.Single {
		return fmt.Errorf("count must be 1 for entity with single instance")
	}
	if len(entity.Files) > 0 && !entity.Single {
		return fmt.Errorf("no files can be defined for multi instance entity")
	}

	// prepare config
	mainConf, err := client.NewDynamicConfig(allModels())
	if err != nil {
		return fmt.Errorf("failed to create dynamic config for all models")
	}

	logrus.Tracef("generating config (count: %d, offset: %d)", opts.Count, opts.IdOffset)

	// repeat for given count
	for i := 0; i < int(opts.Count); i++ {
		idx := i + opts.IdOffset
		id := idx + 1

		// prepare vars
		evars := map[string]string{}
		if !entity.Single {
			evars[varIDX] = fmt.Sprint(idx)
			evars[varID] = fmt.Sprint(id)
		}
		for k, v := range opts.Vars {
			evars[k] = v
		}

		var vars map[string]string
		// render values
		if opts.Interactive {
			color.Fprintf(cli.Out(), "=> <fg=blue;op=bold>Set entity var values for instance</> <bold>ID</>=<white>%v</>\n\n", id)
			vars, err = prepareVarValuesInteractive(cli.Out(), entity, evars)
			if err != nil {
				return fmt.Errorf("failed to render vars (idx: %v): %w", i, err)
			}
		} else {
			vars, err = prepareVarValues(entity, evars)
			if err != nil {
				return fmt.Errorf("failed to render vars (idx: %v): %w", i, err)
			}
		}

		// render config template
		rawConf, err := renderEntityTemplate(entity.Config, vars)
		if err != nil {
			return fmt.Errorf("failed to render config (idx: %v): %w", i, err)
		}

		logrus.Tracef(" - [#%d] final vars:\n%s\nraw config:\n%v", i, yamlTmpl(vars), rawConf)

		// TODO: find some way to preserve the YAML comments from the config template
		//playWithYaml(config)

		// unmarshal into proto config
		conf := mainConf.New().Interface()
		bj, err := yaml2.YAMLToJSON([]byte(rawConf))
		if err != nil {
			return fmt.Errorf("cannot convert YAML to JSON: %w", err)
		}
		err = protojson.Unmarshal(bj, conf)
		if err != nil {
			return fmt.Errorf("cannot unmarshall data into dynamic config due to: %w", err)
		}
		logrus.Tracef("rendered config #%d:\n%v", i, yamlTmpl(conf))

		if err := mergeConfigs(mainConf, conf); err != nil {
			return fmt.Errorf("merging configs failed: %w", err)
		}

		for _, f := range entity.Files {
			rawData, err := renderEntityTemplate(f.Content, vars)
			if err != nil {
				return fmt.Errorf("failed to render file (%v): %w", f.Name, err)
			}

			logrus.Tracef(" - final vars:\n%s\nraw file %s:\n%v", yamlTmpl(vars), f.Name, rawData)

			if err := os.WriteFile(f.Name, []byte(rawData), 0666); err != nil {
				return fmt.Errorf("failed to write file (%v): %w", f.Name, err)
			}
		}
	}

	var finalConf protoreflect.ProtoMessage

	if opts.Target != "" {
		// load config from target file
		config, err := os.ReadFile(opts.Target)
		if err != nil {
			return fmt.Errorf("failed to read target config file: %w", err)
		}

		// unmarshal into proto config
		targetConf := mainConf.New().Interface()
		bj, err := yaml2.YAMLToJSON(config)
		if err != nil {
			return fmt.Errorf("cannot convert YAML to JSON: %w", err)
		}
		err = protojson.Unmarshal(bj, targetConf)
		if err != nil {
			return fmt.Errorf("cannot unmarshall data into dynamic config due to: %w", err)
		}
		logrus.Tracef("target config:\n%s", yamlTmpl(targetConf))

		if action == "ADD" {
			// merge with main with target config
			if err := mergeConfigs(targetConf, mainConf); err != nil {
				return fmt.Errorf("merging configs failed: %w", err)
			}

			logrus.Tracef("merged config:\n%s", yamlTmpl(targetConf))

			// extract items
			items, err := client.DynamicConfigExport(targetConf.(*dynamicpb.Message))
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

		finalConf = targetConf
	} else {
		finalConf = mainConf
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

func mergeConfigs(dst proto.Message, src proto.Message) error {
	dkeys := extractModelKeysFromConfig(dst)
	skeys := extractModelKeysFromConfig(src)

	// check for any conflicts before merging
	conflictKeys := findConflictingKeys(dkeys, skeys)
	if len(conflictKeys) > 0 {
		logrus.Tracef("listing %d conflicting keys:", len(conflictKeys))
		for _, c := range conflictKeys {
			logrus.Tracef(" - conflict key: %v", c)
		}
		return fmt.Errorf("found %d conflicting keys", len(conflictKeys))
	}

	logrus.Tracef("merging configs (dkeys: %d, skeys: %d)", len(dkeys), len(skeys))
	proto.Merge(dst, src)

	return nil
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
	var keys []string
	for i, item := range items {
		l := logrus.WithFields(map[string]interface{}{
			"item":      fmt.Sprintf("%d/%d", i+1, len(items)),
			"protoName": item.ProtoReflect().Descriptor().FullName(),
		})
		model, err := models.GetModelFor(item)
		if err != nil {
			l.Tracef("no model found for item: %v", item)
			continue
		}
		l = l.WithField("model", model.Name())
		name, err := model.InstanceName(item)
		if err != nil {
			l.Tracef("instance name error: %v", err)
			continue
		}
		if name == "" {
			l.Tracef("intance has empty name, skipping item")
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

func renderEntityTemplate(cfg string, evars map[string]string) (string, error) {
	tmpl, err := interpolateStr(cfg, evars)
	if err != nil {
		return "", err
	}
	config, err := renderTmpl(tmpl, evars)
	if err != nil {
		return "", err
	}
	return config, nil
}

func prepareVarValues(e Entity, evars map[string]string) (map[string]string, error) {
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

func prepareVarValuesInteractive(w io.Writer, e Entity, evars map[string]string) (map[string]string, error) {
	for _, v := range e.Vars {
		if v.When != "" {
			tmpl, err := interpolateStr(v.When, evars)
			if err != nil {
				return nil, err
			}
			res, err := renderTmpl(tmpl, evars)
			if err != nil {
				return nil, err
			}
			ok, err := strconv.ParseBool(res)
			if err != nil {
				logrus.Tracef("parse bool err: %v", err)
				continue
			}
			logrus.Tracef("when %s returned %q (%v)", v.Name, res, ok)
			if !ok {
				continue
			}
		}
		vv := v.Value
		if ov, ok := evars[v.Name]; ok {
			vv = ov
		}

		color.Fprintf(w, "-> <lightBlue>Set value for:</> %v <gray>(press ENTER to confirm)</>\n", color.White.Sprint(v.Name))
		if v.Description != "" {
			desc := prefixTmpl(v.Description, "   ")
			color.Fprintln(w, color.Gray.Sprint(desc))
		}
		fmt.Fprintln(w)
		fmt.Fprintf(w, "    ")
		cval, err := promptUserValue(v.Name, vv)
		if err != nil {
			return nil, err
		} else {
			if vv != cval {
				vv = cval
			}
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

		color.Fprintf(w, "%s=%s\n", color.Cyan.Sprint(v.Name), color.LightGreen.Sprint(val))
		fmt.Fprintln(w)
	}
	return evars, nil
}

func promptUserValue(label string, defval string) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		Default:   defval,
		AllowEdit: true,
		Pointer:   promptui.PipeCursor,
		Templates: &promptui.PromptTemplates{
			Prompt: color.Sprintf("    <bold>{{ . }}:</b> "),
		},
		// TODO: validate values
	}
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
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
