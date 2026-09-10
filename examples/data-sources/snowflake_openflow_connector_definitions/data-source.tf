# simple usage
data "snowflake_openflow_connector_definitions" "simple" {
}

# filtering (like)
data "snowflake_openflow_connector_definitions" "like" {
  like = "OPENFLOW_POSTGRES_CDC"
}

# filtering (limit)
data "snowflake_openflow_connector_definitions" "limit" {
  limit {
    rows = 10
    from = "OPENFLOW_"
  }
}
