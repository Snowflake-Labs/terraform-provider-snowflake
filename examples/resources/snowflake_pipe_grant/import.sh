# format is database_name | schema_name | object_name | privilege | with_grant_option | roles
terraform import snowflake_pipe_grant.example 'MY_DATABASE|MY_SCHEMA|MY_OBJECT_NAME|OPERATE|false|role1,role2'
