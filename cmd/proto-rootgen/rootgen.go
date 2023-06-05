// SPDX-License-Identifier: Apache-2.0

// Copyright 2021 PANTHEON.tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/namsral/flag"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	"go.ligato.io/vpp-agent/v3/proto/ligato/generic"
	"google.golang.org/protobuf/encoding/protojson"
)

// TemplateData encapsulates input arguments for the template.
type TemplateData struct {
	Imports     []string
	CnfName     string
	ModelGroups []*ModelGroup
}

type ModelGroup struct {
	Name   string
	Models []*Model
}

type Model struct {
	Repeated     bool
	Name         string
	ProtoMessage string
}

const (
	modelsSpecFileName  = "models.spec.yaml"
	groupSuffix         = "Config"
	repeatedFieldSuffix = "_list"
)

var (
	cnfName = flag.String("cnf-name", "StoneWork", "Name of the CNF for which the root proto file will be generated")
	apiDir  = flag.String("api-dir", "/api", "Directory with CNF API definitions (models.spec.yaml and proto files)")
)

func main() {
	flag.Parse()
	// try to read and parse models.spec.yaml
	var cnfModels []*generic.ModelDetail
	modelsSpecFilePath := path.Join(*apiDir, modelsSpecFileName)
	specData, err := ioutil.ReadFile(modelsSpecFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR ReadFile: ", err)
		os.Exit(1)
	}
	for _, yamlDoc := range bytes.Split(specData, []byte("---")) {
		yamlDoc = bytes.TrimSpace(yamlDoc)
		if len(yamlDoc) == 0 {
			continue
		}
		cnfModel := &generic.ModelDetail{}
		jsonDoc, err := yaml.YAMLToJSON(yamlDoc)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR YAMLToJSON: ", err)
			os.Exit(1)
		}
		err = protojson.Unmarshal(jsonDoc, cnfModel)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR Unmarshal: ", err)
			os.Exit(1)
		}
		cnfModels = append(cnfModels, cnfModel)
	}

	// prepare input data for the template
	inputData := TemplateData{
		CnfName: *cnfName,
	}
	groups := make(map[string]*ModelGroup)
	imports := make(map[string]struct{})
	for _, cnfModel := range cnfModels {
		model := &Model{
			ProtoMessage: cnfModel.ProtoName,
		}
		if protoFile, err := models.ModelOptionFor("protoFile", cnfModel.Options); err == nil {
			imports[protoFile] = struct{}{}
		}
		spec := models.ToSpec(cnfModel.Spec)
		groupName := fmt.Sprintf("%v%v", modulePrefix(spec.ModelName()), groupSuffix)
		model.Name = simpleProtoName(cnfModel.ProtoName)
		if _, err = models.ModelOptionFor("nameTemplate", cnfModel.Options); err == nil {
			model.Name += repeatedFieldSuffix
			model.Repeated = true
		}
		key := fmt.Sprintf("%v.%v", groupName, simpleProtoName(cnfModel.ProtoName))
		if name, found := backwardCompatibleNames[key]; found {
			// using field name from hardcoded ligato.configurator.Config to achieve json/yaml backward
			// compatibility with ligato
			model.Name = name
		}
		if _, groupFound := groups[groupName]; !groupFound {
			groups[groupName] = &ModelGroup{
				Name: groupName,
			}
		}
		group := groups[groupName]
		group.Models = append(group.Models, model)
	}
	for _, group := range groups {
		sort.Slice(group.Models, func(i, j int) bool {
			return group.Models[i].Name < group.Models[j].Name
		})
		inputData.ModelGroups = append(inputData.ModelGroups, group)
	}
	sort.Slice(inputData.ModelGroups, func(i, j int) bool {
		return inputData.ModelGroups[i].Name < inputData.ModelGroups[j].Name
	})
	for protoImport := range imports {
		inputData.Imports = append(inputData.Imports, protoImport)
	}
	sort.Slice(inputData.Imports, func(i, j int) bool {
		return inputData.Imports[i] < inputData.Imports[j]
	})

	// Parse template for the proto file
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"firstUpper": func(s string) string {
			if len(s) > 0 {
				return strings.ToUpper(s[:1]) + s[1:]
			}
			return ""
		},
	}
	var buf bytes.Buffer
	t := template.Must(template.New("proto-root").Funcs(funcMap).Parse(protoTemplate))
	err = t.Execute(&buf, inputData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR template: ", err)
		os.Exit(2)
	}

	// output the generated proto file
	protoFile := path.Join(*apiDir, strings.ToLower(*cnfName)+"-root.proto")
	err = ioutil.WriteFile(protoFile, buf.Bytes(), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR WriteFile: ", err)
		os.Exit(4)
	}
}

func modulePrefix(modelName string) string {
	return strings.Split(modelName, ".")[0]
}

func simpleProtoName(fullProtoName string) string {
	nameSplit := strings.Split(fullProtoName, ".")
	return nameSplit[len(nameSplit)-1]
}
