variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "table" {
  type = string
}

variable "copy_grants" {
  type = bool
}

variable "before" {
  type = map(string)
}

variable "comment" {
  type = string
}
