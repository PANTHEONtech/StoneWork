#!/bin/bash

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

# exit when any command fails
set -e -o pipefail

BUILD=build # must be executed from workspace root
FILENAME=stonework.qcow2
VM_IMAGE=$BUILD/$FILENAME
BASE_OS=ubuntu-20.04
# NOTE: to list possibilities for BASE_OS use: virt-builder --list

rm -rf $BUILD
mkdir -p $BUILD

# parse arguments
ARGUMENT_LIST=(
    "cnfs-spec"
)

opts=$(getopt \
    --longoptions "$(printf "%s:," "${ARGUMENT_LIST[@]}")" \
    --name "$(basename "$0")" \
    --options "" \
    -- "$@"
)

eval set --$opts

while [[ $# -gt 0 ]]; do
    case "$1" in
        --cnfs-spec)
            CNFS_SPEC=$2
            shift 2
            ;;

        *)
            break
            ;;
    esac
done

echo "*** Export stonework docker image."
docker save ghcr.io/pantheontech/stonework \
      > $BUILD/stonework-docker-img.tar

echo "*** Build VM image."
sudo virt-builder $BASE_OS \
--format qcow2 \
--commands-from-file scripts/vm/virt-builder-commands \
--output $VM_IMAGE

echo "*** Install firstboot script."
sudo virt-customize -a $VM_IMAGE --firstboot scripts/vm/firstboot.sh

cp scripts/vm/docker-compose.yaml $BUILD

if [ ! -z ${CNFS_SPEC} ]; then
    echo "*** Adding CNFs."
    sudo scripts/vm/add-cnfs.py \
    --cnfs-spec $CNFS_SPEC \
    --docker-compose $BUILD/docker-compose.yaml \
    --vm-image $VM_IMAGE
fi

sudo virt-customize -a $VM_IMAGE --copy-in build/docker-compose.yaml:/root/

echo "*** Shrink VM image."
qemu-img convert -c -O qcow2 $BUILD/stonework.qcow2 $BUILD/shrunk.qcow2
sudo rm $VM_IMAGE
mv $BUILD/shrunk.qcow2 $VM_IMAGE

echo "*** Create GNS3 appliance file."
cp scripts/vm/stonework-vm.gns3a $BUILD
sed -i "s/<FILENAME>/$FILENAME/" $BUILD/stonework-vm.gns3a
sed -i "s/<VERSION>/$VERSION/" $BUILD/stonework-vm.gns3a

IMAGE_MD5SUM=`md5sum $VM_IMAGE | awk -F" " '{ print $1 }'`
sed -i "s/<IMAGE_MD5SUM>/$IMAGE_MD5SUM/" $BUILD/stonework-vm.gns3a

IMAGE_SIZE=`ls -l $VM_IMAGE | awk -F" " '{ print $5 }'`
sed -i "s/<IMAGE_SIZE>/$IMAGE_SIZE/" $BUILD/stonework-vm.gns3a

echo "*** Remove intermediate results."
rm -f $BUILD/*-docker-img.tar \
      $BUILD/docker-compose.yaml
