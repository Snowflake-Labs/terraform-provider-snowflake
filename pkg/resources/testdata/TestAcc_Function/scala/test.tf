resource "snowflake_function" "f" {
  database = var.database
  schema   = var.schema
  name     = var.name
  arguments {
    name = "x"
    type = "VARCHAR"
  }
  language            = "scala"
  return_type         = "VARCHAR"
  return_behavior     = "VOLATILE"
  null_input_behavior = "CALLED ON NULL INPUT"
  runtime_version     = "2.12"
  handler             = "Echo.echoVarchar"
  comment             = var.comment
  statement           = <<EOT
		class Echo {
			def echoVarchar(x : String): String = {
				return x
			}
		}
  EOT
}
