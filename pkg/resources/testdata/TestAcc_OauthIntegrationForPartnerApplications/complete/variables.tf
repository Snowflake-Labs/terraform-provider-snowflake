variable "name" {
  type = string
}
variable "oauth_client" {
  type = string
}
variable "oauth_redirect_uri" {
  type = string
}
variable "blocked_roles_list" {
  type = set(string)
}
variable "enabled" {
  type = string
}
variable "oauth_issue_refresh_tokens" {
  type = string
}
variable "oauth_refresh_token_validity" {
  type = string
}
variable "oauth_use_secondary_roles" {
  type = string
}
variable "comment" {
  type = string
}
