variable "name" {
  type = string
}

variable "privileges" {
  type = list(string)
}

variable "with_grant_option" {
  type = bool
}
