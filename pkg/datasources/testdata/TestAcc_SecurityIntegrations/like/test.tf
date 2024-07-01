resource "snowflake_scim_integration" "test_1" {
  name        = var.name_1
  enabled     = false
  scim_client = "GENERIC"
  run_as_role = "GENERIC_SCIM_PROVISIONER"
}

resource "snowflake_scim_integration" "test_2" {
  name        = var.name_2
  enabled     = false
  scim_client = "GENERIC"
  run_as_role = "GENERIC_SCIM_PROVISIONER"
}

resource "snowflake_scim_integration" "test_3" {
  name        = var.name_3
  enabled     = false
  scim_client = "GENERIC"
  run_as_role = "GENERIC_SCIM_PROVISIONER"
}

data "snowflake_security_integrations" "test" {
  depends_on = [snowflake_scim_integration.test_1, snowflake_scim_integration.test_2, snowflake_scim_integration.test_3]

  like = var.like
}
