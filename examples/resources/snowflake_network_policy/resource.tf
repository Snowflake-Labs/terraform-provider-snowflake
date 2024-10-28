## Minimal
resource "snowflake_network_policy" "basic" {
  name = "network_policy_name"
}

## Complete (with every optional set)
resource "snowflake_network_policy" "complete" {
  name                      = "network_policy_name"
  allowed_network_rule_list = ["<fully qualified network rule id>"]
  blocked_network_rule_list = ["<fully qualified network rule id>"]
  allowed_ip_list           = ["192.168.1.0/24"]
  blocked_ip_list           = ["192.168.1.99"]
  comment                   = "my network policy"
}