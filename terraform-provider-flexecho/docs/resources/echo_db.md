---
page_title: "flexecho_echo_db Resource - terraform-provider-flexecho"
subcategory: ""
description: |-
  Creates an echo-DB clone on the Silk Flex Echo platform.
---

# flexecho_echo_db (Resource)

Creates (and removes) an echo-DB clone on the Silk Flex Echo platform.

Creation is **asynchronous**: the provider submits the clone job and polls the task to completion. There are two modes:

- **From a snapshot** — set `snapshot_id`; the clone is created from that snapshot.
- **Standalone replicate** — omit `snapshot_id` and set `source_host_id`; a snapshot-less replicate is performed.

The resource `id` is a composite of `destination_host_id/destination_db_id` (delete needs both). The Echo API has no get-by-id for clones, so this resource does not refresh remote state after create — there is no drift detection, and `import` is not supported. All configurable attributes force replacement.

## Example Usage

```terraform
# clone from an existing snapshot
resource "flexecho_echo_db" "orders_clone" {
  snapshot_id         = flexecho_db_snapshot.nightly.id
  destination_host_id = "sql02"
  destination_db_id   = "db-orders"
  destination_db_name = "orders_clone"
  target_state        = "online"
}
```

## Schema

### Required

- `destination_host_id` (String) Host id to place the clone on. Forces replacement.
- `destination_db_id` (String) Database id of the clone. Forces replacement.
- `destination_db_name` (String) Database name of the clone. Forces replacement.

### Optional

- `snapshot_id` (String) If set, clone FROM this snapshot. Otherwise a standalone replicate is performed. Forces replacement.
- `source_host_id` (String) Source host id (required for the standalone replicate path). Forces replacement.
- `target_state` (String) `recovery` or `online`. Defaults to `online`. Forces replacement.

### Read-Only

- `id` (String) Composite id, `destination_host_id/destination_db_id`.
