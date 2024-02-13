variable "proc_name_one" {
  type = string
}

variable "proc_name_two" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

resource "snowflake_procedure" "test_proc_one" {
  name        = var.proc_name_one
  database    = var.database
  schema      = var.schema
  return_type = "VARCHAR"
  language    = "JAVASCRIPT"
  statement   = <<-EOF
		return "Hi"
	EOF
}

resource "snowflake_procedure" "test_proc_two" {
  name     = var.proc_name_two
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

data "snowflake_procedures" "procedures" {
  database   = var.database
  schema     = var.schema
  depends_on = [snowflake_procedure.test_proc_one, snowflake_procedure.test_proc_two]
}
