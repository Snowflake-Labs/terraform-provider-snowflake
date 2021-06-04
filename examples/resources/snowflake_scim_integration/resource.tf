resource "snowflake_scim_integration" "aad" {
  name             = "AAD_PROVISIONING"
  network_policy   = "AAD_NETWORK_POLICY"
  provisioner_role = "AAD_PROVISIONER"
  scim_client      = "AZURE"
}