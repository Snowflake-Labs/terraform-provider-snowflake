resource "snowflake_api_integration" "test_api_int" {
  name                 = var.name
  api_provider         = "aws_api_gateway"
  api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
  api_allowed_prefixes = var.api_allowed_prefixes
  enabled              = true
}

resource "snowflake_function" "test_func_req_translator" {
  name     = "${var.name}_request_translator"
  database = var.database
  schema   = var.schema
  arguments {
    name = "EVENT"
    type = "OBJECT"
  }
  comment     = "Terraform acceptance test"
  return_type = "OBJECT"
  language    = "javascript"
  statement   = <<EOH
		  	let exeprimentName = EVENT.body.data[0][1]
		  	return { "body": { "name": test }}
	  	EOH
}


resource "snowflake_function" "test_func_res_translator" {
  name     = "${var.name}_response_translator"
  database = var.database
  schema   = var.schema
  arguments {
    name = "EVENT"
    type = "OBJECT"
  }
  comment     = "Terraform acceptance test"
  return_type = "OBJECT"
  language    = "javascript"
  statement   = <<EOH
			  return { "body": { "data" :  [[0, EVENT]] } };
		  EOH
}


resource "snowflake_external_function" "external_function" {
  name            = var.name
  database        = var.database
  schema          = var.schema
  comment         = var.comment
  return_type     = "VARIANT"
  return_behavior = "IMMUTABLE"
  api_integration = snowflake_api_integration.test_api_int.name
  header {
    name  = "x-custom-header"
    value = "snowflake"
  }
  max_batch_rows            = 500
  request_translator        = snowflake_function.test_func_req_translator.fully_qualified_name
  response_translator       = snowflake_function.test_func_res_translator.fully_qualified_name
  url_of_proxy_and_resource = var.url_of_proxy_and_resource
}
