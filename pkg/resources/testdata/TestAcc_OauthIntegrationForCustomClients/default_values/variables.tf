variable "name" {
  type = string
}
variable "oauth_client_type" {
  type = string
}
variable "oauth_redirect_uri" {
  type = string
}
variable "blocked_roles_list" {
  type = set(string)
}
variable "enabled" {
  type = bool
}
variable "oauth_allow_non_tls_redirect_uri" {
  type = bool
}
variable "oauth_enforce_pkce" {
  type = bool
}
variable "oauth_use_secondary_roles" {
  type = string
}
variable "pre_authorized_roles_list" {
  type = set(string)
}
variable "oauth_issue_refresh_tokens" {
  type = bool
}
variable "oauth_refresh_token_validity" {
  type = number
}
variable "network_policy" {
  type = string
}
variable "oauth_client_rsa_public_key" {
  type = string
}
variable "oauth_client_rsa_public_key_2" {
  type = string
}
variable "comment" {
  type = string
}
