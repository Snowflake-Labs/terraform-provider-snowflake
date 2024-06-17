data "snowflake_databases" "test" {
  like = "non-existing-database"

  lifecycle {
    postcondition {
      condition     = length(self.databases) > 0
      error_message = "there should be at least one database"
    }
  }
}
