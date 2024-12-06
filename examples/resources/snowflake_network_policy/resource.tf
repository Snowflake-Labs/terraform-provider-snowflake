## Minimal
resource "snowflake_network_policy" "basic" {
  name = "network_policy_name"
}

## Complete (with every optional set)
resource "snowflake_network_policy" "complete" {
  name                      = "network_policy_name"
  allowed_network_rule_list = [snowflake_network_rule.one.fully_qualified_name]
  blocked_network_rule_list = [snowflake_network_rule.two.fully_qualified_name]
  allowed_ip_list           = ["192.168.1.0/24"]
  blocked_ip_list           = ["192.168.1.99"]
  comment                   = "my network policy"
}
