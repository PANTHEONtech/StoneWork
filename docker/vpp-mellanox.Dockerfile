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

################ !!! EXPERIMENTAL IMAGE !!! ###############################3

# This Dockerfile builds VPP image with DPDK Mellanox PMDs
# To use it, rename it to vpp.Dockerfile, as:
# $ mv vpp-mellanox.Dockerfile vpp.Dockerfile
# and then build StoneWork as usual:
# $ make images
# The StoneWork build system will use it instead of standard VPP Dockerfile.

# NOTE: Since Mellanox PMDs uses Netlink, you will additionally need to set
# host network_mode (besides making it privileged and mounting /dev volume as
# ussual with DPDK), to use StoneWork with this VPP under the hood.

ARG VPP_VERSION=23.06
ARG VPP_IMAGE=ligato/vpp-base:$VPP_VERSION

FROM ${VPP_IMAGE} AS base

RUN set -ex; \
    apt-get update && \
	apt-get install -y --no-install-recommends \
		build-essential \
		cmake \
		git \
		ninja-build \
		sudo \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/dev

RUN git clone https://gerrit.fd.io/r/vpp
RUN cd vpp && \
    COMMIT=$(cat /vpp/version | sed -n 's/.*[~-]g\([a-z0-9]*\).*/\1/p' | \
        (grep . || sh -c 'echo "Cant detect commit of VPP from VPP image" 1>&2;exit 1')) && \
    git checkout $COMMIT && git show  -s

COPY docker/enable-mlx-pmds.patch /opt/dev
# NOTE: this patch may change for newer VPP versions or once
# https://gerrit.fd.io/r/c/vpp/+/31876
# is merged into stable/2101, then the update is needed !
RUN cd vpp && git apply ../enable-mlx-pmds.patch

RUN cd vpp && yes | make install-dep install-ext-deps && make pkg-deb

#-----------------
# build ABX plugin
ARG VPP_VERSION=23.06
COPY vpp/abx /tmp/abx
RUN VPPVER=$(echo $VPP_VERSION | tr -d ".") && \
    cp -r /tmp/abx/vpp${VPPVER} /opt/dev/abx

RUN cd abx && ./build.sh /opt/dev/vpp/

FROM ubuntu:22.04

RUN set -ex; \
    apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates \
		curl \
		gnupg \
		iproute2 \
		iputils-ping \
		python3 \
		python3-cffi \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /vpp
COPY --from=base /opt/dev/vpp/build-root/*.deb ./

RUN set -eux; \
    apt-get update && apt-get install -y -V ./*.deb; \
    dpkg-query -f '${Version}\n' -W vpp > /vpp/version; \
    rm -rf /var/lib/apt/lists/*; \
    rm ./*.deb

RUN mkdir -p /var/log/vpp

#-------------------
# install ABX plugin
COPY --from=base \
    /opt/dev/abx/build/lib/vpp_plugins/abx_plugin.so \
    /usr/lib/x86_64-linux-gnu/vpp_plugins/

COPY --from=base \
    /opt/dev/abx/build/abx/abx.api.json \
    /usr/share/vpp/api/core/

COPY --from=base \
    /opt/dev/abx/build/vpp-api/vapi/* \
    /usr/include/vapi/

CMD ["/usr/bin/vpp", "-c", "/etc/vpp/startup.conf"]
