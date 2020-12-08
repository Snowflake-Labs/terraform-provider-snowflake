resource snowflake_network_policy_attachment attach {
  network_policy_name = "policy"
  set_for_account     = false
  users = ["user1", "user2"]
}
