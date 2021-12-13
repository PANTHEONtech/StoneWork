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

const protoTemplate = `// Proto file with the configuration model of {{ .CnfName }}.
syntax = "proto3";

package {{ lower .CnfName }};

option go_package = "pantheon.tech/{{ lower .CnfName }}/proto/{{ lower .CnfName }};{{ lower .CnfName }}";

{{ range $i, $path := .Imports }}
import "{{ $path }}";
{{- end }}

// Configuration root wrapping all models supported by {{ .CnfName }}.
message Root {
{{- range $i, $group := .ModelGroups }}
    message {{ firstUpper $group.Name }} {
    {{- range $j, $model := $group.Models }}
    {{- if .Repeated }}
        repeated {{ $model.ProtoMessage }} {{ $model.Name }} = {{ inc $j }};
    {{- else }}
        {{ $model.ProtoMessage }} {{ $model.Name }} = {{ inc $j }};
    {{- end }}
    {{- end }}
    }
    {{ firstUpper $group.Name }} {{ $group.Name }} = {{ inc $i }};
{{- end }}
}
`
