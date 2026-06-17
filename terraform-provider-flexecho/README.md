# terraform-provider-flexecho

Terraform provider for the Silk Flex **Echo** API — manage hosts, DB snapshots,
echo-db clones, and app tokens. Wraps [`flexecho-go-sdk`](../flexecho-go-sdk).

Named `flexecho` so it can publish to the Terraform registry as
`registry.terraform.io/silk-us/flexecho` (adjust the namespace to your registry org).

## Provider config

```hcl
provider "flexecho" {
  server = "20.84.120.30"   # or SILK_FLEX_SERVER
  token  = var.flex_token   # or SILK_FLEX_TOKEN  (sensitive)
}
```

## Resources

| Resource | Endpoint(s) | Notes |
|---|---|---|
| `flexecho_host` | PUT/GET/DELETE `/flex/api/v1/hosts/{id}` | sync |
| `flexecho_db_snapshot` | POST `/api/echo/v1/db_snapshots`, DELETE `/{id}` | async (task poll) |
| `flexecho_echo_db` | POST `/echo_dbs` or `/db_snapshots/{id}/echo_db`; DELETE `/echo_dbs` | async; clone from snapshot or standalone |
| `flexecho_app_token` | GET/POST/DELETE `/api/v2/flex_app_tokens` | core flex api |

## Data sources

`flexecho_host`, `flexecho_hosts`, `flexecho_db_snapshots`, `flexecho_topology`.

## Build & install locally

Not published to the registry yet, so build it and install into the local plugin
directory under the `localdomain/provider/flexecho` namespace. The build scripts
build against the local [`../flexecho-go-sdk`](../flexecho-go-sdk) checkout and
restore `go.mod` on exit.

```powershell
# Windows
.\build-local.ps1 -Install      # build + install to %APPDATA%\terraform.d\plugins\...
.\build-local.ps1               # just build for the host into .\bin
.\build-local.ps1 -All          # build every release target
```

```bash
# linux / macOS
./build-local.sh --install      # build + install to ~/.terraform.d/plugins/...
./build-local.sh                # host build into ./bin
./build-local.sh --all
```

Install path:

```
# Windows
%APPDATA%\terraform.d\plugins\localdomain\provider\flexecho\<version>\<os>_<arch>\

# linux / macOS
~/.terraform.d/plugins/localdomain/provider/flexecho/<version>/<os>_<arch>/
```

Then reference the local namespace in your `.tf` and run `terraform init`:

```hcl
terraform {
  required_providers {
    flexecho = {
      source  = "localdomain/provider/flexecho"
      version = "0.1.0"
    }
  }
}
```

`terraform init` will pick up the local build (it prints an "unauthenticated"
notice — expected for a local provider). Once published, swap the source to
`silk-us/flexecho`. See [`examples/main.tf`](examples/main.tf).

## Status

First-pass skeleton. Compiles and vets clean. Several behaviors depend on a live
instance to confirm — search both repos for `TODO(verify-live)`:
- async task id/result field mapping,
- snapshot id resolution after create,
- echo-db read-back (no get-by-id in the spec),
- app-token list envelope / delete path.

The underlying API spec is `version 0.0.1` with an empty `servers` block — treat
the surface as not-yet-stable before publishing.
