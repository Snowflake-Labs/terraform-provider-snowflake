# Simple usage
data "snowflake_secrets" "simple" {
}

output "simple_output" {
  value = data.snowflake_secrets.simple.secrets
}

# Filtering (like)
data "snowflake_secrets" "like" {
  like = "secret-name"
}

output "like_output" {
  value = data.snowflake_secrets.like.secrets
}

# Filtering by prefix (like)
data "snowflake_secrets" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_secrets.like_prefix.secrets
}
