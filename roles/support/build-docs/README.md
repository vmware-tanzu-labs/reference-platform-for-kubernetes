# build-docs

An Ansible role to assemble documentation for RPK from common documentation fragments.

* Builds the component documentation tables, which gets inserted into the cloud provider QUICKSTART documentation.  This uses the `rpk_components` data for the `rpk_supported_profiles` variable in [defaults/main.yaml](defaults/main.yaml).

* Builds the main cloud provider QUICKSTART documentation found in [/docs/providers](../../../docs/providers) from the templates and common sections in [templates/sections/main](templates/sections/main) and [templates/sections/common](templates/sections/common) respectively.

* Builds table of contents for the cloud provider QUICKSTART documentation based off of the headings found from the assembled main documentation in [templates/built](templates/built).  The table of contents template can be found in [templates/toc](templates/toc/toc.md.j2).

* Assembles the final version of the documentation with title, table of contents, and content from the provider templates in [templates](templates) with generated content from [templates/built](templates/built).


## Author

[Andrew J. Huffman](mailto:ahuffman@vmware.com)
