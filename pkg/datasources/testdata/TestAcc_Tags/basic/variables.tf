variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "comment" {
  type = string
}

variable "allowed_values" {
  type = set(string)
}
