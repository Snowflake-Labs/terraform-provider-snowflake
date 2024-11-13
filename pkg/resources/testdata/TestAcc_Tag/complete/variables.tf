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
  type = list(string)
}

variable "comment" {
  default = null
  type    = string
}

variable "masking_policies" {
  default = null
  type    = list(string)
}
