# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
import ansiblelint.utils
import re
from ansiblelint import AnsibleLintRule

class EnsureKubectlContext(AnsibleLintRule):
    id = 'RPK002'
    shortdesc = 'k8s modules and commands must use a kubectl context'
    description = 'k8s modules and commands must use a kubectl context'
    tags = ['consistency','safety']

    def matchtask(self, file, task):
        k8s_modules = ['k8s','k8s_auth','k8s_facts','k8s_info','k8s_raw','k8s_scale','k8s_service']

        # check for relevant modules and ensure they use the context key
        if task["action"]["__ansible_module__"] in k8s_modules:
            if task['action'].get('context') is None:
                return True

        # additionally, check the shell module for a kubectl command and make sure we use
        # --context with it
        # TODO: this does not quite work yet
        if task["action"]["__ansible_module__"] == "shell":
            raw_command = self.unjinja(' '.join(task["action"].get("__ansible_arguments__", [])))
            regex_kubectl = re.compile('kubectl')
            regex_context = re.compile('--context')

            if regex_kubectl.search(raw_command):
                if not regex_context.search(raw_command):
                    return True

        return False
