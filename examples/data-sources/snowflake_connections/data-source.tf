# Simple usage
data "snowflake_connections" "simple" {
}

output "simple_output" {
  value = data.snowflake_connections.simple.connections
}

# Filtering (like)
data "snowflake_connections" "like" {
  like = "connection-name"
}

output "like_output" {
  value = data.snowflake_connections.like.connections
}

# Filtering by prefix (like)
data "snowflake_connections" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_connections.like_prefix.connections
}
