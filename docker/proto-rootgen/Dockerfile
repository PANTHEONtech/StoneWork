
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

FROM ubuntu:20.04 as base

ARG DEBIAN_FRONTEND=noninteractive
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

# Install protoc and some plugins to generate docs
ARG PROTOC_VERSION=3.12.4
## protoc-gen-doc of version >1.5.0 use golang "embed" feature -> for that we must use newer golang version
ARG PROTOC_GEN_DOC_VERSION=1.5.0
## has version tagging but without "v" prefix -> need to use prefix of commit hash (f5fcc60 = tag 1.3.5,  dc9f108=1.0.2)
ARG PROTOC_GEN_JSONSCHEMA_PSEUDO_VERSION=dc9f108
WORKDIR /tmp
RUN curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-$(uname -m).zip -o protoc.zip
RUN unzip protoc.zip && mv bin/protoc /usr/local/bin/protoc && mv include/google /usr/local/include/google
RUN go install -v github.com/chrusty/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema@${PROTOC_GEN_JSONSCHEMA_PSEUDO_VERSION}
RUN go install -v github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v${PROTOC_GEN_DOC_VERSION}
COPY docker/proto-rootgen/markdown.tmpl /gendoc/markdown.tmpl
COPY docker/proto-rootgen/pandoc-preamble.tex /gendoc/pandoc-preamble.tex

# Build proto-rootgen
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN make install-proto-rootgen
