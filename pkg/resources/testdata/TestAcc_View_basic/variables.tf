variable "name" {
  type = string
}

variable "comment" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "is_secure" {
  type = bool
}

variable "or_replace" {
  type = bool
}

variable "copy_grants" {
  type = bool
}

variable "statement" {
  type = string
}
