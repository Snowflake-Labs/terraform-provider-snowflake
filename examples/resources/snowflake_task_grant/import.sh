# format is database_name ❄️ schema_name ❄️ task_name ❄️ privilege ❄️ with_grant_option ❄️ roles
terraform import snowflake_task_grant.example 'MY_DATABASE❄️MY_SCHEMA❄️MY_OBJECT❄️OPERATE❄️false❄️role1,role2'
