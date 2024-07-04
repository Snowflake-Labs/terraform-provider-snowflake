variable "name" {
  type = string
}
variable "oauth_client" {
  type = string
}
variable "oauth_use_secondary_roles" {
  type = string
}
variable "blocked_roles_list" {
  type = set(string)
}
