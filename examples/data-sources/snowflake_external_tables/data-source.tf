data "snowflake_external_tables" "current" {
    database = "MYDB"
    schema   = "MYSCHEMA"
}