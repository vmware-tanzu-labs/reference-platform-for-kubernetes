# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
import ansiblelint.utils
from ansiblelint import AnsibleLintRule

class EnsureDebugVerbosity(AnsibleLintRule):
    id = 'RPK004'
    shortdesc = 'ensure the verbosity flag is defined when debugging'
    description = 'ensure the verbosity flag is defined when debugging'
    tags = ['consistency','readability']

    def matchtask(self, file, task):
        target_modules = ['debug']

        # check for relevant modules and ensure they use the changed when
        if task["action"]["__ansible_module__"] in target_modules:
            if task['action'].get('verbosity') is None:
                return True

        return False
