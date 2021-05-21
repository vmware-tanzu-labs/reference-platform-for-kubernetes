# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
import ansiblelint.utils
from ansiblelint import AnsibleLintRule

class EnsureLowercaseNamedTasks(AnsibleLintRule):
    id = 'RPK001'
    shortdesc = 'Named tasks must be all lowercase'
    description = 'Named tasks must be all lowercase'
    tags = ['consistency','readability']

    def matchtask(self, file, task):
        # Task names should be lowercase
        if task.get('name'):
            if not task.get('name').islower():
                return True

        return False
