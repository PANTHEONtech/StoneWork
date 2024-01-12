FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
		build-essential \
		ca-certificates \
		git \
		make \
		nano \
		wget \
 	&& rm -rf /var/lib/apt/lists/*

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
RUN mkdir -p "$GOPATH/bin" && chmod -R 777 "$GOPATH"


# Install agentctl
RUN mkdir -p "/src/ligato" && \
    cd "/src/ligato" && \
    git clone https://github.com/ligato/vpp-agent.git
WORKDIR /src/ligato/vpp-agent
ARG COMMIT
RUN git checkout $COMMIT

ARG VERSION
ARG BRANCH
ARG BUILD_DATE
RUN make agentctl

CMD exec agentctl
