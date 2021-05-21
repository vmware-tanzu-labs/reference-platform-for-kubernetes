# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
#!/usr/bin/env python

import argparse
import os
import sys
import yaml


def get_components(profile_file):
    """Retrieve comonents to install from profile file"""

    components = []

    with open(profile_file, 'r') as f:
        try:
            for c in yaml.safe_load(f)['rpk_components']:
                if c['enabled'] == True:
                    components.append(c['name'])
        except yaml.YAMLError as e:
            print(e)
            sys.exit(1)

    return components

def install_component(component):
    """Run make target to install RPK component with ytt and kapp"""

    cmd = 'ROLE={} make install.role'.format(component)

    resp = os.system(cmd)

    print(resp)

if __name__ == '__main__':
    parser = argparse.ArgumentParser('Install with carvel')
    parser.add_argument('-p', '--profile', type=str, default='profiles/platform.yaml',
            help='path to profiles file with roles to install')

    args = parser.parse_args()
    print('Installing RPK components in profile file: {}'.format(args.profile))

    components = get_components(args.profile)
    print('Components to install: {}'.format(components))

    for c in components:
        install_component(c)

    print('RPK components installed')

