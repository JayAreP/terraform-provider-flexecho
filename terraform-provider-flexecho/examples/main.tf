terraform {
  required_providers {
    flexecho = {
      # local build. run build-local.ps1 -Install (or .sh --install) first.
      # swap to "silk-us/flexecho" once published to the registry
      source  = "localdomain/provider/flexecho"
      version = "0.1.0"
    }
  }
}

provider "flexecho" {
  server = "20.84.120.30"      # or set SILK_FLEX_SERVER
  token  = var.flex_token      # or set SILK_FLEX_TOKEN
}

variable "flex_token" {
  type      = string
  sensitive = true
}

# register an echo host
resource "flexecho_host" "sql01" {
  host_id   = "sql01"
  db_vendor = "mssql"
}

# capture a db snapshot (async, PRovider polls the task to completion)
resource "flexecho_db_snapshot" "nightly" {
  source_host_id = flexecho_host.sql01.host_id
  database_ids   = ["db-orders", "db-inventory"]
  name_prefix    = "nightly"
  consistency_level = "application"
}

# clone an echo db FROM that snapshot onto another host
resource "flexecho_echo_db" "orders_clone" {
  snapshot_id         = flexecho_db_snapshot.nightly.id
  destination_host_id = "sql02"
  destination_db_id   = "db-orders"
  destination_db_name = "orders_clone"
  target_state        = "online"
}

# an app token (core flex api)
resource "flexecho_app_token" "ci" {
  app_name    = "ci-pipeline"
  description = "token for CI snapshot jobs"
}

# data sources
data "flexecho_hosts" "all" {}

data "flexecho_db_snapshots" "for_host" {
  host_id = flexecho_host.sql01.host_id
}

output "all_hosts" {
  value = data.flexecho_hosts.all.hosts
}
