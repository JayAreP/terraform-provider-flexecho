---
page_title: "flexecho Provider"
subcategory: ""
description: |-
  Terraform provider for the Silk Flex "Echo" API — database snapshot and echo-DB clone management.
---

# flexecho Provider

The `flexecho` provider manages Silk Flex "Echo" objects against a Silk Flex deployment — registered hosts, database snapshots, echo-DB clones, and app tokens.

It authenticates with a bearer token and talks to the Flex API over HTTPS (the management certificate is self-signed, so TLS verification is skipped).

## Example Usage

```terraform
terraform {
  required_providers {
    flexecho = {
      source  = "localdomain/provider/flexecho"
      version = "0.1.2"
    }
  }
}

provider "flexecho" {
  server = "10.0.215.133"   # or set SILK_FLEX_SERVER
  token  = var.flex_token   # or set SILK_FLEX_TOKEN
}

variable "flex_token" {
  type      = string
  sensitive = true
}
```

## Authentication

Both `server` and `token` may be supplied in the provider block or through environment variables:

- `SILK_FLEX_SERVER`
- `SILK_FLEX_TOKEN`

## Schema

### Required

- `server` (String) IP address or hostname of the Silk Flex management console. May also be set with the `SILK_FLEX_SERVER` environment variable.
- `token` (String, Sensitive) Bearer token used to authenticate against the Flex API. May also be set with the `SILK_FLEX_TOKEN` environment variable.

## Async operations

`flexecho_db_snapshot` and `flexecho_echo_db` are backed by asynchronous jobs — create (and delete) submit a task that the provider polls to completion (states `pending`/`running` → `completed`, or `failed`/`aborted`). `flexecho_host` and `flexecho_app_token` are synchronous.
