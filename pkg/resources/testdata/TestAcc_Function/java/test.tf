resource "snowflake_function" "f" {
  database = var.database
  schema   = var.schema
  name     = var.name
  arguments {
    name = "x"
    type = "VARCHAR"
  }
  language            = "java"
  return_type         = "VARCHAR"
  return_behavior     = "VOLATILE"
  null_input_behavior = "CALLED ON NULL INPUT"
  handler             = "TestFunc.echoVarchar"
  comment             = var.comment
  statement           = <<EOT
		class TestFunc {
			public static String echoVarchar(String x) {
				return x;
			}
		}
  EOT
}
