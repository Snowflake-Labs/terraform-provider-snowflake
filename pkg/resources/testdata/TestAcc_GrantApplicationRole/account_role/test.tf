resource "snowflake_role" "role" {
  name = var.parent_role_name
}

resource "snowflake_grant_application_role" "g" {
  name             = var.name
  parent_role_name = snowflake_role.role.name
}
