# basic resource
resource "snowflake_external_oauth_integration" "test" {
  enabled                                         = true
  external_oauth_issuer                           = "issuer"
  external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
  external_oauth_token_user_mapping_claims        = ["foo"]
  name                                            = "test"
  external_oauth_type                             = "CUSTOM"
}
# resource with all fields set (jws keys url flow)
resource "snowflake_external_oauth_integration" "test" {
  comment                                         = "foo"
  enabled                                         = true
  external_oauth_allowed_roles_list               = ["foo"]
  external_oauth_any_role_mode                    = "ENABLED"
  external_oauth_audience_list                    = ["foo"]
  external_oauth_blocked_roles_list               = ["bar"]
  external_oauth_issuer                           = "issuer"
  external_oauth_jws_keys_url                     = ["https://example.com"]
  external_oauth_scope_delimiter                  = ","
  external_oauth_scope_mapping_attribute          = "LOGIN_NAME"
  external_oauth_snowflake_user_mapping_attribute = "foo"
  external_oauth_token_user_mapping_claims        = ["foo"]
  name                                            = "foo"
  external_oauth_type                             = "CUSTOM"
}
# resource with all fields set (rsa public key flow)
resource "snowflake_external_oauth_integration" "test" {
  comment                                         = "foo"
  enabled                                         = true
  external_oauth_allowed_roles_list               = ["foo"]
  external_oauth_any_role_mode                    = "ENABLED"
  external_oauth_audience_list                    = ["foo"]
  external_oauth_blocked_roles_list               = ["bar"]
  external_oauth_issuer                           = "issuer"
  external_oauth_rsa_public_key                   = file("key.pem")
  external_oauth_rsa_public_key_2                 = file("key2.pem")
  external_oauth_scope_delimiter                  = "foo"
  external_oauth_scope_mapping_attribute          = "LOGIN_NAME"
  external_oauth_snowflake_user_mapping_attribute = "foo"
  external_oauth_token_user_mapping_claims        = ["foo"]
  name                                            = "foo"
  external_oauth_type                             = "CUSTOM"
}
