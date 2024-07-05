
variable "allowed_email_patterns" {
  type = list(string)
}
variable "allowed_user_domains" {
  type = list(string)
}
variable "comment" {
  type = string
}
variable "enabled" {
  type = bool
}
variable "name" {
  type = string
}
variable "saml2_enable_sp_initiated" {
  type = bool
}
variable "saml2_force_authn" {
  type = bool
}
variable "saml2_issuer" {
  type = string
}
variable "saml2_post_logout_redirect_url" {
  type = string
}
variable "saml2_provider" {
  type = string
}
variable "saml2_requested_nameid_format" {
  type = string
}
variable "saml2_sign_request" {
  type = bool
}
variable "saml2_snowflake_acs_url" {
  type = string
}
variable "saml2_snowflake_issuer_url" {
  type = string
}
variable "saml2_sp_initiated_login_page_label" {
  type = string
}
variable "saml2_sso_url" {
  type = string
}
variable "saml2_x509_cert" {
  type = string
}
