variable "name" {
  type = string
}
variable "oauth_client" {
  type = string
}
variable "blocked_roles_list" {
  type = set(string)
}
