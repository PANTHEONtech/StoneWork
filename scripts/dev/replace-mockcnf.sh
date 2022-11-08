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

# fail in case of error
set -e

# copy files to replace in the image
cp ../../cmd/mockcnf/mockcnf .

# rebuild the production image with replaced agent binary
docker build -f mockcnf.Dockerfile -t stonework-mockcnf:22.10 --no-cache --force-rm=true .

rm -f ./mockcnf
