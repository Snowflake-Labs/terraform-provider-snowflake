# basic resource
resource "snowflake_scim_integration" "test" {
  name          = "test"
  enabled       = true
  scim_client   = "GENERIC"
  sync_password = true
  run_as_role   = "GENERIC_SCIM_PROVISIONER"
}

# resource with all fields set
resource "snowflake_scim_integration" "test" {
  name           = "test"
  enabled        = true
  scim_client    = "GENERIC"
  sync_password  = true
  network_policy = "network_policy_test"
  run_as_role    = "GENERIC_SCIM_PROVISIONER"
  comment        = "foo"
}
