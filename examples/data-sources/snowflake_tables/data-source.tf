data "snowflake_tables" "current" {
    database = "MYDB"
    schema   = "MYSCHEMA"
}