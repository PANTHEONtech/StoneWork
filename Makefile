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

SHELL := /usr/bin/env bash -o pipefail

PROJECT    := StoneWork
VERSION    ?= $(shell git describe --always --tags --dirty)
COMMIT     ?= $(shell git rev-parse HEAD)
BRANCH     ?= $(shell git rev-parse --abbrev-ref HEAD)
DATE       ?= $(shell git log -1 --format="%ct" | xargs -I{} date -d @{} +'%Y-%m-%dT%H:%M%:z')
BUILD_DATE ?= $(shell date +%s)
BUILD_HOST ?= $(shell hostname)
BUILD_USER ?= $(shell id -un)

LDFLAGS = -w -s

CNINFRA := go.ligato.io/cn-infra/v2/agent
LDFLAGS += \
	-X $(CNINFRA).BuildVersion=$(VERSION) \
	-X $(CNINFRA).CommitHash=$(COMMIT) \
	-X $(CNINFRA).BuildDate=$(DATE)

VERSION_PKG := go.pantheon.tech/stonework/pkg/version
LDFLAGS += \
	-X $(VERSION_PKG).version=$(VERSION) \
	-X $(VERSION_PKG).commit=$(COMMIT)-$(DATE) \
	-X $(VERSION_PKG).branch=$(BRANCH) \
	-X $(VERSION_PKG).buildStamp=$(BUILD_DATE) \
	-X $(VERSION_PKG).buildHost=$(BUILD_HOST) \
	-X $(VERSION_PKG).buildUser=$(BUILD_USER)

RELEASE_TAG ?= $(shell git describe --always --tags --dirty --exact-match 2>/dev/null)

TAG_FORMAT="^v([0-9]+\.){2}[0-9]+.*"
RELEASE_TAG_CHECKED = $(shell echo $(RELEASE_TAG) | grep -v "\-dirty" | grep -E ${TAG_FORMAT})
ifneq ($(RELEASE_TAG_CHECKED),)
RELEASE_VERSION_FULL = $(shell echo $(RELEASE_TAG_CHECKED) | cut -c 2-)
RELEASE_VERSION_MAJOR_MINOR = $(shell echo $(RELEASE_VERSION_FULL) | cut -d '.' -f 1-2)
endif

ifeq ($(VPP_VERSION),)
VPP_VERSION="23.06"
endif
ifeq ($(DEV_VERSION),) # for tagging in-development images
DEV_VERSION="23.06"
endif
REPO="ghcr.io/pantheontech"
STONEWORK_VPP_IMAGE="stonework-vpp"
STONEWORK_VPP_TEST_IMAGE="stonework-vpp-test"
STONEWORK_DEV_IMAGE="stonework-dev"
STONEWORK_PROD_IMAGE="stonework"
TESTER_IMAGE="stonework-tester"
MOCK_CNF_IMAGE="stonework-mockcnf"
PROTO_ROOTGEN_IMAGE="stonework-proto-rootgen"

export DOCKER_BUILDKIT=1

help:
	@echo "List of make targets:"
	@grep -E '^[a-zA-Z_-]+:.*?(## .*)?$$' $(MAKEFILE_LIST) | sed 's/^[^:]*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT = help

-include scripts/make/proto.make

build:
	@cd cmd/stonework && go build -v -ldflags "${LDFLAGS}"
	@cd cmd/swctl && go build -v -ldflags "${LDFLAGS}"
	@cd cmd/stonework-init && go build -v -ldflags "${LDFLAGS}"
	@cd cmd/mockcnf && go build -v -ldflags "${LDFLAGS}"
	@cd cmd/proto-rootgen && go build -v -ldflags "${LDFLAGS}"

install:
	go install -v -ldflags "${LDFLAGS}" ./cmd/stonework
	go install -v -ldflags "${LDFLAGS}" ./cmd/swctl
	go install -v -ldflags "${LDFLAGS}" ./cmd/stonework-init

install-mockcnf:
	@cd cmd/mockcnf && go install -v -ldflags "${LDFLAGS}"

install-proto-rootgen:
	@cd cmd/proto-rootgen && go install -v -ldflags "${LDFLAGS}"

install-swctl:
	@cd cmd/swctl && go install -v -ldflags "$(LDFLAGS)"

# -------------------------------
#  Images
# -------------------------------

vpp-image:
	@echo "=> building VPP image"
	VPP_VERSION=${VPP_VERSION} \
	IMAGE_TAG="${STONEWORK_VPP_IMAGE}:${VPP_VERSION}" \
	./scripts/build.sh vpp

vpp-test-image:
	@echo "=> building VPP test image"
	VPP_IMAGE="${STONEWORK_VPP_IMAGE}:${VPP_VERSION}" \
	IMAGE_TAG="${STONEWORK_VPP_TEST_IMAGE}:${VPP_VERSION}" \
	./scripts/build.sh vpp-test

dev-image:
	@echo "=> building development image"
	VPP_IMAGE="${STONEWORK_VPP_IMAGE}:${VPP_VERSION}" \
	VPP_VERSION=${VPP_VERSION} \
	IMAGE_TAG="${STONEWORK_DEV_IMAGE}:${DEV_VERSION}" \
	./scripts/build.sh dev

prod-image:
	@echo "=> building production image"
	VPP_IMAGE="${STONEWORK_VPP_IMAGE}:${VPP_VERSION}" \
	DEV_IMAGE="${STONEWORK_DEV_IMAGE}:${DEV_VERSION}" \
	IMAGE_TAG="${STONEWORK_PROD_IMAGE}:${DEV_VERSION}" \
	./scripts/build.sh prod

tester-image:
	@echo "=> building image with network tools for testing"
	IMAGE_TAG="${TESTER_IMAGE}:${DEV_VERSION}" \
	./scripts/build.sh tester

mockcnf-image:
	@echo "=> building mock CNF"
	VPP_IMAGE="${STONEWORK_VPP_IMAGE}:${VPP_VERSION}" \
	IMAGE_TAG="${MOCK_CNF_IMAGE}:${DEV_VERSION}" \
	./scripts/build.sh mockcnf

proto-rootgen-image:
	@echo "=> building image for building proto file with the config root message"
	IMAGE_TAG="${PROTO_ROOTGEN_IMAGE}:${DEV_VERSION}" \
	./scripts/build.sh proto-rootgen

images: vpp-image dev-image prod-image tester-image mockcnf-image
	docker tag ${STONEWORK_VPP_IMAGE}:${VPP_VERSION} ${REPO}/${STONEWORK_VPP_IMAGE}:${VPP_VERSION}

ifneq ($(RELEASE_TAG_CHECKED),)
	# tag release images
	docker tag ${STONEWORK_PROD_IMAGE}:${DEV_VERSION} ${REPO}/${STONEWORK_PROD_IMAGE}:${RELEASE_VERSION_FULL}
	docker tag ${STONEWORK_PROD_IMAGE}:${DEV_VERSION} ${REPO}/${STONEWORK_PROD_IMAGE}:${RELEASE_VERSION_MAJOR_MINOR}
	docker tag ${STONEWORK_PROD_IMAGE}:${DEV_VERSION} ${REPO}/${STONEWORK_PROD_IMAGE}
endif

push-images:
	docker push ${REPO}/${STONEWORK_VPP_IMAGE}:${VPP_VERSION}

ifneq ($(RELEASE_TAG_CHECKED),)
	docker push ${REPO}/${STONEWORK_PROD_IMAGE}:${RELEASE_VERSION_FULL}
ifneq ($(findstring -,$(RELEASE_TAG_CHECKED)),-)
	docker push ${REPO}/${STONEWORK_PROD_IMAGE}:${RELEASE_VERSION_MAJOR_MINOR}
	docker push ${REPO}/${STONEWORK_PROD_IMAGE}
endif
else
	@echo "Release tag is empty or has incorrect format."
	@echo "Supplied release tag: ${RELEASE_TAG}"
	@echo 'Expected format: ${TAG_FORMAT} ; must not contain "-dirty"'
	@false
endif

cleanup-images:
	docker rmi ${STONEWORK_DEV_IMAGE}:${DEV_VERSION} || \:
	docker rmi ${MOCK_CNF_IMAGE}:${DEV_VERSION} || \:
	docker rmi ${STONEWORK_VPP_IMAGE}:${VPP_VERSION} || \:
	docker rmi ${TESTER_IMAGE}:${DEV_VERSION} || \:

# -------------------------------
#  VM image
# -------------------------------

vm-image: images # unmaintained
	@echo "=> building stonework VM image"
ifdef CNFS_SPEC
	VERSION=${VPP_VERSION} \
	./scripts/vm/create-vm-image.sh --cnfs-spec ${CNFS_SPEC}
else
	VERSION=${VPP_VERSION} \
	./scripts/vm/create-vm-image.sh
endif

# -------------------------------
#  Release
# -------------------------------

md-to-pdf:
	pandoc README.md -o README.pdf "-fmarkdown-implicit_figures -o" --from=markdown -V geometry:margin=.6in -V colorlinks --toc --highlight-style=espresso

generate-config-docs:
	echo "${STONEWORK_PROD_IMAGE}:${DEV_VERSION}" | bash -x ./scripts/gen-docs.sh

release:
ifneq ($(RELEASE_TAG_CHECKED),)
	RELEASE_TAG=$(RELEASE_VERSION_FULL) \
	STONEWORK_IMAGE="$(REPO)/$(STONEWORK_PROD_IMAGE):$(RELEASE_VERSION_FULL)" \
	./scripts/release.sh
else
	@echo "Release tag is empty or has incorrect format."
	@echo "Supplied release tag: ${RELEASE_TAG}"
	@echo 'Expected format: ${TAG_FORMAT} ; must not contain "-dirty"'
	@false
endif

# -------------------------------
#  Testing
# -------------------------------

test: unit-tests e2e-tests

unit-tests:
	go test ./...

e2e-tests:
	STONEWORK_IMAGE="${STONEWORK_PROD_IMAGE}:${DEV_VERSION}" ./scripts/e2e-test.sh

test-vpp-plugins: vpp-image vpp-test-image

test-vpp-plugins-prebuilt: # For running VPP tests repeatedly (saves time by skipping building process)
	docker run --privileged --name=vpp-test -d ${STONEWORK_VPP_TEST_IMAGE}:${VPP_VERSION}
	docker cp ./vpp/isisx/vpp$(shell echo ${VPP_VERSION} | tr -d ".")/isisx vpp-test:/opt/dev/vpp/src/plugins/isisx
	docker cp ./vpp/abx/vpp$(shell echo ${VPP_VERSION} | tr -d ".")/abx vpp-test:/opt/dev/vpp/src/plugins/abx
	docker exec -it vpp-test sh -c "cd /opt/dev/vpp;make test TEST=isisx;make test TEST=abx"
	docker rm -f vpp-test

# -------------------------------
#  Development
# -------------------------------

get-binapi-generator:
	@echo "# installing binary API generator"
	go install go.fd.io/govpp/cmd/binapi-generator

get-descriptor-adapter-generator:
	@echo "# installing descriptor adapter generator"
	go install go.ligato.io/vpp-agent/v3/plugins/kvscheduler/descriptor-adapter

generate-proto: protocgengo ## Generate Protobuf files

# FIXME: Currently generate-binapi-from-system-vpp generates binapi only from VPP .api.json files
# located at /usr/share/vpp/api/. Therefore the binapi will be generated only for a single VPP
# version that user has installed on their system. Modify this to be able to generate binapis for
# different VPP versions. For example by calling a script (with a VPP_VERSION value) that will
# generate the binapi from a docker container that has given version of VPP .api.json files with
# volume mounted to this repository.
# Consider using stonework-dev:<VPP_VERSION> image for this purpose.
#
# NOTE: Before running this make sure that the VPP .api.json files in /usr/share/vpp/api on your
# system belong to VPP version 23.06.x and not other version of VPP! Also do not forget that
# StoneWork VPP plugins (abx and isisx) .api.json files have to be copied into the
# /usr/share/vpp/api/plugins directory as well.
generate-binapi-from-system-vpp: get-binapi-generator
	@echo "=> generating binary API"
	@cd plugins/binapi/vpp2306 && VPP_VERSION=23.06 go generate .

generate-descriptor-adapters: get-descriptor-adapter-generator
	@echo "# generating descriptor adapters"
	go generate -x -run=descriptor-adapter ./plugins/...

replace-stonework:
	@cd ./scripts/dev && ./replace-stonework.sh
replace-mockcnf:
	@cd ./scripts/dev && ./replace-mockcnf.sh
