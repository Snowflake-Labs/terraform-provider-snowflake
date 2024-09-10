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

variable "is_recursive" {
  type = bool
}

variable "columns" {
  type = set(map(string))
}
