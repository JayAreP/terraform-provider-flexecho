---
page_title: "flexecho_host Data Source - terraform-provider-flexecho"
subcategory: ""
description: |-
  Reads a single registered host from the Silk Flex Echo platform.
---

# flexecho_host (Data Source)

Reads a single registered host by id. Useful for referencing an existing host (one not managed by this configuration) — the read errors if the host does not exist.

## Example Usage

```terraform
data "flexecho_host" "sql01" {
  host_id = "sql01"
}

resource "flexecho_db_snapshot" "nightly" {
  source_host_id = data.flexecho_host.sql01.host_id
  database_ids   = ["SilkEDW"]
}
```

## Schema

### Required

- `host_id` (String) The host id to look up.

### Read-Only

- `agent_version` (String)
- `db_engine_version` (String)
- `db_vendor` (String)
- `host_name` (String)
- `id` (String) The host id.
- `is_connected` (Boolean)
