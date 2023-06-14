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

ARG VPP_VERSION=23.06
ARG VPP_IMAGE=ligato/vpp-base:$VPP_VERSION

FROM ${VPP_IMAGE}

ARG DEBIAN_FRONTEND=noninteractive
RUN set -ex; \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    	git \
    	ca-certificates \
    	build-essential \
    	sudo \
    	cmake \
    	nasm \
    	ninja-build \
    	python3-ply \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/dev

RUN git clone https://gerrit.fd.io/r/vpp
RUN cd vpp && \
    COMMIT=$(cat /vpp/version | sed -n 's/.*[~-]g\([a-z0-9]*\).*/\1/p' | \
        (grep . || sh -c 'echo "Cant detect commit of VPP from VPP image" 1>&2;exit 1')) && \
    git checkout $COMMIT && git show  -s

#----------------------
# build & install external plugins (ABX, ISISX)
ARG VPP_VERSION
COPY vpp/abx /tmp/abx
COPY vpp/isisx /tmp/isisx

# Plugins
RUN set -ex; \
    VPPVER=$(echo $VPP_VERSION | tr -d ".") && \
    cp -r /tmp/abx/vpp${VPPVER} /opt/dev/abx && \
    cp -r /tmp/isisx/vpp${VPPVER} /opt/dev/isisx    

# Plugin ABX
RUN set -ex; \
	cd abx; \
    ./build.sh /opt/dev/vpp/
RUN set -ex; \
	cp /opt/dev/abx/build/lib/vpp_plugins/abx_plugin.so \
       /usr/lib/x86_64-linux-gnu/vpp_plugins/; \
	cp /opt/dev/abx/build/abx/abx.api.json \
       /usr/share/vpp/api/core/

# Plugin plugin
RUN set -ex; \
	cd ./isisx; \
    ./build.sh /opt/dev/vpp/
RUN set -ex; \
	cp /opt/dev/isisx/build/lib/vpp_plugins/isisx_plugin.so \
       /usr/lib/x86_64-linux-gnu/vpp_plugins/; \
	cp /opt/dev/isisx/build/isisx/isisx.api.json \
       /usr/share/vpp/api/core/

# there is a bug in VPP 21.06 that api files are not built on standard location
# for external plugins, to reproduce it is enough to try to build sample-plugin
RUN set -ex; \
    if [ "$VPP_VERSION" = "23.06" ] || [ "$VPP_VERSION" = "22.10" ] || [ "$VPP_VERSION" = "22.02" ]; \
    then \
      cp abx/build/CMakeFiles/vpp-api/vapi/* /usr/include/vapi/; \
    elif [ "$VPP_VERSION" = "21.06" ]; \
	then \
	  cp /vpp-api/vapi/* /usr/include/vapi/; \
	else \
      cp /opt/dev/abx/build/vpp-api/vapi/* /usr/include/vapi/; \
    fi

COPY docker/vpp-startup.conf /etc/vpp/startup.conf

CMD /bin/bash -c "mkdir -p /run/stonework/vpp; \
	exec /usr/bin/vpp -c /etc/vpp/startup.conf"
