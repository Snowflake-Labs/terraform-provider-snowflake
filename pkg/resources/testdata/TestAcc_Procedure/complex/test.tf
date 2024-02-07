resource "snowflake_procedure" "p" {
  database            = var.database
  schema              = var.schema
  name                = var.name
  arguments {
		name = "arg1"
		type = "VARCHAR"
	}
	arguments {
		name = "arg2"
		type = "DATE"
	}
  language            = "JAVASCRIPT"
  return_type         = "VARCHAR"
  execute_as          = "CALLER"
  null_input_behavior = "RETURNS NULL ON NULL INPUT"
  comment             = var.comment
  statement           = <<EOT
var x = 1
return x
  EOT
}
