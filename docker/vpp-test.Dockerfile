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

ARG VPP_IMAGE=vpp:23.06

FROM ${VPP_IMAGE}

# `netbase` has to be installed because it provides /etc/protocols file.
# The file is used by Scapy which the VPP Test Framework depends on.
RUN set -ex; \
    apt-get update && \
    apt-get install -y --no-install-recommends \
		libssl-dev \
		netbase \
		pkg-config \
		python3-venv \
    && rm -rf /var/lib/apt/lists/*

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

RUN set -ex; \
    cp -r /opt/dev/abx/abx /opt/dev/vpp/src/plugins/abx && \
    cp -r /opt/dev/isisx/isisx /opt/dev/vpp/src/plugins/isisx && \
    cp -r /opt/dev/isisx/test /opt/dev/vpp/ && \
    cp -r /opt/dev/abx/test /opt/dev/vpp/

WORKDIR /opt/dev/vpp

ARG CACHEBUST=1
RUN make test TEST=isisx
RUN make test TEST=abx
