# basic resource
resource "snowflake_secret_with_basic_authentication" "test" {
  name     = "EXAMPLE_SECRET"
  database = "EXAMPLE_DB"
  schema   = "EXAMPLE_SCHEMA"
  username = var.username
  password = var.password
}


# resource with all fields set
resource "snowflake_secret_with_basic_authentication" "test" {
  name     = "EXAMPLE_SECRET"
  database = "EXAMPLE_DB"
  schema   = "EXAMPLE_SCHEMA"
  username = var.username
  password = var.password
  comment  = "EXAMPLE_COMMENT"
}

variable "username" {
  type      = string
  sensitive = true
}

variable "password" {
  type      = string
  sensitive = true
}
