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
