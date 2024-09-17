resource "snowflake_masking_policy" "test_1" {
  name     = var.name_1
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.arguments
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body             = var.body
  return_data_type = var.return_data_type
}

resource "snowflake_masking_policy" "test_2" {
  name     = var.name_2
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.arguments
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body             = var.body
  return_data_type = var.return_data_type
}

resource "snowflake_masking_policy" "test_3" {
  name     = var.name_3
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.arguments
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body             = var.body
  return_data_type = var.return_data_type
}


data "snowflake_masking_policies" "test" {
  depends_on = [snowflake_masking_policy.test_1, snowflake_masking_policy.test_2, snowflake_masking_policy.test_3]

  like = var.like
}
