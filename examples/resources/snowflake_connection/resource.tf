## Minimal
resource "snowflake_connection" "basic" {
  name = "connection_name"
}

## Complete (with every optional set)
resource "snowflake_connection" "complete" {
  name    = "connection_name"
  comment = "my complete connection"
  enable_failover_to_accounts = [
    "<secondary_account_organization_name>.<secondary_account_name>"
  ]
}
