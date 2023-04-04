resource "snowflake_masking_policy_grant" "example" {
  masking_policy_name    = "EXAMPLE_MASKING_POLICY_NAME"
  database_name          = "EXAMPLE_DB_NAME"
  schema_name            = "EXAMPLE_SCHEMA_NAME"
  privilege              = "APPLY"
  roles                  = ["ROLE1_NAME", "ROLE2_NAME"]
  with_grant_option      = true
  enable_multiple_grants = true
}
