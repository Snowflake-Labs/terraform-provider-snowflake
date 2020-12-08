resource snowflake_network_poilcy policy {
  name    = "policy"
  comment = "A policy."

  allowed_ip_list = ["192.168.0.100/24"]
  blocked_ip_list = ["192.168.0.101"]
}
