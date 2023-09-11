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
	"context"
	"fmt"
	"io/fs"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"go.ligato.io/cn-infra/v2/agent"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	"go.ligato.io/vpp-agent/v3/pkg/util"
	"go.ligato.io/vpp-agent/v3/proto/ligato/generic"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"go.pantheon.tech/stonework/cmd/stonework/app"
)

const (
	protoFileExt     = ".proto"
	ligatoApiFileDir = "/api"
	cnfApiFileDir    = "/cnfapi"

	modelSpecExtNum = 50222
	modelTmplExtNum = 50223

	logo = `
   ______               _      __         __
  / __/ /____  ___  ___| | /| / /__  ____/ /__
 _\ \/ __/ _ \/ _ \/ -_) |/ |/ / _ \/ __/  '_/
/___/\__/\___/_//_/\__/|__/|__/\___/_/ /_/\_\  %s
`
)

func findProtoFiles(dir string) []string {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(d.Name(), protoFileExt) {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		// TODO: better error handling here
		logging.DefaultLogger.Errorf("error when searching for proto files: %w", err)
		return nil
	}
	return files
}

func compileProtoFiles(files []string, importPaths ...string) (linker.Files, error) {
	resolver := protocompile.WithStandardImports(&protocompile.SourceResolver{
		ImportPaths: importPaths,
	})
	compiler := protocompile.Compiler{
		Resolver:       resolver,
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	return compiler.Compile(context.Background(), files...)
}

func extractMessageOptions(files linker.Files) map[protowire.Number]protoreflect.ExtensionTypeDescriptor {
	result := make(map[protowire.Number]protoreflect.ExtensionTypeDescriptor)
	optsMsgName := (&descriptorpb.MessageOptions{}).ProtoReflect().Descriptor().FullName()
	for _, f := range files {
		exts := f.Extensions()
		for i := 0; i < exts.Len(); i++ {
			fieldNum := exts.Get(i).Number()
			extTypeDesc := f.FindExtensionByNumber(optsMsgName, fieldNum)
			if extTypeDesc != nil {
				result[fieldNum] = extTypeDesc
			}
		}
	}
	return result
}

func registerModels(files linker.Files, specExtDesc, tmplExtDesc protoreflect.ExtensionTypeDescriptor) error {
	for _, f := range files {
		msgs := f.Messages()
		for i := 0; i < msgs.Len(); i++ {
			msgDesc := msgs.Get(i)
			opts := msgDesc.Options().(*descriptorpb.MessageOptions)
			if !proto.HasExtension(opts, specExtDesc.Type()) {
				// this msg is not a model, continue with next message
				continue
			}
			dynSpec, ok := proto.GetExtension(opts, specExtDesc.Type()).(*dynamicpb.Message)
			if !ok {
				// no model spec detected, continue with next message
				continue
			}
			specMsg, err := util.ConvertProto(&generic.ModelSpec{}, dynSpec)
			if err != nil {
				return err
			}
			var modelOpts []models.ModelOption
			if proto.HasExtension(opts, tmplExtDesc.Type()) {
				tmpl := proto.GetExtension(opts, tmplExtDesc.Type()).(string)
				modelOpts = append(modelOpts, models.WithNameTemplate(tmpl))
			}
			spec := models.ToSpec(specMsg.(*generic.ModelSpec))
			modelName := spec.ModelName()
			modelMsg := dynamicpb.NewMessage(msgDesc)
			if _, err := models.DefaultRegistry.GetModel(modelName); err == nil {
				// model already registered, print warning and continue with next message
				logging.DefaultLogger.Warnf("Cannot register cnf model: model with name %s is already registered", modelName)
				continue
			}
			_, err = models.DefaultRegistry.Register(modelMsg, spec, modelOpts...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	if _, err := os.Stat(cnfApiFileDir); !os.IsNotExist(err) {
		ligatoApiFiles, err := compileProtoFiles(findProtoFiles(ligatoApiFileDir), ligatoApiFileDir)
		if err != nil {
			logging.DefaultLogger.Fatal(err)
		}
		files, err := compileProtoFiles(findProtoFiles(cnfApiFileDir), cnfApiFileDir, ligatoApiFileDir)
		if err != nil {
			logging.DefaultLogger.Fatal(err)
		}
		ligatoMsgExtDescs := extractMessageOptions(ligatoApiFiles)
		modelSpecExtDesc := ligatoMsgExtDescs[modelSpecExtNum]
		modelTmplExtDesc := ligatoMsgExtDescs[modelTmplExtNum]
		if err = registerModels(files, modelSpecExtDesc, modelTmplExtDesc); err != nil {
			logging.DefaultLogger.Fatal(err)
		}
	}
	if _, err := fmt.Fprintf(os.Stdout, logo, agent.BuildVersion); err != nil {
		logging.DefaultLogger.Fatal(err)
	}

	swAgent := app.New()
	a := agent.NewAgent(agent.AllPlugins(swAgent), agent.StartTimeout(30*time.Second))
	if err := a.Run(); err != nil {
		logging.DefaultLogger.Fatal(err)
	}
}
