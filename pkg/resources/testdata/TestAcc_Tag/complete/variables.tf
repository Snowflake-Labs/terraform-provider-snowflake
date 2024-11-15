variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "name" {
  type = string
}

variable "allowed_values" {
  type = set(string)
}

variable "comment" {
  default = null
  type    = string
}

variable "masking_policies" {
  default = null
  type    = set(string)
}
