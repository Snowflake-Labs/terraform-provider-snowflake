data "snowflake_row_access_policies" "current" {
    database = "MYDB"
    schema   = "MYSCHEMA"
}