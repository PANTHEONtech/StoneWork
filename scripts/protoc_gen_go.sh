#!/usr/bin/env bash

# SPDX-License-Identifier: Apache-2.0

# Copyright 2022 PANTHEON.tech
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eo pipefail

fail() {
  echo "$@" >&2
  exit 1
}

if [ -z "${1}" ] || [ -z "${2}" ]; then
  fail "usage: ${0} <proto_path> <protoc_gen_go_out> [protoc_gen_go_parameter]"
fi

protoc --version

PROTO_PATH="${1}"
PROTOC_GEN_GO_OUT="${2}"
PROTOC_GEN_GO_PARAMETER="${3}"

PROTO_PATH_LIGATO="$( go list -f '{{.Dir}}' -m go.ligato.io/vpp-agent/v3 )/proto"

PROTOC_GEN_GO_ARGS="${PROTOC_GEN_GO_OUT}"
if [ -n "${PROTOC_GEN_GO_PARAMETER}" ]; then
  PROTOC_GEN_GO_ARGS="${PROTOC_GEN_GO_PARAMETER}:${PROTOC_GEN_GO_ARGS}"
fi

mkdir -p "${PROTOC_GEN_GO_OUT}"

# all directories with proto files
protodirs=$(find "${PROTO_PATH}" -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $protodirs; do
	  protofiles=$(find "${dir}" -maxdepth 1 -name '*.proto')
	  echo "$dir | $protofiles"

    protoc --proto_path="${PROTO_PATH}" \
  	  --proto_path="${PROTO_PATH_LIGATO}" \
  		--go_out="${PROTOC_GEN_GO_ARGS}" \
  		--go-grpc_out="${PROTOC_GEN_GO_ARGS}" $protofiles
done
