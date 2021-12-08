resource "snowflake_role" "role" {
  name    = "rking_test_role"
  comment = "for testing"
}

resource "snowflake_user" "user" {
  name    = "rking_test_user"
  comment = "for testing"
}

resource "snowflake_user" "user2" {
  name    = "rking_test_user2"
  comment = "for testing"
}

resource "snowflake_role" "other_role" {
  name = "rking_test_role2"
}

resource "snowflake_role_grants" "grants" {
  role_name = "${snowflake_role.role.name}"

  roles = [
    "${snowflake_role.other_role.name}",
  ]

  users = [
    "${snowflake_user.user.name}",
    "${snowflake_user.user2.name}",
  ]
}
