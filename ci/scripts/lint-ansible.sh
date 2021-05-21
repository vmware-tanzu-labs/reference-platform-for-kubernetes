#!/usr/bin/env bash
# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT

# set the exit code and adjust if there are failures
EXIT_CODE=0

echo 'linting ansible using ansible-lint rules at .ansible-lint'

# lint each role independently
for ROLE in `find roles/ -mindepth 1 -maxdepth 2 -type d`; do
  ansible-lint "${ROLE}/" -vvvvvvvvvvvvvv -R -r ./.ansible-lint-rules
  if [[ $? != 0 ]]; then
    EXIT_CODE=$?
  fi
done

# lint the top-level playbook
ansible-lint site.yaml -vvvvvvvvvvvvvv -R -r ./.ansible-lint-rules
if [[ $? != 0 ]]; then
  EXIT_CODE=$?
fi

exit ${EXIT_CODE}
