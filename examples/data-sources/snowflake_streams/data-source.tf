data "snowflake_streams" "current" {
    database = "MYDB"
    schema   = "MYSCHEMA"
}