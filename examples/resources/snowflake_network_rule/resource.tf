resource "snowflake_network_rule" "rule" {
  name       = "rule"
  database   = "EXAMPLE_DB"
  schema     = "EXAMPLE_SCHEMA"
  comment    = "A rule."
  type       = "IPV4"
  mode       = "INGRESS"
  value_list = ["192.168.0.100/24", "29.254.123.20"]
}
