# format is database_name | schema_name | privilege | with_grant_option | roles | shares
terraform import snowflake_schema_grant.example 'MY_DATABASE|MY_SCHEMA|MONITOR|false|role1,role2|share1,share2'
