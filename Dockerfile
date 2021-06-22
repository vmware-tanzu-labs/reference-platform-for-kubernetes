# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
#
# RUN
#
FROM projects.registry.vmware.com/rpk/rpk-base:v1.4.3

# update the image
RUN apk update && apk add bash bind-tools openssl --no-cache

# configure user for runtime environment
# NOTE: because we mount /etc/hosts to ensure DNS resolution during the deployment process,
#       we cannot use the ansible user since it cannot write to /etc/hosts.  leaving the below
#       in place as a comment so we can revert to this behavior once the DNS resolution problem
#       is resolved.
# RUN addgroup -S ansible
# RUN adduser -S -G ansible ansible

# copy dependencies from build image
# NOTE: because we mount /etc/hosts to ensure DNS resolution during the deployment process,
#       we cannot use the ansible user since it cannot write to /etc/hosts.  leaving the below
#       in place as a comment so we can revert to this behavior once the DNS resolution problem
#       is resolved.
# COPY --from=build --chown ansible:ansbile /ansible /ansible

# copy the rest of the directory
COPY . /ansible

# set pathing
ENV PATH=/ansible/.pip/usr/local/bin:$PATH
ENV PYTHONPATH=/usr/local/lib/python3.7/site-packages/:/ansible/.pip/usr/local/lib/python3.7/site-packages/

# set kubeconfig and ansible options
ENV KUBECONFIG=/ansible/.kube/config
ENV ANSIBLE_FORCE_COLOR=1

# set runtime user
# NOTE: because we mount /etc/hosts to ensure DNS resolution during the deployment process,
#       we cannot use the ansible user since it cannot write to /etc/hosts.  leaving the below
#       in place as a comment so we can revert to this behavior once the DNS resolution problem
#       is resolved.
# USER ansible:ansible

WORKDIR /ansible
ENTRYPOINT ["./bin/rpk"]
