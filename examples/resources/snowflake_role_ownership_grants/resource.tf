resource "snowflake_role" "role" {
  name    = "rking_test_role"
  comment = "for testing"
}

resource "snowflake_role" "other_role" {
  name = "rking_test_role2"
}

# ensure the Terraform user inherits ownership privileges for the rking_test_role role
# otherwise Terraform will fail to destroy the rking_test_role2 role due to insufficient privileges
resource "snowflake_role_grants" "grants" {
  role_name = snowflake_role.role.name

  roles = [
    "ACCOUNTADMIN",
  ]
}

resource "snowflake_role_ownership_grant" "grant" {
  on_role_name   = snowflake_role.role.name
  to_role_name   = snowflake_role.other_role.name
  current_grants = "COPY"
}