data "snowflake_materialized_views" "current" {
    database = "MYDB"
    schema   = "MYSCHEMA"
}