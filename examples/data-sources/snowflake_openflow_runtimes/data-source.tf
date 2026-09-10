# simple usage
data "snowflake_openflow_runtimes" "simple" {
}

# filtering (like)
data "snowflake_openflow_runtimes" "like" {
  like = "runtime-name"
}

# filtering (in)
data "snowflake_openflow_runtimes" "in" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

# filtering (starts_with)
data "snowflake_openflow_runtimes" "starts_with" {
  starts_with = "runtime-"
}

# filtering (limit)
data "snowflake_openflow_runtimes" "limit" {
  limit {
    rows = 10
    from = "runtime-name"
  }
}

# without additional data (to limit the number of calls make sure to set all of these to false)
data "snowflake_openflow_runtimes" "only_show" {
  with_describe = false
}
