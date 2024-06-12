resource "snowflake_share" "example" {
  name = "test"
}

resource "snowflake_database" "example" {
  # remember to define dependency between objects on a share, because shared objects have to be dropped before dropping share
  depends_on = [snowflake_share.example]
  name       = "test"
}

##################################
### on database
##################################

resource "snowflake_grant_privileges_to_share" "example" {
  to_share    = snowflake_share.example.name
  privileges  = ["USAGE"]
  on_database = snowflake_database.example.name
}

## ID: "\"share_name\"|USAGE|OnDatabase|\"database_name\""

##################################
### on schema
##################################

resource "snowflake_grant_privileges_to_share" "example" {
  to_share   = snowflake_share.example.name
  privileges = ["USAGE"]
  on_schema  = "${snowflake_database.example.name}.${snowflake_schema.example.name}"
}

## ID: "\"share_name\"|USAGE|OnSchema|\"database_name\".\"schema_name\""

##################################
### on table
##################################

resource "snowflake_grant_privileges_to_share" "example" {
  to_share   = snowflake_share.example.name
  privileges = ["SELECT"]
  on_table   = "${snowflake_database.example.name}.${snowflake_schema.example.name}.${snowflake_table.example.name}"
}

## ID: "\"share_name\"|SELECT|OnTable|\"database_name\".\"schema_name\".\"table_name\""

##################################
### on all tables in schema
##################################

resource "snowflake_grant_privileges_to_share" "example" {
  to_share                = snowflake_share.example.name
  privileges              = ["SELECT"]
  on_all_tables_in_schema = "${snowflake_database.example.name}.${snowflake_schema.example.name}"
}

## ID: "\"share_name\"|SELECT|OnAllTablesInSchema|\"database_name\".\"schema_name\""

##################################
### on tag
##################################

resource "snowflake_grant_privileges_to_share" "example" {
  to_share   = snowflake_share.example.name
  privileges = ["READ"]
  on_tag     = "${snowflake_database.example.name}.${snowflake_schema.example.name}.${snowflake_tag.example.name}"
}

## ID: "\"share_name\"|READ|OnTag|\"database_name\".\"schema_name\".\"tag_name\""

##################################
### on view
##################################

resource "snowflake_grant_privileges_to_share" "example" {
  to_share   = snowflake_share.example.name
  privileges = ["SELECT"]
  on_view    = "${snowflake_database.example.name}.${snowflake_schema.example.name}.${snowflake_view.example.name}"
}

## ID: "\"share_name\"|SELECT|OnView|\"database_name\".\"schema_name\".\"view_name\""
