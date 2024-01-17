locals {
  snowflake_database_sandbox = "snow_test_db"
  snowflake_schema_list = [
    "snow_test_sch_1",
    "snow_test_sch_2",
  ]
}

resource "snowflake_database" "db" {
  name = local.snowflake_database_sandbox
}

resource "snowflake_schema" "sandbox_raw" {
  depends_on = [snowflake_database.db]
  for_each = toset(local.snowflake_schema_list)

  database     = local.snowflake_database_sandbox
  name         = each.value
  is_transient = false
  is_managed   = true
}

resource "snowflake_grant_privileges_to_role" "schema_usage_sandbox_dev" {
  depends_on = [snowflake_schema.sandbox_raw]
  for_each = snowflake_schema.sandbox_raw

  privileges = ["CREATE TABLE"]
  role_name  = "\"custom.role-123\""
  on_schema {
    schema_name = "\"${each.value.database}\".\"${each.value.name}\""
  }
}