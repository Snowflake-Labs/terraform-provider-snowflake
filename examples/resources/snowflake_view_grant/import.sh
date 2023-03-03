# format is database_name ❄️ schema_name ❄️ view_name ❄️ privilege ❄️ with_grant_option ❄️ roles ❄️ shares
terraform import snowflake_view_grant.example 'MY_DATABASE❄️MY_SCHEMA❄️MY_OBJECT❄️USAGE❄️false❄️role1,role2❄️share1,share2'
