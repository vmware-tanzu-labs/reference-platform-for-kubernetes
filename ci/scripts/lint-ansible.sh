#!/usr/bin/env bash
# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT

echo 'linting ansible using ansible-lint rules at .ansible-lint'

# lint each role independently
for ROLE in `find roles/ -mindepth 1 -maxdepth 2 -type d`; do
  ansible-lint "${ROLE}/" -vvvvvvvvvvvvvv -R -r ./.ansible-lint-rules
done

# lint the top-level playbook
ansible-lint site.yaml -vvvvvvvvvvvvvv -R -r ./.ansible-lint-rules
