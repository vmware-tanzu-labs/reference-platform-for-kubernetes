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
module: kapp

short_description: This module runs Carvel's kapp tool.

# If this is part of a collection, you need to use semantic versioning,
# i.e. the version is of the form "2.5.0" and not "2.4".
version_added: "1.0.0"

description: This module runs Carvel's kapp to deploy Kubernetes manifests

options:
    name:
        description: The name of the App to modified by kapp
        required: true
        type: str
    state:
        description: The desired state of the App
        required: false
        type: str
        choices: ['present','absent']
        default: present
    kubeconfig:
        description: The location of the kubeconfig file to use
        required: false
    context:
        description: Override kubeconfig context
        required: false
    namespace:
        description: namespace to deploy application into
        required: true
    src:
        description: One or more files, directories, URLs to process (required when state=present)
        required: false
    kapp:
        description: Path to kapp binary (if its not in your $PATH)
        required: false
    wait:
        description: Wait until changes are applied
        required: false
        type: bool
        default: true
    wait_timeout:
        description: Time to wait for resources to apply
        required: false
        default: 2m0s

# Specify this value according to your collection
# in format of namespace.collection.doc_fragment_name
extends_documentation_fragment:
    - vmware.rpk.kapp

author:
    - Paul Czarkowski (@paulczar)
'''

EXAMPLES = r'''
# deploy an application to Kubernetes

- name: "ensure kapp app foo exists"
  kapp:
    state:        "present"
    src:          "/tmp/manifests/foo"
    namespace:    "foo-namespace"
    name:         "foo"
    wait:         false

# delete an application from Kubernetes
- name: "ensure kapp bar foo doesn't exist"
  kapp:
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
        state=dict(default='present', required=False, choices=['present','absent']),
        kubeconfig=dict(required=False),
        context=dict(required=False),
        name=dict(required=True),
        namespace=dict(required=True),
        src=dict(type='list', required=False),
        kapp=dict(required=False),
        wait=dict(type='bool', required=False, default=True),
        wait_timeout=dict(required=False, default='2m0s')
    )

    # seed the result dict in the object
    # we primarily care about changed and state
    # changed is if this module effectively modified the target
    # state will include any data that you want your module to pass back
    # for consumption, for example, in a subsequent task
    result = dict(
        changed=False,
        message='',
        rc=0, # return code
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
    # result['message'] = 'goodbye'

    args = []
    ## check that kapp binary exists
    kapp = module.params['kapp']
    if not kapp:
        kapp = module.get_bin_path("kapp", required=True)
    if not os.path.exists(kapp):
        result['message'] = kapp
        module.fail_json(msg='the specified kapp binary does not exist', **result)
    if not is_executable(kapp):
        result['message'] = kapp
        module.fail_json(msg='The specified kapp binary is not executable', **result)
    args.append(kapp)

    # Create or Destroy ?
    state = module.params['state']
    kapp_command = 'deploy'
    if state == 'absent':
        kapp_command = 'delete'
        # TEMP: ensure we don't kill antrea components. see https://github.com/vmware-tanzu/carvel-kapp/issues/160
        args.extend(["--filter", '{"and":[{"not":{"resource":{"kinds":["Antrea%"]}}}]}'])
    args.append(kapp_command)

    # Mandatory Args
    args.extend(["--yes", "--json", "--diff-changes"])
    args.extend(["--namespace", module.params['namespace']])
    args.extend(["--app", module.params['name']])

    # Optional args
    if module.params['wait']:
        args.extend(["--wait-timeout", module.params['wait_timeout']])
    else:
        args.append("--wait=false")
    if module.params['wait'] and module.params['wait_timeout']:
        args.extend([])
    if module.params['kubeconfig']:
        args.extend(["--kubeconfig", module.params['kubeconfig']])
    if module.params['context']:
        args.extend(["--kubeconfig-context", module.params['context']])


    # check that provided files exist.
    if kapp_command == 'deploy':
        if not module.params['src']:
            module.fail_json(msg='you must provide at least one value for "src"', **result)
        files = module.params['src']
        file_not_found = False
        for file in module.params['src']:
            if not os.path.exists(file):
                result['message'] = result['message'] + "the provided file " + file + " does not exist.\n"
                file_not_found = True
            args.extend(["-f", file])
        if file_not_found:
            module.fail_json(msg='One or more files not found', **result)
        # result['message'] = kapp + " " + " ".join(args)

    # Run it!
    print("Running " + " ".join(args))
    rc, out, err = module.run_command(args)

    if rc != 0:
        module.fail_json(msg="failed to run %s" % (" ".join(args)), rc=rc, err=err, out=out)

    parsed = json.loads(out)

    # attempt to process if changes happened, this involves some ugly text parsing
    if kapp_command == "delete":
        # check in parse lines to verify it wasn't already deleted
        if "App '%s' (namespace: %s) does not exist" % (
            module.params['name'], module.params['namespace']
            ) not in parsed['Lines']:
                result['changed'] = True

    if kapp_command == "deploy":
        # scan output for changes
        for line in parsed['Lines']:
            if re.match("applying \d+ changes", line):
                result['changed'] = True

    # import pprint
    # pp = pprint.PrettyPrinter(indent=4)
    # pp.pprint(parsed)

    result['rc'] = rc
    result['out'] = out
    result['err'] = err
    # result['message'] = parsed['Lines']
    # result['changed'] = True

    # during the execution of the module, if there is an exception or a
    # conditional state that effectively causes a failure, run
    # AnsibleModule.fail_json() to pass in the msg and the result
    # if module.params['name'] == 'fail me':
    #     module.fail_json(msg='You requested this to fail', **result)

    # in the event of a successful module execution, you will want to
    # simple AnsibleModule.exit_json(), passing the key/value results
    module.exit_json(**result)

# def execute_kapp():
#     kapp = module.get_bin_path("lvremove", required=True)
#     cmd = "%s %s -%s %s%s %s/%s %s" % (tool, test_opt, size_opt, size, size_unit, vg, this_lv['name'], pvs)
#                 rc, out, err = module.run_command(cmd)

def main():
    run_module()


if __name__ == '__main__':
    main()
