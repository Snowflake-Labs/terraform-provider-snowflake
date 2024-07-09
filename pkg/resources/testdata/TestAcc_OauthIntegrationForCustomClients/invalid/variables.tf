variable "name" {
  type = string
}
variable "oauth_client_type" {
  type = string
}
variable "oauth_redirect_uri" {
  type = string
}
variable "oauth_use_secondary_roles" {
  type = string
}
variable "blocked_roles_list" {
  type = set(string)
}
