resource "snowflake_row_access_policy_grant" "grant" {
  database_name          = "db"
  schema_name            = "schema"
  row_access_policy_name = "row_access_policy"

  privilege = "APPLY"
  roles = [
    "role1",
    "role2",
  ]

  with_grant_option = false
}
