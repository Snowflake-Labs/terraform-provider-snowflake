# format is warehouse_name|privilege|with_grant_option|roles
terraform import snowflake_warehouse_grant.example "MY_WAREHOUSE|MODIFY|false|role1,role2"
