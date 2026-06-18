---
page_title: "flexecho_app_token Resource - terraform-provider-flexecho"
subcategory: ""
description: |-
  Manages a Silk Flex app token.
---

# flexecho_app_token (Resource)

Manages a Silk Flex app token (the core `/api/v2/flex_app_tokens` API, separate from the Echo endpoints). Create and delete are synchronous. Read walks the token list to match by id.

The TTL is always `-1` (non-expiring). All configurable attributes force replacement; `import` is not supported.

## Example Usage

```terraform
resource "flexecho_app_token" "ci" {
  app_name    = "ci-pipeline"
  description = "token for CI snapshot jobs"
}
```

## Schema

### Required

- `app_name` (String) Application name for the token. Forces replacement.

### Optional

- `description` (String) Free-text description. Forces replacement.

### Read-Only

- `expire_ts` (Number) Expiry timestamp.
- `id` (String) The token id.
- `ttl` (Number) Token time-to-live. Always `-1`.
- `valid` (Boolean) Whether the token is currently valid.
