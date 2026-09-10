# simple usage
data "snowflake_openflow_connectors" "simple" {
}

# filtering (like)
data "snowflake_openflow_connectors" "like" {
  like = "connector-name"
}

# filtering (in)
data "snowflake_openflow_connectors" "in" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

# filtering (starts_with)
data "snowflake_openflow_connectors" "starts_with" {
  starts_with = "connector-"
}

# filtering (limit)
data "snowflake_openflow_connectors" "limit" {
  limit {
    rows = 10
    from = "connector-name"
  }
}

# without additional data (to limit the number of calls make sure to set all of these to false)
data "snowflake_openflow_connectors" "only_show" {
  with_describe = false
}
