## Minimal
resource "snowflake_account" "minimal" {
  name                 = "ACCOUNT_NAME"
  admin_name           = var.admin_name
  admin_password       = var.admin_password
  email                = var.email
  edition              = "STANDARD"
  grace_period_in_days = 3
}

## Complete (with SERVICE user type)
resource "snowflake_account" "complete" {
  name                 = "ACCOUNT_NAME"
  admin_name           = var.admin_name
  admin_rsa_public_key = "<public_key>"
  admin_user_type      = "SERVICE"
  email                = var.email
  edition              = "STANDARD"
  region_group         = "PUBLIC"
  region               = "AWS_US_WEST_2"
  comment              = "some comment"
  is_org_admin         = "true"
  grace_period_in_days = 3
}

## Complete (with PERSON user type)
resource "snowflake_account" "complete" {
  name                 = "ACCOUNT_NAME"
  admin_name           = var.admin_name
  admin_password       = var.admin_password
  admin_user_type      = "PERSON"
  first_name           = var.first_name
  last_name            = var.last_name
  email                = var.email
  must_change_password = "false"
  edition              = "STANDARD"
  region_group         = "PUBLIC"
  region               = "AWS_US_WEST_2"
  comment              = "some comment"
  is_org_admin         = "true"
  grace_period_in_days = 3
}

variable "admin_name" {
  type      = string
  sensitive = true
}

variable "email" {
  type      = string
  sensitive = true
}

variable "admin_password" {
  type      = string
  sensitive = true
}

variable "first_name" {
  type      = string
  sensitive = true
}

variable "last_name" {
  type      = string
  sensitive = true
}
