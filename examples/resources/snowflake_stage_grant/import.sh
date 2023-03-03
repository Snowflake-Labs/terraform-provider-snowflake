# format is database_name | schema_name | stage_name | privilege | with_grant_option | roles
terraform import snowflake_stage_grant.example 'MY_DATABASE|MY_SCHEMA|MY_OBJECT|USAGE|false|role1,role2'
