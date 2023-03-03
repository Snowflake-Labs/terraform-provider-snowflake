# format is database_name | schema_name | tag_name | privilege | with_grant_option | roles
terraform import snowflake_tag_grant.example 'MY_DATABASE|MY_SCHEMA|MY_OBJECT|APPLY|false|role1,role2'
