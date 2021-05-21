#!/usr/bin/env bash
# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT

# set the exit code and adjust if there are failures
EXIT_CODE=0
BASE_PLAYBOOK_LIST="site.yaml"

echo 'linting ansible using ansible-lint rules at .ansible-lint'

# lint each role independently
for ROLE in $BASE_PLAYBOOK_LIST `find roles/ -mindepth 1 -maxdepth 2 -type d`; do
  ansible-lint "${ROLE}/" -vvvvvvvvvvvvvv -R -r ./.ansible-lint-rules
  RC=$?
  if [ $RC -ne 0 ]; then
    EXIT_CODE=$RC
  fi
done

echo "exiting with code: ${EXIT_CODE}"

exit ${EXIT_CODE}
