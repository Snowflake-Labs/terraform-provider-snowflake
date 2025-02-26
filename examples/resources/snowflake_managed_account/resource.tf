resource "snowflake_managed_account" "account" {
  name           = "managed account"
  admin_name     = "admin"
  admin_password = var.admin_password
  type           = "READER"
  comment        = "A managed account."
  cloud          = "aws"
  region         = "us-west-2"
  locator        = "managed-account"
}

variable "admin_password" {
  type      = string
  sensitive = true
}
