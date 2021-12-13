#!/usr/bin/env bash

# SPDX-License-Identifier: Apache-2.0

# Copyright 2021 PANTHEON.tech
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

set -euo pipefail

PROTOC_VERSION="3.13.0"
PROTOC_OS="linux"
PROTOC_ARCH="x86_64"
PROTOC_GEN_GO_VERSION="1.25.0"
PROTOC_GEN_GO_GRPC_COMMIT="ad51f572fd270f2323e3aa2c1d2775cab9087af2" # before first version tagging

PROTO_DIR=${1:-proto}
ROOT_DIR=$(pwd)

# install correct version of protoc
PROTOC_DIR="/tmp/cached-protoc-"${PROTOC_VERSION}
PROTOC_BIN_DIR=${PROTOC_DIR}"/bin"
if [ ! -d ${PROTOC_DIR} ]
then
  mkdir -p ${PROTOC_DIR}
  cd ${PROTOC_DIR}
  curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-${PROTOC_OS}-${PROTOC_ARCH}.zip -o protoc.zip
  unzip -o protoc.zip
fi
echo "using "`${PROTOC_BIN_DIR}"/protoc" --version`

# install correct version of go plugin for protoc (can't detect version of installed -> install again)
TMP_DIR=$(mktemp -d) # for go >= 1.16 we can replace this with "go install package@version"
cd ${TMP_DIR}
GO111MODULE=on go get google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}
rm -rf ${TMP_DIR}

# install correct version of go-grpc plugin for protoc (can't detect version of installed -> install again)
TMP_DIR=$(mktemp -d)
cd ${TMP_DIR}
GO111MODULE=on go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_COMMIT}
rm -rf ${TMP_DIR}

# get stonework protos
cd ${ROOT_DIR}
protos=$(find $PROTO_DIR -type f -name '*.proto')

# Set parent directory to hold all the symlinks
TMP_DEP_DIR=$(mktemp -d -t stonework-deps-XXXXXX)
mkdir -p "${TMP_DEP_DIR}"
trap "{ rm -rf $TMP_DEP_DIR; }" EXIT

# Download all the required dependencies
go mod download

# Get all the modules we use and create required directory structure
go list -f "${TMP_DEP_DIR}/{{ .Path }}" -m all \
  | xargs -L1 dirname | sort | uniq | xargs mkdir -p

# Create symlinks
go list -f "{{ .Dir }} ${TMP_DEP_DIR}/{{ .Path }}" -m all \
  | xargs -L1 -- ln -s

# fix for vpp-agent protos that reference some other vpp-agent protos with relative path instead of absolute
ln -s go.ligato.io/vpp-agent/v3/proto/ligato ${TMP_DEP_DIR}/ligato

for proto in $protos; do
	echo " - $proto";
	${PROTOC_BIN_DIR}"/protoc" \
		-I ${TMP_DEP_DIR} \
		--proto_path=${PROTO_DIR} \
		--go_out=paths=source_relative:${PROTO_DIR} "$proto" \
		--go-grpc_out=paths=source_relative:${PROTO_DIR} "$proto";
done

rm -rf ${TMP_DEP_DIR}
