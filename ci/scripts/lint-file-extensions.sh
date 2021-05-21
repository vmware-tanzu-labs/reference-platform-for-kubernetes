#!/usr/bin/env bash
# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT

LINT_EXTENSION_EXCLUSIONS="./ansible_virtualenv"

EXCLUSION_SYNTAX=""
for EXT in ${LINT_EXTENSION_EXCLUSIONS}; do
  EXCLUSION_SYNTAX="${EXCLUSION_SYNTAX} -and -path ${EXT} -prune"
done

VIOLATIONS=`find . -name "*.yml*" -type f ${EXCLUSION_SYNTAX}`

if [[ ${VIOLATIONS} == "" ]]; then
  echo "no yaml versus yml violations detected"
  exit 0
else
  echo "the following files have use a .yml extension as opposed to a .yaml extension and do not meet validation criteria: ${VIOLATIONS}"
  exit 1
fi
