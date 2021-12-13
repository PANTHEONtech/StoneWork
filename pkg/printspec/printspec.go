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

package printspec

import (
	"io"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"

	"go.ligato.io/vpp-agent/v3/pkg/models"
)

// Print specification of all known models.
func AllKnownModels(writer io.Writer) error {
	return SelectedModels(writer, models.DefaultRegistry.RegisteredModels()...)
}

// Print specification of only selected subset of models.
func SelectedModels(writer io.Writer, models ...models.KnownModel) error {
	for i, model := range models {
		if model.Spec().Class != "config" {
			continue
		}
		m := jsonpb.Marshaler{Indent: "  "}
		jsonOut, err := m.MarshalToString(model.ModelDetail())
		if err != nil {
			return err
		}
		yamlOut, err := yaml.JSONToYAML([]byte(jsonOut))
		if err != nil {
			return err
		}
		_, err = writer.Write(append([]byte("---\n"), yamlOut...))
		if err != nil {
			return err
		}
		if i < len(models)-1 {
			_, err = writer.Write([]byte("\n\n"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
