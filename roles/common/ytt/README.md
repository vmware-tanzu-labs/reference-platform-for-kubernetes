# ytt

An Ansible role to template files with ytt (yaml templating tool).

## Variables

### Input Variables

The following variables are common input variables:

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| component | Name of the component to template.  In most cases this will be the component name.  Refers to kapp application name. | "" | string | yes |
| values | List of values to write to the ytt values file. Will default to the play's variables if none provided. | [] | list | no |

### Default Variables

The following variables are defaulted in `vars/main.yaml` and require no additional user input.  **WARN:** change with caution:

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| output_dir | Directory where manifests end up | "{{ staging_dir }}/manifests" | string (path) | no |
| discovered_values | Values discovered during pre-flight processing available to YTT via `data.values.discovered_values` | "{{ rpk_staging_dir }}/discovered-values.yaml" | string (path) | no |
| component_base | The component base (raw, unmodified YAML) directory output location. | "{{ staging_dir }}/base" | string (path) | no |
| component_overlays | The component overlays directory output location. | "{{ staging_dir }}/overlays" | string (path) | no |
| component_templates | The component templates directory output location. | "{{ staging_dir }}/templates" | string (path) | no |
| component_values | The component values directory output location. | "{{ staging_dir }}/values" | string (path) | no |
| common_files | Common base, overlays, templates, and values which are used by all components. | See `vars/main.yaml` | list | no |

```yaml
  - name: "generate manifests with ytt"
    import_role:
      name: "common/ytt"
    vars:
      component: "tanzu-storage"
      values:
        - tanzu_storage:    "{{ tanzu_storage }}"
        - component_values: "{{ tanzu_storage }}"
```


## Author

[Paul Czarkowski](mailto:pczarkowski@vmware.com)

## Editor(s)

- [Dustin Scott](@mailto:sdustin@vmware.com)
