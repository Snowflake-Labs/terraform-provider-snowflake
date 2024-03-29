data "snowflake_grants" "test" {
  grants_on {
    object_name = var.database
    object_type = "DATABASE"
  }
}
