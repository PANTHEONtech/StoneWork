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
	"os"
	"strconv"

	"go.ligato.io/vpp-agent/v3/pkg/models"
)

var (
	ModelMockCnf1 = models.Register(&MockCnf1{}, models.Spec{
		Module:  "mock1",
		Version: "v1",
		Type:    "mock-type",
	}, models.WithNameTemplate("{{.IpProtocol}}"))

	ModelMockCnf2 = models.Register(&MockCnf2{}, models.Spec{
		Module:  "mock2",
		Version: "v1",
		Type:    "mock-type",
	}, models.WithNameTemplate("{{.VppInterface}}"))
)

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
