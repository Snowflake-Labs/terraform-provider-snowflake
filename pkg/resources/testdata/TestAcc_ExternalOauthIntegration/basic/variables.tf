variable "name" {
  type = string
}
variable "type" {
  type = string
}
variable "enabled" {
  type = bool
}
variable "external_oauth_issuer" {
  type = string
}
variable "external_oauth_snowflake_user_mapping_attribute" {
  type = string
}
variable "external_oauth_token_user_mapping_claim" {
  type = set(any)
}
