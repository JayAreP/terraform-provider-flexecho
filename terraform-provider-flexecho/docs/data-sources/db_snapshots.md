---
page_title: "flexecho_db_snapshots Data Source - terraform-provider-flexecho"
subcategory: ""
description: |-
  Lists database snapshots on the Silk Flex Echo platform.
---

# flexecho_db_snapshots (Data Source)

Lists database snapshots, optionally narrowed to a single host.

## Example Usage

```terraform
# all snapshots
data "flexecho_db_snapshots" "all" {}

# snapshots for one host
data "flexecho_db_snapshots" "for_host" {
  host_id = "sql-silk-edw"
}
```

## Schema

### Optional

- `host_id` (String) If set, only return snapshots for this host. Omit to return all.

### Read-Only

- `id` (String) The ID of this data source.
- `snapshots` (List of Object) The matching snapshots. (see [below for nested schema](#nestedatt--snapshots))

<a id="nestedatt--snapshots"></a>
### Nested Schema for `snapshots`

Read-Only:

- `consistency_level` (String)
- `host_id` (String)
- `host_name` (String)
- `id` (String) The snapshot id.
- `timestamp` (Number)
