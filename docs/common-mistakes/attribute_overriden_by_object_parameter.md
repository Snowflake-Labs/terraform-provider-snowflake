# Resource attribute overriden by object parameter resource

## Problem

Adding an object parameter resource for a Snowflake object which already is managed by another resource might result in
unexpected behaviours.

### Example

```terraform
resource "snowflake_database" "d" {
  name = "TEST_DB"
}

resource "snowflake_schema" "s" {
  name     = "TEST_SCHEMA"
  database = snowflake_database.d.name
}

resource "snowflake_table" "t" {
  name     = "TEST_TABLE"
  database = snowflake_database.d.name
  schema   = snowflake_schema.s.name
  column {
    name = "id"
    type = "NUMBER"
  }
}

resource "snowflake_object_parameter" "o" {
  key         = "DATA_RETENTION_TIME_IN_DAYS"
  value       = "89"
  object_type = "TABLE"
  object_identifier {
    database = snowflake_database.d.name
    schema   = snowflake_schema.s.name
    name     = snowflake_table.t.name
  }
}
```

In the example above we define a Snowflake table and an object parameter that sets `DATA_RETENTION_TIME_IN_DAYS`
parameter of that table. Since table resource schema also allows to set `DATA_RETENTION_TIME_IN_DAYS` parameter with an
optional `data_retention_time_in_days` attribute, value for the parameter will be managed by two different resources.
This can cause confusion, as one value will be overriden by another.

## Solution 1
First solution is to use only table resource and set the value for `DATA_RETENTION_TIME_IN_DAYS` parameter using table's `data_retention_time_in_days` attribute.
```terraform
resource "snowflake_database" "d" {
  name = "TEST_DB"
}

resource "snowflake_schema" "s" {
  name     = "TEST_SCHEMA"
  database = snowflake_database.d.name
}

resource "snowflake_table" "t" {
  name     = "TEST_TABLE"
  database = snowflake_database.d.name
  schema   = snowflake_schema.s.name
  data_retention_time_in_days = 89
  column {
    name = "id"
    type = "NUMBER"
  }
}
```

## Solution 2
If we do want to have a separate object parameter resource for managing `DATA_RETENTION_TIME_IN_DAYS` parameter, we can instruct terraform to ignore `data_retention_time_in_days` attribute of the table resource by adding a lifecycle argument

```terraform
resource "snowflake_database" "d" {
  name = "TEST_DB"
}

resource "snowflake_schema" "s" {
  name     = "TEST_SCHEMA"
  database = snowflake_database.d.name
}

resource "snowflake_table" "t" {
  name     = "TEST_TABLE"
  database = snowflake_database.d.name
  schema   = snowflake_schema.s.name
  column {
    name = "id"
    type = "NUMBER"
  }
  lifecycle {
    ignore_changes = [
      "data_retention_time_in_days"
    ]
  }
}

resource "snowflake_object_parameter" "o" {
  key         = "DATA_RETENTION_TIME_IN_DAYS"
  value       = "89"
  object_type = "TABLE"
  object_identifier {
    database = snowflake_database.d.name
    schema   = snowflake_schema.s.name
    name     = snowflake_table.t.name
  }
}
```
Fore more info on lifecyle attributes, check https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#ignore_changes
