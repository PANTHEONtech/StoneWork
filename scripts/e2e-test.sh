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

RV=0

echo "==============================================================================="
echo "StoneWork basic - StoneWork as cross-connect"
echo "==============================================================================="
pushd examples/testing/010-xconnect
make test || RV=$?
popd

echo "==============================================================================="
echo "StoneWork basic - StoneWork as switch"
echo "==============================================================================="
pushd examples/testing/020-switch
make test || RV=$?
popd

echo "==============================================================================="
echo "StoneWork basic - StoneWork as router"
echo "==============================================================================="
pushd examples/testing/030-router
make test || RV=$?
popd

echo "==============================================================================="
echo "StoneWork basic - StoneWork as router (IPv6)"
echo "==============================================================================="
pushd examples/testing/040-router6
make test || RV=$?
popd

echo "==============================================================================="
echo "CNF as StoneWork module"
echo "==============================================================================="
pushd examples/testing/100-mock-cnf-module
make test || RV=$?
popd

echo "==============================================================================="
echo "Standalone CNF"
echo "==============================================================================="
pushd examples/testing/110-mock-cnf-standalone
make test || RV=$?
popd

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

if [ $RV -ne 0 ]; then
    echo "===============================================================================";
    echo -e "${RED}[FAIL] Some tests failed${NC}";
    echo "===============================================================================";
else
    echo "===============================================================================";
    echo -e "${GREEN}[OK] All tests passed${NC}";
    echo "===============================================================================";
fi

exit $RV
