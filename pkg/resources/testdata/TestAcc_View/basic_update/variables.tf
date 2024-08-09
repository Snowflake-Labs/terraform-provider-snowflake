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

variable "row_access_policy" {
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

variable "comment" {
  type = string
}
