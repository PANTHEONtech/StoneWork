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

ARG VPP_IMAGE=vpp:22.10
ARG DEV_IMAGE=stonework-dev:22.10

FROM $VPP_IMAGE as vpp
FROM $DEV_IMAGE as dev
FROM ubuntu:20.04 as base

# utils just for testing - to be removed
RUN set -ex; \
    apt-get update && apt-get install -y \
		binutils \
		curl \
		iproute2 \
		iputils-ping \
		lshw \
		net-tools \
		netcat-openbsd; \
    rm -rf /var/lib/apt/lists/*

# Install vpp (except for ikev2 plugin which we are in a conflict with)
RUN mkdir -p /vpp
COPY --from=vpp \
		/vpp/vpp_*.deb \
		/vpp/libvppinfra_*.deb \
		/vpp/vpp-plugin-dpdk_*.deb \
		/vpp/vpp-plugin-core_*.deb \
	/vpp/

RUN set -ex; \
    cd /vpp; \
    apt-get update; \
    apt-get install -y ./*.deb; \
    rm /usr/lib/x86_64-linux-gnu/vpp_plugins/ikev2_plugin.so; \
    rm *.deb; \
    rm -rf /var/lib/apt/lists/*;

# install custom built vpp plugins
COPY --from=vpp \
    /usr/lib/x86_64-linux-gnu/vpp_plugins/abx_plugin.so \
    /usr/lib/x86_64-linux-gnu/vpp_plugins/isisx_plugin.so \
    /usr/lib/x86_64-linux-gnu/vpp_plugins/

# Install control-plane
COPY --from=dev \
		/usr/local/bin/stonework-init \
		/usr/local/bin/stonework \
		/usr/local/bin/agentctl \
	/usr/local/bin/

# Install config files
COPY --from=dev /etc/stonework /etc/stonework
COPY --from=dev /etc/vpp /etc/vpp
COPY ./docker/init_hook.sh /usr/bin/
ENV CONFIG_DIR /etc/stonework/

# Install API definitions
COPY --from=dev /api /api

# Install script for packet tracing on VPP
COPY ./docker/vpptrace.sh /usr/bin/vpptrace.sh
RUN chmod u+x /usr/bin/vpptrace.sh

# Final image
FROM scratch
COPY --from=base / /

ENV CONFIG_DIR /etc/stonework/
ENV CNF_MODE STONEWORK
CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
    mkdir -p /run/vpp /run/stonework/vpp && \
    exec stonework-init
