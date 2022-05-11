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

ARG VPP_IMAGE=vpp:21.06

FROM ${VPP_IMAGE}

RUN apt-get update && \
    apt-get install -y python3-venv libssl-dev && \
    apt-get install -y pkg-config

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

RUN cp -r /opt/dev/abx/abx /opt/dev/vpp/src/plugins/abx && \
    cp -r /opt/dev/isisx/isisx /opt/dev/vpp/src/plugins/isisx && \
    cp -r /opt/dev/isisx/test /opt/dev/vpp/test

WORKDIR /opt/dev/vpp

ARG CACHEBUST=1
RUN make test TEST=isisx
