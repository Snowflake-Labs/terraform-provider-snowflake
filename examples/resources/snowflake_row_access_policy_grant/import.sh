# format is database_name | schema_name | object_name | privilege | with_grant_option | roles
terraform import snowflake_row_access_policy_grant.example 'MY_DATABASE|MY_SCHEMA|MY_OBJECT_NAME|SELECT|false|role1,role2'
