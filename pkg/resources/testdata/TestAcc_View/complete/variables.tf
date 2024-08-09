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

variable "change_tracking" {
  type = string
}

variable "or_replace" {
  type = bool
}

variable "copy_grants" {
  type = bool
}

variable "row_access_policy" {
  type = string
}

variable "is_temporary" {
  type = string
}

variable "row_access_policy_on" {
  type = list(string)
}

variable "aggregation_policy" {
  type = string
}

variable "aggregation_policy_entity_key" {
  type = list(string)
}

variable "statement" {
  type = string
}
