# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
import re
import sys
import os
import yaml

def _is_true(check):
    if type(check) is str:
        if check.lower() in ['true', '1', 't', 'y', 'yes']:
            return True
    if type(check) is bool:
        if check:
            return True
    return False

def rpk_component_is_enabled(role, components):
    for component in rpk_component_list(components):
        if component == role:
            return True
    return False

def rpk_component_list(components):
    list = []
    for index in range(len(components)):
        if _is_true(components[index]["enabled"]):
            list.append(components[index]["name"])
    return list

def eskape(string, fix=True):
    fail = False
    if fix:
        string = string.replace(".", "-")
        string = string.replace("_", "-")
        string = string.lower()
    # https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
    # contain at most 63 characters
    if len(string) > 63:
        if fix:
            string = string[:63]
        else:
            print(string + " must contain at most 63 characters")
    # contain only lowercase alphanumeric characters or '-'
    reg = re.compile("^[a-z0-9\-]+$")
    if not reg.match(string):
        print(string + " must contain only lowercase alphanumeric characters or '-'")
        fail = True
    # start with an alphanumeric character
    if string[0] == "-":
        print(string + " must start with an alphanumeric character")
        fail = True
    # end with an alphanumeric character
    if string[-1] == "-":
        print(string + " must end with an alphanumeric character")
        fail = True
    if fail:
        print("Reference https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names")
        raise
    return string

def _get_dependencies(component_name):
    # these paths are not required to have a dependencies file and are skipped
    blacklisted_paths = ["support", "common"]

    component_path = ""
    component_base_dir = component_name.split("/")[-1]
    for root, dirs, files in os.walk(os.getcwd() + "/roles/", topdown=False):
        for name in dirs:
            if name == component_base_dir:
                component_path = os.path.join(root, component_base_dir + "/")

                # return an empty dict if the path is blacklisted
                for path in blacklisted_paths:
                    if path in component_path:
                        return {"dependencies": []}

                break

    dependency_file = component_path + ".dependencies.yaml"

    with open(dependency_file) as file:
        dependency_yaml = yaml.safe_load(file)

    return dependency_yaml

# resolves dependencies from an individual component
def rpk_resolve_deps(component, resolve=True):
    # if we don't need to resolve, return the component in it's original form
    if not resolve:
        return component

    # quick pre-validation
    #  NOTE: just return the original value back to the user if we don't have what we need
    if type(component) is not list:
        return component
    if not len(component) == 1:
        return component
    if component[0]['name'] == None:
        return component
    if component[0]['enabled'] == None:
        return component

    # get the dependencies from the calling module and set the children dependencies
    metadata = _get_dependencies(component[0]['name'])
    dependencies = metadata['dependencies']

    # order dependencies by priority (lowest to highest)
    dependencies = sorted(dependencies, key = lambda i: i['priority'])

    # add required fields to each dependency
    for dep in dependencies:
        dep['name'] = dep['component']
        dep['enabled'] = True

    # add in the original component at the end of the dictionary and return
    dependencies.append(component[0])
    return dependencies

class FilterModule(object):
    ''' rpk custom filters '''

    def filters(self):
        return {
            'eskape': eskape,
            'rpk_resolve_deps': rpk_resolve_deps,
            'rpk_component_is_enabled': rpk_component_is_enabled,
            'rpk_component_list': rpk_component_list,
        }
