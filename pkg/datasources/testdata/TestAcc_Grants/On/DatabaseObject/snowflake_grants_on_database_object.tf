data "snowflake_grants" "test" {
  grants_on {
    object_name = "\"${var.database}\".\"${var.schema}\""
    object_type = "SCHEMA"
  }
}
