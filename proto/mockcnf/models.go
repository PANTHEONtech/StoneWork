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

package mockcnf

import (
	"fmt"
	"os"
	"strconv"

	"github.com/golang/protobuf/proto"
	"go.ligato.io/vpp-agent/v3/pkg/models"
)

var (
	ModelMockCnf = models.Register(getModelProto(), models.Spec{
		Module:  fmt.Sprintf("mock%d", MockCnfIndex()),
		Version: "v1",
		Type:    "mock-type",
	}, models.WithNameTemplate(getModelNameTemplate()))
)

func getModelProto() proto.Message {
	switch MockCnfIndex() {
	case 1:
		return &MockCnf1{}
	case 2:
		return &MockCnf2{}
	}
	return nil
}

func getModelNameTemplate() string {
	switch MockCnfIndex() {
	case 1:
		return "{{.IpProtocol}}"
	case 2:
		return "{{.VppInterface}}"
	}
	return ""
}

func MockCnfIndex() int {
	index, err := strconv.Atoi(os.Getenv("MOCK_CNF_INDEX"))
	if err != nil {
		panic(err.Error())
	}
	if index < 1 || index > 2 {
		panic("mock cnf index out of range")
	}
	return index
}
