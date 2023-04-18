# format is database_name|schema_name|materialized_view_name|privilege|with_grant_option|on_future|on_all|roles|shares
terraform import snowflake_materialized_view_grant.example "MY_DATABASE|MY_SCHEMA|MY_MV_NAME|SELECT|false|false|role1,role2|share1,share2"

