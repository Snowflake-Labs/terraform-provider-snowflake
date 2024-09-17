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

variable "data_metric_schedule_using_cron" {
  type = string
}

variable "data_metric_function" {
  type = string
}

variable "data_metric_function_on" {
  type = list(string)
}

variable "schedule_status" {
  type = string
}

variable "column" {
  type = set(map(string))
}
