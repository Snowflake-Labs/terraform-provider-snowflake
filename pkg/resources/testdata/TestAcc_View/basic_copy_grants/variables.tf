variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "statement" {
  type = string
}

variable "copy_grants" {
  type = bool
}

variable "is_secure" {
  type = bool
}

variable "columns" {
  type = set(map(string))
}
