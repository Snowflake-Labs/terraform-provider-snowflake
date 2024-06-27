
variable "blocked_roles_list" {
  type = set(string)
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
variable "oauth_client" {
  type = string
}
variable "oauth_issue_refresh_tokens" {
  type = string
}
variable "oauth_refresh_token_validity" {
  type = number
}
variable "oauth_use_secondary_roles" {
  type = string
}
