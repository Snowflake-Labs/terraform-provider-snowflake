resource "snowflake_secondary_database" "test" {
  name        = var.name
  as_replica_of  = var.as_replica_of
  comment     = var.comment
}
