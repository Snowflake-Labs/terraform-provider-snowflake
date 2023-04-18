# format is database|schema|external_table|privilege|with_grant_option|on_future|roles|shares
terraform import snowflake_external_table_grant.example "MY_DATABASE|MY_SCHEMA|MY_TABLE_NAME|SELECT|false|false|role1,role2|share1,share2"
