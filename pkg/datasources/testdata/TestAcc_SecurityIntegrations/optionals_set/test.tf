resource "snowflake_scim_integration" "test" {
  name           = var.name
  enabled        = false
  scim_client    = "GENERIC"
  run_as_role    = "GENERIC_SCIM_PROVISIONER"
  network_policy = var.network_policy
  comment        = var.comment
}

data "snowflake_security_integrations" "test" {
  depends_on = [snowflake_scim_integration.test]

  like = var.name
}
