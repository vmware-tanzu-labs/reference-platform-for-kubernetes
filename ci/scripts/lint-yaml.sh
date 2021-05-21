#!/usr/bin/env bash
# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT

echo 'linting yaml using yamllint rules at .yamllint'
yamllint . --strict
