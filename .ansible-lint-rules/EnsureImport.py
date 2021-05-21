# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
import ansiblelint.utils
from ansiblelint import AnsibleLintRule

class EnsureImport(AnsibleLintRule):
    id = 'RPK005'
    shortdesc = 'prefer import_[role,tasks,playbook] over include unless looping'
    description = 'prefer import_[role,tasks,playbook] over include unless looping'
    tags = ['consistency','readability','performance']

    def matchtask(self, file, task):
        target_modules = ['include_playbook','include_role','include_tasks']
        target_loops = ['loop','with_items','with_dict','with_fileglob','with_list']

        # check for relevant modules and ensure they use import_* if not looping
        if task["action"]["__ansible_module__"] in target_modules:
            failed = True
            for loop_var in target_loops:
                if task.get(loop_var) is None:
                    continue
                else:
                    failed = False
                    break
            return failed
        return False
