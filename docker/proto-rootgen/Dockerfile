
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

FROM ubuntu:18.04 as base

RUN apt-get update && apt-get install -y \
    git \
    gcc \
    make \
    wget \
    curl \
    zip \
    rsync \
    pandoc \
    texlive &&\
    rm -rf /var/lib/apt/lists/*

# Install Go
ENV GOLANG_VERSION 1.14.1
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

# Install protoc and some plugins to generate docs
ARG PROTOC_VERSION=3.12.4
WORKDIR /tmp
RUN curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-$(uname -m).zip -o protoc.zip
RUN unzip protoc.zip && mv bin/protoc /usr/local/bin/protoc && mv include/google /usr/local/include/google
RUN go get -v github.com/chrusty/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema
RUN go get -v github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
COPY docker/proto-rootgen/markdown.tmpl /gendoc/markdown.tmpl
COPY docker/proto-rootgen/pandoc-preamble.tex /gendoc/pandoc-preamble.tex

# Build proto-rootgen
ENV GO111MODULE=on
RUN mkdir -p $GOPATH/src/pantheon.tech/StoneWork
WORKDIR $GOPATH/src/pantheon.tech/StoneWork
COPY . ./
RUN make dep-install
RUN make install-proto-rootgen
