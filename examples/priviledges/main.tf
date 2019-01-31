# NOTE that this is an RFC for how things *could* work.

# This is the schema for granting priviledges
# https://docs.snowflake.net/manuals/sql-reference/sql/grant-privilege.html
# 
# GRANT {  { globalPrivileges | ALL [ PRIVILEGES ] } ON ACCOUNT
#        | { accountObjectPrivileges | ALL [ PRIVILEGES ] } ON { RESOURCE MONITOR | WAREHOUSE | DATABASE } <object_name>
#        | { schemaPrivileges | ALL [ PRIVILEGES ] } ON { SCHEMA <schema_name> | ALL SCHEMAS IN DATABASE <db_name> }
#        | { schemaObjectPrivileges | ALL [ PRIVILEGES ] } ON { <object_type> <object_name> | ALL <object_type>S IN SCHEMA <schema_name> }  }
#        | { schemaObjectPrivileges | ALL [ PRIVILEGES ] } ON FUTURE <object_type>S IN SCHEMA <schema_name> }  }
#   TO [ ROLE ] <role_name> [ WITH GRANT OPTION ]

# globalPrivileges ::=
#   { { CREATE { ROLE | USER | WAREHOUSE | DATABASE } } | MANAGE GRANTS | MONITOR USAGE } [ , ... ] }

# accountObjectPrivileges ::=
# -- For RESOURCE MONITOR
#   { MODIFY | MONITOR } [ , ... ]
# -- For WAREHOUSE
#   { MODIFY | MONITOR | USAGE | OPERATE } [ , ... ]
# -- For DATABASE
#   { MODIFY | MONITOR | USAGE | CREATE SCHEMA | IMPORTED PRIVILEGES } [ , ... ]

# schemaPrivileges ::=
#   { MODIFY | MONITOR | USAGE | CREATE { TABLE | VIEW | FILE FORMAT | STAGE | PIPE | SEQUENCE | FUNCTION } } [ , ... ]

# schemaObjectPrivileges ::=
# -- For TABLE
#   { SELECT | INSERT | UPDATE | DELETE | TRUNCATE | REFERENCES } [ , ... ]
# -- For VIEW
#     SELECT
# -- For internal STAGE
#     READ [ , WRITE ]
# -- For external STAGE
#     USAGE
# -- For FILE FORMAT, UDF, or SEQUENCE
#     USAGE

# Given that schema, here are some terraform resources we could support. We should map this out semi-fully, even though
# we will build incrementally.

# set global privs
resource "snowflake_priviledge_grants" "a" {
  name      = "foo1"
  role_name = "role1"

  global {
    all = true

    # these two ^ v are mutually excluseive
    # priviledges = ["create_role", "create_user", "create_warehouse", "create_database", "manage_grants", "monitor_usage"]
    # any subset of these can be supplied ^
  }
}

# set account object privs

# resource monitor
resource "snowflake_priviledge_grants" "b" {
  name      = "foo2"
  role_name = "role2"

  resource_monitor {
    name        = "bar"                 // name of the resource monitor
    priviledges = ["modify", "monitor"]
  }
}

# warehouse
resource "snowflake_priviledge_grants" "c" {
  name      = "foo3"
  role_name = "role3"

  warehouse {
    name        = "bar"                                     // name of the warehouse
    priviledges = ["modify", "monitor", "usage", "operate"]
  }
}

# database
resource "snowflake_priviledge_grants" "d" {
  name      = "foo4"
  role_name = "role4"

  database {
    name        = "bar"                                                                   // name of the database
    priviledges = ["modify", "monitor", "usage", "create_schema", "imported_priviledges"]
  }
}

# database
resource "snowflake_priviledge_grants" "e" {
  name      = "foo5"
  role_name = "role5"

  schema {
    schema_name = "bar"

    # these two ^ v are mutually exclusive
    # database_name = "bar"
    #  if ^ then it means "all schemas in db"

    # not sure if we should make these mutually exclusive or allow you to create one resource w/ a bunch of privs
    table {
      priviledges = ["modify", "monitor", "usage", "create"]
    }
    view = {
      priviledges = ["modify", "monitor", "usage", "create"]
    }
    file_format = {
      priviledges = ["modify", "monitor", "usage", "create"]
    }
    stage = {
      priviledges = ["modify", "monitor", "usage", "create"]
    }
    pipe = {
      priviledges = ["modify", "monitor", "usage", "create"]
    }
    sequence = {
      priviledges = ["modify", "monitor", "usage", "create"]
    }
    function = {
      priviledges = ["modify", "monitor", "usage", "create"]
    }
  }
}

# objects in schema
resource "snowflake_priviledge_grants" "e" {
  name      = "foo6"
  role_name = "role6"

  schema {
    schema_name = "bar"

    #   This has two forms–
    # 1. all objects of a type in a schema
    object_type = "table" # or view, internal_stage, external_stage, file_format, udf or sequence

    # 2a. A specific object in that schema. Option a is to have the name outside the sub-resource…

    table_name          = ".."
    view_name           = ".."
    internal_stage_name = ".."
    external_stage_name = ".."
    file_format_name    = ".."
    udf_name            = "..."
    sequence_name       = "..."
    table {
      priviledges = []
    }
    view = {
      priviledges = []
    }
    internal_stage = {
      priviledges = []
    }
    external_stage = {
      priviledges = []
    }
    file_format = {
      priviledges = []
    }

    # 2b. option b is to nest it


    # also not certain if we should make these mutually exclusive, or even to allow specifying multiple of the same
    # object type in one resource

    table {
      name        = ""
      priviledges = []
    }
    view = {
      name        = ""
      priviledges = []
    }
    internal_stage = {
      name        = ""
      priviledges = []
    }
    external_stage = {
      name        = ""
      priviledges = []
    }
    file_format = {
      name        = ""
      priviledges = []
    }
  }
}
