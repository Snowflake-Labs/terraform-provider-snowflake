data "snowflake_current_role" "test" {}

data "snowflake_grants" "test" {
  grants_of {
    role = data.snowflake_current_role.test.name
  }
}
