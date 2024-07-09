variable "name" {
  type = string
}
variable "saml2_issuer" {
  type = string
}
variable "saml2_provider" {
  type = string
}
variable "saml2_sso_url" {
  type = string
}
variable "saml2_x509_cert" {
  type = string
}
variable "saml2_sp_initiated_login_page_label" {
  type = string
}
variable "saml2_snowflake_issuer_url" {
  type = string
}
variable "saml2_snowflake_acs_url" {
  type = string
}
variable "allowed_email_patterns" {
  type = list(string)
}
variable "allowed_user_domains" {
  type = list(string)
}
