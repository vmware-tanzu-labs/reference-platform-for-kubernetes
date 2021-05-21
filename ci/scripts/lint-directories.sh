#!/usr/bin/env bash
# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT

LINT_DIRECTORY_EXCLUSIONS="
  ./.git
  ./roles/support/role-skeleton
  ./ansible-virtualenv
  ./.ansible-lint-rules/__pycache__
  ./build"

EXCLUSION_SYNTAX=""
for DIR in ${LINT_DIRECTORY_EXCLUSIONS}; do
  EXCLUSION_SYNTAX="${EXCLUSION_SYNTAX} -and -path ${DIR} -prune"
done

VIOLATIONS=`find . -name "*_*" -type d ${EXCLUSION_SYNTAX}`

if [[ ${VIOLATIONS} == "" ]]; then
  echo "no directories with underscores detected"
  exit 0
else
  echo "the following directories have underscores and do not meet validation criteria: ${VIOLATIONS}"
  exit 1
fi
