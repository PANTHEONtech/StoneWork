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

UNAME_OS   ?= $(shell uname -s)
UNAME_ARCH ?= $(shell uname -m)

ifndef CACHE_BASE
CACHE_BASE := $(HOME)/.cache/$(PROJECT)
endif
CACHE := $(CACHE_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
CACHE_BIN := $(CACHE)/bin
CACHE_INCLUDE := $(CACHE)/include
CACHE_VERSIONS := $(CACHE)/versions

export PATH := $(abspath $(CACHE_BIN)):$(PATH)

# https://github.com/protocolbuffers/protobuf-go
PROTOC_GEN_GO_VERSION ?= v1.27.1
# https://github.com/grpc/grpc-go
PROTOC_GEN_GO_GRPC_VERSION ?= v1.38.0
# https://github.com/protocolbuffers/protobuf/releases
PROTOC_VERSION ?= 3.17.3

GO_BINS := $(GO_BINS) \
	buf \
	protoc-gen-buf-check-breaking \
	protoc-gen-buf-check-lint

PROTO_PATH := proto
PROTOC_GEN_GO_OUT := proto

PROTOC_GEN_GO_PARAMETER ?= paths=source_relative

ifeq ($(UNAME_OS),Darwin)
PROTOC_OS := osx
PROTOC_ARCH := x86_64
endif
ifeq ($(UNAME_OS),Linux)
PROTOC_OS = linux
PROTOC_ARCH := $(UNAME_ARCH)
endif

PROTOC := $(CACHE_VERSIONS)/protoc/$(PROTOC_VERSION)
$(PROTOC):
	@if ! command -v curl >/dev/null 2>/dev/null; then echo "error: curl must be installed"  >&2; exit 1; fi
	@if ! command -v unzip >/dev/null 2>/dev/null; then echo "error: unzip must be installed"  >&2; exit 1; fi
	@rm -f $(CACHE_BIN)/protoc
	@rm -rf $(CACHE_INCLUDE)/google
	@mkdir -p $(CACHE_BIN) $(CACHE_INCLUDE)
	$(eval PROTOC_TMP := $(shell mktemp -d))
	cd $(PROTOC_TMP); curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-$(PROTOC_ARCH).zip -o protoc.zip
	cd $(PROTOC_TMP); unzip protoc.zip && mv bin/protoc $(CACHE_BIN)/protoc && mv include/google $(CACHE_INCLUDE)/google
	@rm -rf $(PROTOC_TMP)
	@rm -rf $(dir $(PROTOC))
	@mkdir -p $(dir $(PROTOC))
	@touch $(PROTOC)

PROTOC_GEN_GO := $(CACHE_VERSIONS)/protoc-gen-go/$(PROTOC_GEN_GO_VERSION)
$(PROTOC_GEN_GO):
	@rm -f $(GOBIN)/protoc-gen-go
	$(eval PROTOC_GEN_GO_TMP := $(shell mktemp -d))
	cd $(PROTOC_GEN_GO_TMP); go get google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@rm -rf $(PROTOC_GEN_GO_TMP)
	@rm -rf $(dir $(PROTOC_GEN_GO))
	@mkdir -p $(dir $(PROTOC_GEN_GO))
	@touch $(PROTOC_GEN_GO)

PROTOC_GEN_GO_GRPC := $(CACHE_VERSIONS)/protoc-gen-go-grpc/$(PROTOC_GEN_GO_GRPC_VERSION)
$(PROTOC_GEN_GO_GRPC):
	@if ! command -v git >/dev/null 2>/dev/null; then echo "error: git must be installed"  >&2; exit 1; fi
	@rm -f $(GOBIN)/protoc-gen-go-grpc
	$(eval PROTOC_GEN_GO_GRPC_TMP := $(shell mktemp -d))
	#cd $(PROTOC_GEN_GO_GRPC_TMP); go get -u -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	cd $(PROTOC_GEN_GO_GRPC_TMP); git clone -b $(PROTOC_GEN_GO_GRPC_VERSION) https://github.com/grpc/grpc-go
	cd $(PROTOC_GEN_GO_GRPC_TMP); cd grpc-go/cmd/protoc-gen-go-grpc && go install .
	@rm -rf $(PROTOC_GEN_GO_GRPC_TMP)
	@rm -rf $(dir $(PROTOC_GEN_GO_GRPC))
	@mkdir -p $(dir $(PROTOC_GEN_GO_GRPC))
	@touch $(PROTOC_GEN_GO_GRPC)

.PHONY: protocgengoclean
protocgengoclean:
	find "$(PROTOC_GEN_GO_OUT)" -name "*.pb.go" -exec rm -drfv '{}' \;

.PHONY: protocgengo
protocgengo: protocgengoclean $(PROTOC) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)
	bash scripts/protoc_gen_go.sh "$(PROTO_PATH)" "$(PROTOC_GEN_GO_OUT)" "$(PROTOC_GEN_GO_PARAMETER)"
