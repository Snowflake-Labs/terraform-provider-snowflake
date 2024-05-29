variable "name" {
  type = string
}

variable "as_replica_of" {
  type = string
}

variable "transient" {
  type = bool
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
