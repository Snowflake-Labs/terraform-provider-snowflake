##################################
### using network rules
##################################

resource "snowflake_network_rule" "rule" {
  name       = "rule"
  database   = "EXAMPLE_DB"
  schema     = "EXAMPLE_SCHEMA"
  comment    = "A rule."
  type       = "IPV4"
  mode       = "INGRESS"
  value_list = ["192.168.0.100/24", "29.254.123.20"]
}

resource "snowflake_network_policy" "policy" {
  name    = "policy"
  comment = "A policy."

  allowed_network_rule_list = [snowflake_network_rule.rule.qualified_name]
}


##################################
### using ip lists
##################################

resource "snowflake_network_policy" "policy" {
  name    = "policy"
  comment = "A policy."

  allowed_ip_list = ["192.168.0.100/24"]
  blocked_ip_list = ["192.168.0.101"]
}
