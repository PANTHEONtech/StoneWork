# syntax = docker/dockerfile:1.2

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
		rsync \
		# for debugging
		binutils \
		curl \
		wget \
		tcpdump \
		iproute2 \
		iputils-ping \
		# stats client
		python3 \
		python3-cffi \
	&& rm -rf /var/lib/apt/lists/*

# Install Go
ENV GOLANG_VERSION 1.18.3
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

ENV CGO_ENABLED=0

# Build agent
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x
COPY . ./
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    set -eux; \
    make install; \
	mv $GOPATH/bin/stonework /usr/local/bin/; \
	mv $GOPATH/bin/swctl /usr/local/bin/; \
	mv $GOPATH/bin/stonework-init /usr/local/bin/

# Build api directory
RUN mkdir /api
RUN VPPAGENT_DIR=$(go list -f "{{ .Dir }}" "go.ligato.io/vpp-agent/v3") && \
    rsync -v --recursive --chmod=D2775,F444 --exclude '*.go' "${VPPAGENT_DIR}/proto/" /api/
RUN rsync -v --recursive --chmod=D2775,F444 \
    --exclude '*.go' --exclude 'puntmgr*' --exclude 'cnfreg*' --exclude 'mockcnf*' proto/ /api/
RUN /usr/local/bin/stonework-init --print-spec > /api/models.spec.yaml

# Install agentctl
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    set -eux; \
    go install go.ligato.io/vpp-agent/v3/cmd/agentctl; \
	mv $GOPATH/bin/agentctl /usr/local/bin/agentctl

# Install config files
RUN mkdir -p /etc/stonework /etc/vpp

COPY ./docker/vpp-startup.conf /etc/vpp/vpp.conf
COPY ./docker/etcd.conf /etc/stonework/etcd.conf
COPY ./docker/grpc.conf /etc/stonework/grpc.conf
COPY ./docker/aggregator.conf /etc/stonework/aggregator.conf
COPY ./docker/initfileregistry.conf /etc/stonework/initfileregistry.conf
COPY ./docker/supervisor.conf /etc/stonework/supervisor.conf
COPY ./docker/init_hook.sh /usr/bin/

ENV CONFIG_DIR /etc/stonework/
ENV CNF_MODE STONEWORK

# handle differences in vpp.conf which are between supported VPP versions
ARG VPP_VERSION
COPY ./docker/legacy-nat.conf /tmp/legacy-nat.conf
RUN bash -c "if [[ \"$VPP_VERSION\" < "21.01" ]]; then cat /tmp/legacy-nat.conf >> /etc/vpp/vpp.conf; fi"
RUN rm /tmp/legacy-nat.conf

# Install script for packet tracing on VPP
COPY ./docker/vpptrace.sh /usr/bin/vpptrace.sh
RUN chmod u+x /usr/bin/vpptrace.sh

COPY ./plugins/binapi/vpp2306/api/abx.api.json /usr/share/vpp/api/plugins/
COPY ./plugins/binapi/vpp2306/api/isisx.api.json /usr/share/vpp/api/plugins/

CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
    mkdir -p /run/vpp /run/stonework/vpp && \
    exec stonework-init
