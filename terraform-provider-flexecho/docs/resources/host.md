---
page_title: "flexecho_host Resource - terraform-provider-flexecho"
subcategory: ""
description: |-
  Registers a host with the Silk Flex Echo platform.
---

# flexecho_host (Resource)

Registers (and removes) a host with the Silk Flex Echo platform. Registration is synchronous.

The `token` returned at creation is the agent token, and it is only available once — at create time. All configurable attributes force replacement; the resource has no in-place update.

## Example Usage

```terraform
resource "flexecho_host" "sql01" {
  host_id   = "sql01"
  db_vendor = "mssql"
}
```

## Schema

### Required

- `host_id` (String) The host id. Must be 3-32 chars and start with a letter. Forces replacement.

### Optional

- `db_vendor` (String) Database vendor: `mssql` or `oracledb`. Defaults to `mssql`. Forces replacement.
- `sdp_id` (String) Optional SDP id to associate with the host. Forces replacement.

### Read-Only

- `agent_version` (String)
- `db_engine_version` (String)
- `host_name` (String)
- `id` (String) The host id.
- `is_connected` (Boolean)
- `token` (String, Sensitive) Agent token returned at host creation. Only populated on the initial create.

## Import

Import is supported using the host id:

```shell
terraform import flexecho_host.sql01 sql01
```
