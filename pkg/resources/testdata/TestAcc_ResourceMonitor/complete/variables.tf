variable "name" {
  type = string
}

variable "notify_users" {
  type = set(string)
}

variable "credit_quota" {
  type = number
}

variable "frequency" {
  type = string
}

variable "start_timestamp" {
  type = string
}

variable "end_timestamp" {
  type = string
}

variable "trigger" {
  type = set(object({
    threshold = number
    on_threshold_reached = string
  }))
}
