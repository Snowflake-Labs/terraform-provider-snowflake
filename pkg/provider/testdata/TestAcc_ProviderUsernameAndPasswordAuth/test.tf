variable "account" {
  type = string
}
variable "user" {
  type = string
}

variable "password" {
  type = string
}

variable "role" {
  type = string
}

provider "snowflake" {
  account  = var.account
  user     = var.user
  password = var.password
  role     = var.role
}

data "snowflake_current_account" "p" {}
