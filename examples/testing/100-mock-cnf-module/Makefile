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

test:
	@make exec-test || (make dump-logs-ignore-errors && make stop-example-ignore-errors && make environment-cleanup && exit 1)

dump-logs-ignore-errors:
	@make dump-logs || true

exec-test: stop-example-ignore-errors environment-cleanup start-example test-stonework stop-example

start-example:
	@docker-compose up -d
	@../utils.sh waitForAgentConfig stonework 52 20

stop-example:
	@docker-compose down --remove-orphans -v

stop-example-ignore-errors:
	@docker-compose down --remove-orphans -v || true

# fix for "existing container with the same name" (docker-compose down of this test can't match and clean them
# if it is for example other image. This points to some environment cleanup bug in some other test)
environment-cleanup:
	@cat docker-compose.yaml | grep container_name: | awk '{print $$2}' | xargs docker rm -f 1> /dev/null 2> /dev/null

test-stonework:
	@./test-stonework.sh

LOGFILE=./example.log

dump-logs:
	@sh -c "(\
		echo ===============================================================================;\
		echo StoneWork logs;\
		echo ===============================================================================;\
		docker logs stonework || true;\
		echo ===============================================================================;\
		echo StoneWork running configuration;\
		echo ===============================================================================;\
		docker exec stonework agentctl config retrieve || true;\
		echo ===============================================================================;\
		echo Mock CNF 1 logs;\
		echo ===============================================================================;\
		docker logs mockcnf1 || true;\
		echo ===============================================================================;\
		echo Mock CNF 2 logs;\
		echo ===============================================================================;\
		docker logs mockcnf2 || true;\
	) >${LOGFILE} 2>&1"
