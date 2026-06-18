---
page_title: "flexecho_hosts Data Source - terraform-provider-flexecho"
subcategory: ""
description: |-
  Lists all registered hosts on the Silk Flex Echo platform.
---

# flexecho_hosts (Data Source)

Lists all hosts registered with the Silk Flex Echo platform.

## Example Usage

```terraform
data "flexecho_hosts" "all" {}

output "all_hosts" {
  value = data.flexecho_hosts.all.hosts
}
```

## Schema

### Read-Only

- `hosts` (List of Object) The registered hosts. (see [below for nested schema](#nestedatt--hosts))
- `id` (String) The ID of this data source.

<a id="nestedatt--hosts"></a>
### Nested Schema for `hosts`

Read-Only:

- `db_vendor` (String)
- `host_id` (String)
- `host_name` (String)
- `is_connected` (Boolean)
