# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
import ansiblelint.utils
from ansiblelint import AnsibleLintRule

class EnsureNoChange(AnsibleLintRule):
    id = 'RPK003'
    shortdesc = 'only report changes on k8s and infrastructure providers'
    description = 'modules which modify files should not report a change; only kubernetes changes or other api changes should report a change'
    tags = ['consistency','readability']

    def matchtask(self, file, task):
        target_modules = ['template','file','lineinfile']

        # check for relevant modules and ensure they use the changed when
        if task["action"]["__ansible_module__"] in target_modules:
            if task.get('changed_when') is None:
                return True
            if task.get('changed_when') is not False:
                return True

        return False
