variable "name" {
  type = string
}

variable "as_replica_of" {
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

variable "default_ddl_collation" {
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
