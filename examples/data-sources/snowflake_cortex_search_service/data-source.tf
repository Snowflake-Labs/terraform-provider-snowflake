# Simple usage
data "snowflake_cortex_search_services" "simple" {
}

output "simple_output" {
  value = data.snowflake_cortex_search_services.simple.cortex_search_services
}

# Filtering (like)
data "snowflake_cortex_search_services" "like" {
  like = "some-name"
}

output "like_output" {
  value = data.snowflake_cortex_search_services.like.cortex_search_services
}

# Filtering (starts_with)
data "snowflake_cortex_search_services" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_cortex_search_services.starts_with.cortex_search_services
}

# Filtering (limit)
data "snowflake_cortex_search_services" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_cortex_search_services.limit.cortex_search_services
}
