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

case $1 in
dev)
  docker build --build-arg VPP_IMAGE="${VPP_IMAGE}" \
    --build-arg VPP_VERSION="${VPP_VERSION}" \
    -t ${IMAGE_TAG} \
    -f ./docker/dev.Dockerfile .
  ;;
prod)
  docker build --build-arg VPP_IMAGE="${VPP_IMAGE}" \
    --build-arg DEV_IMAGE="${DEV_IMAGE}" \
    -t ${IMAGE_TAG} \
    -f ./docker/prod.Dockerfile .
  ;;
vpp)
  docker build --build-arg VPP_VERSION="${VPP_VERSION}" \
    -t ${IMAGE_TAG} \
    -f ./docker/vpp.Dockerfile .
  ;;
vpp-test)
  docker build --progress=plain --build-arg VPP_IMAGE="${VPP_IMAGE}" \
    -t ${IMAGE_TAG} \
    -f ./docker/vpp-test.Dockerfile .
  ;;
tester)
  docker build \
    -t ${IMAGE_TAG} \
    -f ./docker/tester.Dockerfile .
  ;;
mockcnf)
  docker build --build-arg VPP_IMAGE="${VPP_IMAGE}" \
    -t ${IMAGE_TAG} \
    -f ./docker/mockcnf/Dockerfile .
  ;;
proto-rootgen)
  docker build \
    -t ${IMAGE_TAG} \
    -f ./docker/proto-rootgen/Dockerfile .
esac
