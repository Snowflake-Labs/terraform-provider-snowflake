variable "name" {
  type = string
}

variable "transient" {
  type = bool
}

variable "data_retention_time_in_days" {
  type = string
}

variable "max_data_extension_time_in_days" {
  type = string
}

variable "external_volume" {
  type = string
}

variable "catalog" {
  type = string
}

variable "replace_invalid_characters" {
  type = string
}

variable "default_ddl_collation" {
  type = string
}

variable "storage_serialization_policy" {
  type = string
}

variable "log_level" {
  type = string
}

variable "trace_level" {
  type = string
}

variable "comment" {
  type = string
}
