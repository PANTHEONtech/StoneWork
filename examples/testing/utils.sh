#!/usr/bin/env bash

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

# parameters:
# - name of the container to check for config items
# - desired number of occurences
# - timeout, in seconds
function waitForAgentConfig {
	SLEPT=0
	while ! [ $(docker compose exec -T $1 agentctl values 2>/dev/null | grep -c -E "CONFIGURED|obtained") -ge $2 ]
	do
		sleep 1
		SLEPT=$((SLEPT+1))
		if [ $SLEPT -ge $3 ]; then
			break
		fi
	done
	sleep 1 # extra wait in case not all interfaces are up yet
}

$*
