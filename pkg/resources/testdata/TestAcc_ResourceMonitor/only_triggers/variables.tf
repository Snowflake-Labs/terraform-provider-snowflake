variable "name" {
  type = string
}

variable "trigger" {
  type = set(object({
    threshold = number
    on_threshold_reached = string
  }))
}
