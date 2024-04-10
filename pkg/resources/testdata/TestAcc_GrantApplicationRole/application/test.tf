resource "snowflake_grant_application_role" "g" {
  name             = var.name
  application_name = var.application_name
}
