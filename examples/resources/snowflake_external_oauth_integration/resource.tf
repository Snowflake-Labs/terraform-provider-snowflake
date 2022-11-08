resource "snowflake_external_oauth_integration" "azure" {
  name                             = "AZURE_POWERBI"
  type                             = "AZURE"
  enabled                          = true
  issuer                           = "https://sts.windows.net/00000000-0000-0000-0000-000000000000"
  snowflake_user_mapping_attribute = "LOGIN_NAME"
  jws_keys_urls                    = ["https://login.windows.net/common/discovery/keys"]
  audience_urls                    = ["https://analysis.windows.net/powerbi/connector/Snowflake"]
  token_user_mapping_claims        = ["upn"]
}