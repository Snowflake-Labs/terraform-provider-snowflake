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

# Filtering (in)
data "snowflake_secrets" "in" {
  in {
    schema = "snowflake_schema.test.fully_qualified_name"
  }
}

output "in_output" {
  value = data.snowflake_secrets.in.secrets
}

# Without additional data (to limit the number of calls make for every found secret)
data "snowflake_secrets" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SECRET for every secret found and attaches its output to secrets.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_secrets.only_show.secrets
}
