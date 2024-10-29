## Minimal
resource "snowflake_connection" "basic" {
  name = "connection_name"
}

## Enable failover to accounts
resource "snowflake_connection" "with_enable_failover_list" {
  name                        = "connection_name"
  enable_failover_to_accounts = [
    "<secondary_account_organization_name>.<secondary_account_name>"
  ]
}

## As replica
resource "snowflake_connection" "replica" {
  name          = "connection_name"
  as_replica_of = "<organization_name>.<account_name>.<connection_name>"
}

# As replica with promotion to primary
resource "snowflake_connection" "replica_with_promotion" {
  name          = "connection_name"
  as_replica_of = "<organization_name>.<account_name>.<connection_name>"
  is_primary    = true
}

## Complete (with every optional set)
resource "snowflake_connection" "complete" {
  name          = "connection_name"
  as_replica_of = "<organization_name>.<account_name>.<connection_name>"
  is_primary    = true
  enable_failover_to_accounts = [
    "<secondary_account_organization_name>.<secondary_account_name>"
  ]
  comment = "my complete connection"
}
