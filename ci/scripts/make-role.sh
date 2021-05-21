#!/usr/bin/env bash
# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT

if [ -z $RPK_ANSIBLE_ROLE ]; then
	echo "No role name was provided."
	echo "Usage is: ./make-role.sh <new-role-name>"
else
	ansible-galaxy init --role-skeleton=/ansible/roles/support/role-skeleton $RPK_ANSIBLE_ROLE
    mv $RPK_ANSIBLE_ROLE /ansible/roles/$RPK_ANSIBLE_ROLE
fi
