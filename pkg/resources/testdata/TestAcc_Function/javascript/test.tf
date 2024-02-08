resource "snowflake_function" "f" {
  database            = var.database
  schema              = var.schema
  name                = var.name
  arguments {
    name = "d"
    type = "FLOAT"
  }
  language            = "javascript"
  return_type         = "FLOAT"
  return_behavior     = "VOLATILE"
  null_input_behavior = "CALLED ON NULL INPUT"
  comment             = var.comment
  statement           = <<EOT
		if (D <= 0) {
			return 1;
		} else {
			var result = 1;
			for (var i = 2; i <= D; i++) {
				result = result * i;
			}
			return result;
		}
  EOT
}
