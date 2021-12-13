#!/usr/bin/env python3

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

import argparse
import yaml
import re
import docker
import ntpath
import os

ap = argparse.ArgumentParser()

# This script adds CNFs to the StoneWork VM and updates the docker-compose.yaml
# It takes an input file as --cnfs-spec <yaml-file> in following format, where
# volumes part is optional and may contain multiple entries

# - image: <cnf-repo>/<cnf-name>:<tag>
#   license: <path-to-license-file>
#   volumes:
#    - <host-path>:<cnf-path>

# NOTE: The path to the license file and host paths of volumes are paths
# relative to the StoneWork workspace root since make is executed from there.
# Absolute paths can be used as well.

ap.add_argument('-c', '--cnfs-spec', required=True,
                help='CNFs specification YAML file')
ap.add_argument('-d', '--docker-compose', required=True,
                help='Docker compose file to extend')
ap.add_argument('-i', '--vm-image', required=True,
                help='VM image to customize')

args = vars(ap.parse_args())


def save_docker_img(name, path):
    """
    Saves docker image the same way as 'docker save ..' command.
    """
    client = docker.from_env()
    image = client.images.get(name)
    save_gen = image.save(named=True)

    with open(path, 'wb') as img_f:
        for chunk in save_gen:
            img_f.write(chunk)


def path_leaf(path):
    head, tail = ntpath.split(path)
    return tail or ntpath.basename(head)


with open(args['cnfs_spec'], 'r') as cnfs_spec, \
     open(args['docker_compose'], 'a') as dc:
    y = yaml.safe_load(cnfs_spec)
    if y is None:
        raise Exception('Empty CNF specification file')
    for cnf in y:
        print(cnf)
        m = re.match(r'(.*)/(.*):(.*)', cnf['image'])
        cnf_repo = m.group(1)
        cnf_name = m.group(2)
        tag = m.group(3)
        img = 'build/{}-docker-img.tar'.format(cnf_name)
        save_docker_img(cnf['image'], img)
        os.system('virt-customize -a {0} --copy-in {1}:/root/'
                  .format(args['vm_image'], img))
        os.system('virt-customize -a {0} --copy-in {1}:/root/'
                  .format(args['vm_image'], cnf['license']))
        dc.write('\n'
                 '  {1}:\n'
                 '    container_name: {1}\n'
                 '    image: "{0}/{1}:{2}"\n'
                 '    depends_on:\n'
                 '      - stonework\n'
                 '    privileged: true\n'
                 '    env_file:\n'
                 '      - {3}\n'
                 '    environment:\n'
                 '      CNF_MODE: "STONEWORK_MODULE"\n'
                 '      INITIAL_LOGLVL: "debug"\n'
                 '      MICROSERVICE_LABEL: "{1}"\n'
                 '      ETCD_CONFIG: ""\n'
                 '    volumes:\n'
                 '      - runtime_data:/run/stonework\n'
                 .format(cnf_repo, cnf_name, tag, path_leaf(cnf['license'])))
        for v in cnf['volumes']:
            host_path, cnf_path = v.split(':')
            dc.write('      - {0}:{1}\n'.format(path_leaf(host_path),
                                                cnf_path))
            os.system('virt-customize -a {0} --copy-in {1}:/root/'
                      .format(args['vm_image'], host_path))
