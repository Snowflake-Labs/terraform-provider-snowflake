data "snowflake_grants" "test" {
  grants_on {
    object_name = var.fully_qualified_function_name
    object_type = "FUNCTION"
  }
}
