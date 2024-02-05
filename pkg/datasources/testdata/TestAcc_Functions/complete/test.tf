variable "function_name_one" {
  type = string
}

variable "function_name_two" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

resource "snowflake_function" "test_function_one" {
  name        = var.function_name_one
  database    = var.database
  schema      = var.schema
  return_type = "VARCHAR"
  language    = "JAVASCRIPT"
  statement   = <<-EOF
		return "Hi"
	EOF
}

resource "snowflake_function" "test_function_two" {
  name     = var.function_name_two
  database = var.database
  schema   = var.schema
  arguments {
    name = "arg1"
    type = "varchar"
  }
  comment     = "Terraform acceptance test"
  return_type = "varchar"
  language    = "JAVASCRIPT"
  statement   = <<-EOF
		var X=1
		return X
  EOF
}

data "snowflake_functions" "functions" {
  database   = var.database
  schema     = var.schema
  depends_on = [snowflake_function.test_function_one, snowflake_function.test_function_two]
}
