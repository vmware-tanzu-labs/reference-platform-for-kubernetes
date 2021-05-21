# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
#!/usr/bin/env python

# Copyright: (c) 2018, Terry Jones <terry.jones@example.org>
# Copyright: (c) 2020, Paul Czarkowski <pczarkowski@vmware.com>
# GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
from __future__ import (absolute_import, division, print_function)
__metaclass__ = type

DOCUMENTATION = r'''
---
module: ytt

short_description: This module runs Carvel's ytt tool.

# If this is part of a collection, you need to use semantic versioning,
# i.e. the version is of the form "2.5.0" and not "2.4".
version_added: "1.0.0"

description: This module runs Carvel's ytt to template Kubernetes manifests

options:
    src:
        description: One or more files or directories to process with ytt
        required: false
    output_directory:
        description: The path to place processed files
        required: true
    ytt:
        description: Path to ytt binary (if its not in your $PATH)
        required: false
    ignore_unknown_comments:
        description: "Configure whether unknown comments are considered as errors (comments that do not start with '#@' or '#!')"
        required: false
        default: true
    values:
        description: a list of values
        required: false
    file_marks:
        description: a list of file marks to apply (see ytt template -f - --file-mark)
        required: false

# Specify this value according to your collection
# in format of namespace.collection.doc_fragment_name
extends_documentation_fragment:
    - vmware.rpk.ytt

author:
    - Paul Czarkowski (@paulczar)
'''

EXAMPLES = r'''
# deploy an application to Kubernetes

- name: "ensure ytt app foo exists"
  ytt:
    state:        "present"
    src:          "/tmp/manifests/foo"
    namespace:    "foo-namespace"
    name:         "foo"
    wait:         false

# delete an application from Kubernetes
- name: "ensure ytt bar foo doesn't exist"
  ytt:
    state:        "absent"
    namespace:    "bar-namespace"
    name:         "bar"
    wait:         true
    wait_timeout: "5m"
'''

RETURN = r'''
# These are examples of possible return values, and in general should use other names for return values.
msg='',
rc=0, # return code
out='',  # stdout
err='', # stderr
'''

from ansible.module_utils.basic import *
import re


def run_module():
    # define available arguments/parameters a user can pass to the module
    module_args = dict(
        src=dict(type='list', required=True),
        output_directory=dict(required=True),
        values=dict(type='list', required=False, default=[]),
        ytt=dict(required=False),
        ignore_unknown_comments=dict(type='bool', required=False, default=True),
        file_marks=dict(type='list', required=False, default=[])
    )

    # seed the result dict in the object
    # we primarily care about changed and state
    # changed is if this module effectively modified the target
    # state will include any data that you want your module to pass back
    # for consumption, for example, in a subsequent task
    result = dict(
        changed=False,
        message='',
        # rc=0, # return code
        out='',  # stdout
        err='', # stderr
    )

    # the AnsibleModule object will be our abstraction working with Ansible
    # this includes instantiation, a couple of common attr would be the
    # args/params passed to the execution, as well as if the module
    # supports check mode
    module = AnsibleModule(
        argument_spec=module_args,
        supports_check_mode=True
    )

    # if the user is working with this module in only check mode we do not
    # want to make any changes to the environment, just return the current
    # state with no modifications
    if module.check_mode:
        module.exit_json(**result)

    # manipulate or modify the state as needed (this is going to be the
    # part where your module will do what it needs to do)
    # result['original_msg'] = module.params['name']
    # result['msg'] = 'goodbye'

    args = []
    ## check that ytt binary exists
    ytt = module.params['ytt']
    if not ytt:
        ytt = module.get_bin_path("ytt", required=True)
    if not os.path.exists(ytt):
        module.fail_json(msg='the specified ytt binary does not exist', **result)
    if not is_executable(ytt):
        module.fail_json(msg='The specified ytt binary is not executable', **result)
    args.append(ytt)

    # check that provided source files exist.
    files = module.params['src']
    file_not_found = False
    for file in module.params['src']:
        if not os.path.exists(file):
            result['message'] = result['message'] + "the provided file " + file + " does not exist.\n"
            file_not_found = True
        args.extend(["-f", file])
    if file_not_found:
        module.fail_json(msg='One or more files not found', **result)

    # check the destination directory exists
    dest = module.params['output_directory']
    if not os.path.exists(dest):
        result['message'] = "the provided destination directory " + dest + " does not exist.\n"
        module.fail_json(msg='Error processing args', **result)
    if not os.path.isdir(dest):
        result['message'] = "the provided destination directory " + dest + " is a file.\n"
        module.fail_json(msg='Error processing args', **result)
    args.extend(["--output-files", dest])

    if module.params['ignore_unknown_comments']:
        args.append("--ignore-unknown-comments=true")

    if len(module.params['values']) > 0:
        for value in module.params['values']:
            args.extend(['--data-value', value])

    if len(module.params['file_marks']) > 0:
        for mark in module.params['file_marks']:
            args.extend(['--file-mark', mark])

    # Run it!
    # import pprint
    # pp = pprint.PrettyPrinter(indent=4)
    # pp.pprint(args)
    # print("Running " + " ".join(args))
    rc, out, err = module.run_command(args)

    if rc != 0:
        module.fail_json(msg="failed to run %s" % (" ".join(args)), rc=rc, err=err, out=out)

    result['rc'] = rc
    result['out'] = out
    result['err'] = err
    result['changed'] = True

    # during the execution of the module, if there is an exception or a
    # conditional state that effectively causes a failure, run
    # AnsibleModule.fail_json() to pass in the msg and the result
    # if module.params['name'] == 'fail me':
    #     module.fail_json(msg='You requested this to fail', **result)

    # in the event of a successful module execution, you will want to
    # simple AnsibleModule.exit_json(), passing the key/value results
    module.exit_json(**result)

def main():
    run_module()


if __name__ == '__main__':
    main()
