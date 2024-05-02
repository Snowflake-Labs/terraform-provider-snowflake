variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "name" {
  type = string
}

variable "comment" {
  type = string
}

variable "allowed_values" {
  type = list(string)
}