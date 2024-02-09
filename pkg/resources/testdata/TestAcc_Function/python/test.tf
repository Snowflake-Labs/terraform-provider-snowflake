resource "snowflake_function" "f" {
  database = var.database
  schema   = var.schema
  name     = var.name
  arguments {
    name = "x"
    type = "NUMBER"
  }
  language            = "python"
  return_type         = "VARIANT"
  return_behavior     = "VOLATILE"
  null_input_behavior = "CALLED ON NULL INPUT"
  runtime_version     = "3.8"
  handler             = "dump"
  comment             = var.comment
  statement           = <<EOT
def dump(i):
	print("Hello World!")
  EOT
}
