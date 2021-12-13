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

set -euo pipefail

OUTDIR=${1:-./docs/config}
CNFNAME=${2:-StoneWork}

# create output directory
mkdir -p ${OUTDIR}

# Create a temporary directory and store its name in a variable
TMPDIR=$(mktemp -d)
echo "Temporary directory is ${TMPDIR}"

# Bail out if the temp directory wasn't created successfully.
if [ ! -e $TMPDIR ]; then
    >&2 echo "Failed to create temp directory"
    exit 1
fi

# Make sure it gets removed even if the script exits abnormally.
trap "exit 1"           HUP INT PIPE QUIT TERM
trap 'rm -rf "$TMPDIR"' EXIT

while IFS= read -r line; do
  image="${line}"
  echo "CNF ${image} will be included in the generated docs."
  mkdir ${TMPDIR}/tmp
  # copy CNF API definitions
  docker create -ti --name dummy-cnf-image ${image} bash
  trap 'docker rm -f dummy-cnf-image 2>/dev/null || true' EXIT
  docker cp dummy-cnf-image:/api ${TMPDIR}/tmp
  rsync -v --recursive --exclude '*.yaml' ${TMPDIR}/tmp/api/ /${TMPDIR}
  cat ${TMPDIR}/tmp/api/models.spec.yaml >> ${TMPDIR}/models.spec.yaml
  # cleanup
  docker rm -f dummy-cnf-image
  rm -rf ${TMPDIR}/tmp
done

# generate root proto message
docker run --rm -v ${TMPDIR}:/api ghcr.io/pantheontech/proto-rootgen proto-rootgen --cnf-name ${CNFNAME}

# generate json schema
docker run --rm -v ${TMPDIR}:/api ghcr.io/pantheontech/proto-rootgen \
  protoc \
    --jsonschema_out="json_fieldnames:/api" \
    --proto_path=/api /api/${CNFNAME,,}-root.proto
cp ${TMPDIR}/Root.jsonschema ${OUTDIR}/${CNFNAME^^}-CONFIG.jsonschema

# generate docs (as markdown & pdf)
docker run --rm -v ${TMPDIR}:/api ghcr.io/pantheontech/proto-rootgen \
  bash -x -c "\
    sed -i 's/{{.CnfName}}/${CNFNAME}/g' /gendoc/markdown.tmpl &&
    protos=\$(find /api -name \"*.proto\" | grep -vF \"/api/${CNFNAME,,}-root.proto\") &&
    protoc \
      --doc_out=/api \
      --doc_opt=/gendoc/markdown.tmpl,CONFIG.md \
      --proto_path=/api /api/${CNFNAME,,}-root.proto \${protos}
    pandoc /api/CONFIG.md -o /api/CONFIG.pdf \
      \"-fmarkdown-implicit_figures -o\" --from=markdown -V geometry:margin=.6in \
      -V colorlinks -H /gendoc/pandoc-preamble.tex --highlight-style=espresso
  "
cp ${TMPDIR}/CONFIG.md ${OUTDIR}/${CNFNAME^^}-CONFIG.md
cp ${TMPDIR}/CONFIG.pdf ${OUTDIR}/${CNFNAME^^}-CONFIG.pdf

rm -rf ${TMPDIR}
exit 0
