
# Spiff-based Transfer Handler

This package provides a `TransferHandler` programmable using the
[Spiff Templating Engine](https://github.com/mandelsoft/spiff).
The spiff template get a flat binding `mode` describing the operation mode
and it must provide a top-level node `process` containing
the processing result.

The following modes are used:

## Resource Mode

This mode is used to decide on the by-value option for a resource. It gets the
following bindings:

- `mode` *&lt;string>*: `resource`
- `values` *&lt;map>*:
  - `component` *&lt;map>*:  the meta dats of the component version carrying the resource
    - `name` *&lt;string>*: component name
    - `version` *&lt;string>*: component version
    - `provider` *&lt;string>*: provider name
    - `labels` *&lt;map[string]>*: labels of the component version (deep)
  - `element` *&lt;map>*:  the meta data of the resource
  - `access` *&lt;map>*:  the access specification of the resource

The result value (field `process`) must be a boolean describing whether the
resource should be transported ny-value.

## Source Mode

This mode is used to decide on the by-value option for a source. It gets the
following bindings:

- `mode` **&lt;string>**: `resource`
- `values` **&lt;map>**: (see [Resource Mode](#resource-mode))

The result value (field `process`) must be a boolean describing whether the
resource should be transported ny-value.

## Component Version Mode

This mode is used to decide on the recursion option for a referenced component
version. It gets the  following bindings:

- `mode` **&lt;string>**: `componentversion`
- `values` **&lt;map>**: (see [Resource Mode](#resource-mode))
  - `component` **&lt;map>**:  the meta dats of the component version carrying the reference
    - `name` **&lt;string>**: component name
    - `version` **&lt;string>**: component version
    - `provider` **&lt;string>**: provider name
    - `labels` **&lt;map[string]>**: labels of the component version (deep)
  - `element` **&lt;map>**:  the meta data of the component reference 

The result value (field `process`) can either be a simple boolean value
or a map with the following fields:

- `process` **&lt;bool>**: `true` indicates to follow the reference
- `repospec` **&lt;map>** *(optional)*: the specification of the repository to use to follow the reference
- `script` **&lt;template>** *(optional)*: the script to use instead of the current one.

If no new repository spec is given, the actual repository is used. If no new
script is given, the actual one is used for sub sequent processing.