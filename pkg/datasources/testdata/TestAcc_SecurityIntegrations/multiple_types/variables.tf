# saml2
variable "name_1" {
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

# scim
variable "name_2" {
  type = string
}
variable "scim_client" {
  type = string
}
variable "run_as_role" {
  type = string
}
variable "enabled" {
  type = bool
}

variable "like" {
  type = string
}
