---
page_title: "flexecho_echo_db Resource - terraform-provider-flexecho"
subcategory: ""
description: |-
  Replicates (clones) databases from a source host to one or more destinations.
---

# flexecho_echo_db (Resource)

Replicates databases from a source host out to one or more destination hosts (the `/echo_dbs` standalone replicate).

Creation is **asynchronous**: the provider submits the replicate job and polls the task to completion. Destinations are given as repeated `destination` blocks; on delete, the provider tears down each destination clone individually. The Echo API has no get-by-id for clones, so this resource does not refresh remote state after create — there is no drift detection, and `import` is not supported. All attributes force replacement.

## Example Usage

```terraform
resource "flexecho_echo_db" "clone" {
  source_host_id = "sql-silk-edw"
  database_ids   = ["SilkEDW"]

  destination {
    host_id = "sql-silk-prod"
    db_id   = "SilkEDW"
    db_name = "SilkEDW-test"
  }

  destination {
    host_id = "sql-silk-uat"
    db_id   = "SilkEDW"
    db_name = "SilkEDW-test"
  }

  name_prefix       = "snap"
  consistency_level = "application"
  target_state      = "online"
}
```

## Schema

### Required

- `source_host_id` (String) Source host id the databases live on. Forces replacement.
- `database_ids` (List of String) Database ids to replicate from the source host. Forces replacement.
- `destination` (Block List, min: 1) One or more clone destinations. (see [below for nested schema](#nestedblock--destination)) Forces replacement.

### Optional

- `consistency_level` (String) `crash` or `application`. Defaults to `application`. Forces replacement.
- `name_prefix` (String) Snapshot name prefix (4-20 chars, `^[a-z][a-z0-9_-]+$`). Defaults to `snap`. Forces replacement.
- `use_vss` (Boolean) Use VSS for the capture. Defaults to `false`. Forces replacement.
- `target_state` (String) `recovery` or `online`. Defaults to `online`. Forces replacement.
- `backup_session_timeout_sec` (Number) Optional backup session timeout, in seconds. Forces replacement.
- `restore_session_timeout_sec` (Number) Optional restore session timeout, in seconds. Forces replacement.

### Read-Only

- `id` (String) Composite id derived from the source host and destinations.

<a id="nestedblock--destination"></a>
### Nested Schema for `destination`

Required:

- `host_id` (String) Destination host id.
- `db_id` (String) Destination database id.
- `db_name` (String) Destination database name.
