---
page_title: "flexecho_db_snapshot Resource - terraform-provider-flexecho"
subcategory: ""
description: |-
  Captures a database snapshot on the Silk Flex Echo platform.
---

# flexecho_db_snapshot (Resource)

Captures a database snapshot on the Silk Flex Echo platform.

Creation is **asynchronous**: the provider submits the capture job and polls the task to completion before returning. The resource `id` is the resolved snapshot id (e.g. `snap_<epoch>`). All configurable attributes force replacement; the resource has no in-place update.

## Example Usage

```terraform
resource "flexecho_db_snapshot" "nightly" {
  source_host_id    = "sql-silk-edw"
  database_ids      = ["SilkEDW"]
  name_prefix       = "snap"
  consistency_level = "application"
  use_vss           = false
}
```

## Schema

### Required

- `database_ids` (List of String) Database ids to capture in the snapshot. Forces replacement.
- `source_host_id` (String) Host id the databases live on. Forces replacement.

### Optional

- `consistency_level` (String) `crash` or `application`. Defaults to `application`. Forces replacement.
- `name_prefix` (String) Snapshot name prefix (4-20 chars, `^[a-z][a-z0-9_-]+$`). Defaults to `snap`. Forces replacement.
- `use_vss` (Boolean) Use VSS for the capture. Defaults to `false`. Forces replacement.

### Read-Only

- `host_name` (String)
- `id` (String) The snapshot id (e.g. `snap_<epoch>`).
- `timestamp` (Number)

## Import

Import is supported using the snapshot id:

```shell
terraform import flexecho_db_snapshot.nightly snap_1781794269
```
