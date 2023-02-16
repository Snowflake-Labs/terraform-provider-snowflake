# format is database_name ❄️ schema_name ❄️ table_name ❄️ privilege ❄️ with_grant_option ❄️ roles ❄️ shares
terraform import snowflake_table_grant.example 'MY_DATABASE❄️MY_SCHEMA❄️MY_OBJECT❄️MODIFY❄️false❄️role1,role2❄️share1,share2'
