## Minimal
resource "snowflake_role" "minimal" {
  name = "role_name"
}

## Complete (with every optional set)
resource "snowflake_role" "complete" {
  name    = "role_name"
  comment = "my account role"
}
