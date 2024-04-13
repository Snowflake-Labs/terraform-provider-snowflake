data "snowflake_grants" "test" {
  grants_to {
    user = var.user
  }
}
