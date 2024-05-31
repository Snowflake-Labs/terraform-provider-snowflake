resource "snowflake_scim_integration" "test" {
  name        = var.name
  enabled     = var.enabled
  scim_client = var.scim_client
  run_as_role = var.run_as_role
}
