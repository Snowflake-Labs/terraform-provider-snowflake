# format is database_name | schema_name | object_name | argument_data_types | privilege | with_grant_option | roles | shares
terraform import snowflake_procedure_grant.example 'MY_DATABASE|MY_SCHEMA|MY_OBJECT_NAME|ARG1TYPE,ARG2TYPE|USAGE|false|role1,role2|share1,share2'
