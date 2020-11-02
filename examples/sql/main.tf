/*
 * simple example
 *
 * You would NEVER use this module to manage a warehouse. Instead use snowflake_warehouse...
 *
 * This resource should only be used when existing resources do not exist or
 * do not support the required functionality.
 */
# resource "snowflake_sql_script" "script" {
#   name = "testing"
#   lifecycle_commands {
#     create = "CREATE OR REPLACE WAREHOUSE TESTING;"
#     #read   = "SHOW WAREHOUSES LIKE TESTING;"
#     delete = "DROP WAREHOUSE TESTING;"
#   }
# }


/*
 * grant all on all tables example
 */
resource "snowflake_database" "database" {
  name = "TESTING_DATABASE"
}

resource "snowflake_schema" "schema" {
  name = "TESTING_SCHEMA"
  database = snowflake_database.database.name
}

resource "snowflake_table" "table" {
  name = "TESTING_TABLE"
  database = snowflake_database.database.name
  schema = snowflake_schema.schema.name
  
  column {
    name = "col1"
    type = "VARIANT"
  }
}

resource "snowflake_role" "role" {
  name = "TESTING_ROLE"
}

resource "snowflake_sql_script" "script" {
  depends_on = [
    snowflake_table.table,
  ]
  name = "grant-all-on-all-tables-on-database-to-role"
  lifecycle_commands {
    create = join("", ["GRANT ALL ON ALL TABLES IN DATABASE ", snowflake_database.database.name, " TO ROLE ", snowflake_role.role.name, ";"])
    read   = join("", ["SHOW GRANTS TO ROLE ", snowflake_role.role.name, ";"])
    delete = join("", ["REVOKE ALL ON ALL TABLES IN DATABASE ", snowflake_database.database.name, " FROM ROLE ", snowflake_role.role.name, ";"])
  }
}

/*
 * sad path
 */
# resource "snowflake_sql_script" "script" {
#   name = "testing"
#   lifecycle_commands {
#     create = "bad query"
#     #read   = "SHOW WAREHOUSES LIKE TESTING;"
#     delete = "bad query"
#   }
# }


