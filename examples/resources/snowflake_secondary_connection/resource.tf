## Minimal
resource "snowflake_secondary_connection" "basic" {
  name          = "connection_name"
  as_replica_of = "<organization_name>.<account_name>.<connection_name>"
}

## Complete (with every optional set)
resource "snowflake_secondary_connection" "complete" {
  name          = "connection_name"
  as_replica_of = "<organization_name>.<account_name>.<connection_name>"
  comment       = "my complete secondary connection"
}
