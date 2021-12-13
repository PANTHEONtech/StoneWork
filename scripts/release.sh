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

set -o errexit
set -o nounset
set -o pipefail


RELEASE_DIR="./stonework-release-${RELEASE_TAG}"
STONEWORK_IMAGE="${STONEWORK_IMAGE}"

if [ -d "${RELEASE_DIR}" ]; then
    echo "ERROR: Release directory already exists (${RELEASE_DIR}). Please remove it first."
    exit 1
fi

echo " => Creating release directory ${RELEASE_DIR}."
mkdir "${RELEASE_DIR}"

echo " => Archiving ${STONEWORK_IMAGE} docker image."
docker save --output "${RELEASE_DIR}/stonework.image" "${STONEWORK_IMAGE}"

echo " => Copying documentation"
cp README.md "${RELEASE_DIR}/"
cp EULA "${RELEASE_DIR}/"
cp LICENSE "${RELEASE_DIR}/"
cp THIRD_PARTY_LICENSES "${RELEASE_DIR}/"
mkdir -p "${RELEASE_DIR}/docs/"
cp -r docs/config "${RELEASE_DIR}/docs/"

echo " => Copying examples"
cp -r ./examples "${RELEASE_DIR}/"
rm -r "${RELEASE_DIR}/examples/testing"

echo " => Archiving release to ${RELEASE_DIR}.tar.gz"
tar -czvf ${RELEASE_DIR}.tar.gz "${RELEASE_DIR}"

echo " => Cleaning"
rm -r ${RELEASE_DIR}
