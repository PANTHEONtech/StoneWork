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
ARG VPPAGENT_IMAGE=ligato/vpp-agent:v3.4.0

FROM $VPP_IMAGE as vpp
FROM $VPPAGENT_IMAGE as vppagent
FROM ubuntu:20.04 as base

RUN apt-get update && apt-get install -y \
		git \
		gcc \
		make \
		iptables \
		net-tools \
		# for debugging
		binutils \
		curl \
		wget \
		tcpdump \
		iproute2 \
		iputils-ping \
		# stats client
		python3 \
		python3-cffi && \
    rm -rf /var/lib/apt/lists/*

# Install Go
ENV GOLANG_VERSION 1.20.7
RUN set -eux; \
	dpkgArch="$(dpkg --print-architecture)"; \
		case "${dpkgArch##*-}" in \
			amd64) goRelArch='linux-amd64'; ;; \
			armhf) goRelArch='linux-armv6l'; ;; \
			arm64) goRelArch='linux-arm64'; ;; \
	esac; \
 	wget -nv -O go.tgz "https://golang.org/dl/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; \
 	tar -C /usr/local -xzf go.tgz; \
 	rm go.tgz;

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Install vpp (except for ikev2 plugin which we are in a conflict with)
RUN mkdir -p /vpp
COPY --from=vpp /vpp/vpp_*.deb \
    /vpp/libvppinfra_*.deb \
    /vpp/vpp-plugin-dpdk_*.deb \
    /vpp/vpp-plugin-core_*.deb \
    /vpp/

RUN cd /vpp/ \
    && apt-get update \
    && apt-get install -y ./*.deb \
    && rm /usr/lib/x86_64-linux-gnu/vpp_plugins/ikev2_plugin.so \
    && rm *.deb \
    && rm -rf /var/lib/apt/lists/*

# install custom built vpp plugins
COPY --from=vpp \
    /usr/lib/x86_64-linux-gnu/vpp_plugins/abx_plugin.so \
    /usr/lib/x86_64-linux-gnu/vpp_plugins/isisx_plugin.so \
    /usr/lib/x86_64-linux-gnu/vpp_plugins/

# Install agentctl
COPY --from=vppagent /bin/agentctl /usr/local/bin/agentctl

# Build agent
WORKDIR /src/
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN make install # stonework-init
RUN make install-mockcnf
RUN cp $GOPATH/bin/mockcnf /usr/local/bin/mockcnf
RUN cp $GOPATH/bin/stonework-init /usr/local/bin/stonework-init

# Install config files
RUN mkdir -p /etc/mockcnf /etc/mockcnf-novpp /etc/vpp
# -> STANDALONE
COPY ./docker/vpp-startup.conf /etc/vpp/vpp.conf
COPY ./docker/etcd.conf /etc/mockcnf/etcd.conf
COPY ./docker/grpc.conf /etc/mockcnf/grpc.conf
COPY ./docker/mockcnf/initfileregistry.conf /etc/mockcnf/initfileregistry.conf
COPY ./docker/mockcnf/supervisor.conf /etc/mockcnf/supervisor.conf
COPY ./docker/mockcnf/init_hook.sh /usr/bin/init_hook.sh
# -> SW MODULE
COPY ./docker/etcd.conf /etc/mockcnf-novpp/etcd.conf
COPY ./docker/grpc.conf /etc/mockcnf-novpp/grpc.conf
COPY ./docker/mockcnf/initfileregistry-novpp.conf /etc/mockcnf-novpp/initfileregistry.conf
COPY ./docker/mockcnf/supervisor-novpp.conf /etc/mockcnf-novpp/supervisor.conf
COPY ./docker/mockcnf/init_hook-novpp.sh /usr/bin/init_hook-novpp.sh

# Copy proto api (mockcnf models)
RUN mkdir -p /api/mockcnf
RUN cp -r ./proto/mockcnf /api/mockcnf

# Install script for packet tracing on VPP
COPY ./docker/vpptrace.sh /usr/bin/vpptrace.sh
RUN chmod u+x /usr/bin/vpptrace.sh

CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
    mkdir -p /run/vpp /run/stonework/vpp && \
    if [ "$CNF_MODE" = "STONEWORK_MODULE" ]; then CONFIG_DIR="/etc/mockcnf-novpp/"; else CONFIG_DIR="/etc/mockcnf/"; fi && \
    export CONFIG_DIR && \
    exec stonework-init
