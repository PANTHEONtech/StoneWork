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

terminate_process() {
  if [ -n "NOTERM_INIT_PROC" ] && { echo "WARNING: Termination of init process is disabled via NOTERM_INIT_PROC env var"; return }
  PID=$(pidof $1)
  if [[ ${PID} != "" ]]; then
    kill ${PID}
    echo "process $1 terminated"
  fi
}

if [[ "${SUPERVISOR_PROCESS_NAME}" == "agent" && "${SUPERVISOR_PROCESS_STATE}" == "terminated" ]]; then
  terminate_process stonework-init
fi

if [[ "${SUPERVISOR_PROCESS_NAME}" == "vpp" && "${SUPERVISOR_PROCESS_STATE}" == "terminated" ]]; then
  terminate_process stonework-init
fi
