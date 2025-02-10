# basic resource
resource "snowflake_secret_with_generic_string" "test" {
  name          = "EXAMPLE_SECRET"
  database      = "EXAMPLE_DB"
  schema        = "EXAMPLE_SCHEMA"
  secret_string = var.secret_string
}

# resource with all fields set
resource "snowflake_secret_with_generic_string" "test" {
  name          = "EXAMPLE_SECRET"
  database      = "EXAMPLE_DB"
  schema        = "EXAMPLE_SCHEMA"
  secret_string = var.secret_string
  comment       = "EXAMPLE_COMMENT"
}

variable "secret_string" {
  type      = string
  sensitive = true
}
