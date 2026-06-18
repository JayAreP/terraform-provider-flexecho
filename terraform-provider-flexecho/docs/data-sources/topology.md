---
page_title: "flexecho_topology Data Source - terraform-provider-flexecho"
subcategory: ""
description: |-
  Returns the Silk Flex Echo topology as a JSON string.
---

# flexecho_topology (Data Source)

Returns the Echo topology (host, SDPs, and databases) as a JSON string. The nested shape is large and loosely typed, so it is exposed raw for you to parse with `jsondecode()`.

## Example Usage

```terraform
data "flexecho_topology" "current" {}

output "topology" {
  value = jsondecode(data.flexecho_topology.current.json)
}
```

## Schema

### Read-Only

- `id` (String) The ID of this data source.
- `json` (String) The topology, serialized as a JSON string.
