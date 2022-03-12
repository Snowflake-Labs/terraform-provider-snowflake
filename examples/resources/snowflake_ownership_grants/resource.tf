resource "snowflake_role" "role" {
  name    = "rking_test_role"
  comment = "for testing"
}

resource "snowflake_role" "other_role" {
  name = "rking_test_role2"
}

resource "snowflake_ownership_grants" "grants" {
  roles = [
    "${snowflake_role.other_role.name}",
  ]
  owner = "${snowflake_role.role.name}"
  
  current_grants = "COPY"
}
