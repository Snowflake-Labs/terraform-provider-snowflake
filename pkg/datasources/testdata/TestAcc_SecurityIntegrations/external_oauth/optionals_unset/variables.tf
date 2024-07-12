variable "comment" {
  type = string
}
variable "enabled" {
  type = bool
}
variable "external_oauth_allowed_roles_list" {
  type = set(string)
}
variable "external_oauth_any_role_mode" {
  type = string
}
variable "external_oauth_audience_list" {
  type = set(string)
}
variable "external_oauth_issuer" {
  type = string
}
variable "external_oauth_jws_keys_url" {
  type = set(string)
}
variable "external_oauth_scope_delimiter" {
  type = string
}
variable "external_oauth_scope_mapping_attribute" {
  type = string
}
variable "external_oauth_snowflake_user_mapping_attribute" {
  type = string
}
variable "external_oauth_token_user_mapping_claim" {
  type = set(string)
}
variable "name" {
  type = string
}
variable "external_oauth_type" {
  type = string
}
