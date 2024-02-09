resource "snowflake_function" "f" {
  database = var.database
  schema   = var.schema
  name     = var.name
  arguments {
    name = "x"
    type = "FLOAT"
  }
  language            = "SQL"
  return_type         = "FLOAT"
  return_behavior     = "VOLATILE"
  null_input_behavior = "CALLED ON NULL INPUT"
  comment             = var.comment
  statement           = <<EOT
		3.141592654::FLOAT
  EOT
}
