# simple usage
data "snowflake_openflow_deployments" "simple" {
}

# filtering (like)
data "snowflake_openflow_deployments" "like" {
  like = "deployment-name"
}

# filtering (starts_with)
data "snowflake_openflow_deployments" "starts_with" {
  starts_with = "deployment-"
}

# filtering (limit)
data "snowflake_openflow_deployments" "limit" {
  limit {
    rows = 10
    from = "deployment-name"
  }
}

# without additional data (to limit the number of calls make sure to set all of these to false)
data "snowflake_openflow_deployments" "only_show" {
  with_describe   = false
  with_parameters = false
}
