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

VERSION ?= $(shell git describe --always --tags --dirty)
COMMIT  ?= $(shell git rev-parse HEAD)
DATE    ?= $(shell git log -1 --format="%ct" | xargs -I{} date -d @{} +'%Y-%m-%dT%H:%M%:z')

CNINFRA := go.ligato.io/cn-infra/v2/agent
LDFLAGS = -X $(CNINFRA).BuildVersion=$(VERSION) -X $(CNINFRA).CommitHash=$(COMMIT) -X $(CNINFRA).BuildDate=$(DATE)

ifeq ($(VPP_VERSION),)
VPP_VERSION="21.06"
endif
REPO="ghcr.io/pantheontech/"
STONEWORK_VPP_IMAGE=${REPO}"vpp"
STONEWORK_VPP_TEST_IMAGE=${REPO}"vpp-test"
STONEWORK_DEV_IMAGE=${REPO}"stonework-dev"
STONEWORK_PROD_IMAGE=${REPO}"stonework"
TESTER_IMAGE=${REPO}"stonework-tester"
MOCK_CNF_IMAGE=${REPO}"stonework-mockcnf"
PROTO_ROOTGEN_IMAGE=${REPO}"proto-rootgen"

# Go.mod unrelated version locking
BINAPI_GENERATOR_COMMIT="4c1cccf48cd144414c7233f167087aff770ef67b" # no actual tag, newest tag is "0.3.5" and it is older commit and it is incompatible


build:
	@cd cmd/stonework && go build -v -ldflags "${LDFLAGS}"
	@cd cmd/stonework-init && go build -v -ldflags "${LDFLAGS}"
	@cd cmd/mockcnf && go build -v -ldflags "${LDFLAGS}"
	@cd cmd/proto-rootgen && go build -v -ldflags "${LDFLAGS}"

install:
	@cd cmd/stonework && go install -v -ldflags "${LDFLAGS}"
	@cd cmd/stonework-init && go install -v -ldflags "${LDFLAGS}"

install-mockcnf:
	@cd cmd/mockcnf && go install -v -ldflags "${LDFLAGS}"

install-proto-rootgen:
	@cd cmd/proto-rootgen && go install -v -ldflags "${LDFLAGS}"

# -------------------------------
#  Images
# -------------------------------

vpp-image:
	@echo "=> building VPP image -- version=${VPP_VERSION}"
	IMAGE_TAG=${STONEWORK_VPP_IMAGE} \
	VERSION=${VPP_VERSION} \
	./scripts/build.sh vpp

vpp-test-image:
	@echo "=> building VPP test image -- version=${VPP_VERSION}"
	IMAGE_TAG=${STONEWORK_VPP_TEST_IMAGE} \
	VERSION=${VPP_VERSION} \
	./scripts/build.sh vpp-test

dev-image:
	@echo "=> building development image, VPP version=${VPP_VERSION}"
	VPP_IMAGE="${STONEWORK_VPP_IMAGE}:${VPP_VERSION}" \
	IMAGE_TAG=${STONEWORK_DEV_IMAGE} \
	VERSION=${VPP_VERSION} \
	./scripts/build.sh dev

prod-image:
	@echo "=> building production image, VPP version=${VPP_VERSION}"
	VPP_IMAGE="${STONEWORK_VPP_IMAGE}:${VPP_VERSION}" \
	IMAGE_TAG=${STONEWORK_PROD_IMAGE} \
	DEV_IMAGE_TAG=${STONEWORK_DEV_IMAGE} \
	VERSION=${VPP_VERSION} \
	./scripts/build.sh prod

tester-image:
	@echo "=> building image with network tools for testing"
	IMAGE_TAG=${TESTER_IMAGE} \
	./scripts/build.sh tester

mockcnf-image:
	@echo "=> building mock CNF, VPP version=${VPP_VERSION}"
	VPP_IMAGE="${STONEWORK_VPP_IMAGE}:${VPP_VERSION}" \
	IMAGE_TAG=${MOCK_CNF_IMAGE} \
	VERSION=${VPP_VERSION} \
	./scripts/build.sh mockcnf

proto-rootgen-image:
	@echo "=> building image for building proto file with the config root message"
	IMAGE_TAG=${PROTO_ROOTGEN_IMAGE} \
	./scripts/build.sh proto-rootgen

images: vpp-image dev-image prod-image tester-image mockcnf-image
	# tag latest images
	docker tag ${STONEWORK_PROD_IMAGE}:${VPP_VERSION} ${STONEWORK_PROD_IMAGE}
	docker tag ${MOCK_CNF_IMAGE}:${VPP_VERSION} ${MOCK_CNF_IMAGE}

push-images:
	docker push ${STONEWORK_PROD_IMAGE}:${VPP_VERSION}
	docker push ${STONEWORK_PROD_IMAGE}
	docker push ${TESTER_IMAGE}

cleanup-images:
	docker rmi ${STONEWORK_DEV_IMAGE}:${VPP_VERSION} || \:
	docker rmi ${MOCK_CNF_IMAGE}:${VPP_VERSION} || \:
	docker rmi ${STONEWORK_VPP_IMAGE}:${VPP_VERSION} || \:
	docker rmi ${TESTER_IMAGE} || \:

# -------------------------------
#  VM image
# -------------------------------

vm-image: images
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
	echo "${STONEWORK_PROD_IMAGE}:${VPP_VERSION}" | bash -x ./scripts/gen-docs.sh

release: dev-image prod-image
	RELEASE_TAG=$(VPP_VERSION) \
	STONEWORK_IMAGE="$(STONEWORK_PROD_IMAGE):$(VPP_VERSION)" \
	./scripts/release.sh

# -------------------------------
#  Testing
# -------------------------------

test: unit-tests e2e-tests

unit-tests:
	go test ./...

e2e-tests:
	./scripts/e2e-test.sh

test-vpp-plugins: vpp-image vpp-test-image

test-vpp-plugins-prebuilt: # For running VPP tests repeatedly (saves time by skipping building process)
	docker run --privileged --name=vpp-test -d ${STONEWORK_VPP_TEST_IMAGE}:${VPP_VERSION}
	docker cp ./vpp/isisx/vpp$(shell echo ${VPP_VERSION} | tr -d ".")/isisx vpp-test:/opt/dev/vpp/src/plugins/isisx
	docker exec -it vpp-test sh -c "cd /opt/dev/vpp;make test TEST=isisx"
	docker rm -f vpp-test

# -------------------------------
#  Development
# -------------------------------

get-binapi-generator:
	@# temp directory is "go install" fix for <go1.16 (go.mod in root is changed but shouldn't)
	@# when upgraded to >=go1.16 use "go install" as is instead of "go get" + temp directory
	@echo "# installing binary API generator"
	$(eval TMP_DIR := $(shell mktemp -d))
	cd $(TMP_DIR);GO111MODULE=on go get git.fd.io/govpp.git/cmd/binapi-generator@$(BINAPI_GENERATOR_COMMIT)
	rm -rf $(TMP_DIR)

get-descriptor-adapter-generator:
	@echo "# installing descriptor adapter generator"
	cd submodule/vpp-agent;go install ./plugins/kvscheduler/descriptor-adapter

dep-install:
	@go mod download

generate-proto:
	@echo "=> generating proto files"
	./scripts/gen-proto.sh

generate-binapi: get-binapi-generator
    # generated from vpp json api files copied into Stonework repository (plugins/binapi/vppXXXX/api)
    # from VPP (/usr/share/vpp/api/(core|plugins))
	@echo "=> generating binary API"
	@cd plugins/binapi/vpp2009 && VPP_VERSION=20.09 go generate .
	@cd plugins/binapi/vpp2101 && VPP_VERSION=21.01 go generate .
	@cd plugins/binapi/vpp2106 && VPP_VERSION=21.06 go generate .

generate-descriptor-adapters: get-descriptor-adapter-generator
	@echo "# generating descriptor adapters"
	go generate -x -run=descriptor-adapter ./plugins/...

replace-stonework:
	@cd ./scripts/dev && ./replace-stonework.sh
replace-mockcnf:
	@cd ./scripts/dev && ./replace-mockcnf.sh
