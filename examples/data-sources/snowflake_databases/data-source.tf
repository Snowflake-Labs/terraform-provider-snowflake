data "snowflake_databases" "test" {
  with_describe   = false
  with_parameters = false
  like            = "database-name"
  starts_with     = "database-name"
  limit {
    rows = 20
    from = "database-name"
  }

  lifecycle {
    postcondition {
      condition     = length(self.databases) > 0
      error_message = "there should be at least one database"
    }
  }
}

check "database_check" {
  data "snowflake_databases" "test" {
    like = "database-name"
  }

  assert {
    condition     = length(data.snowflake_databases.test.databases) == 1
    error_message = "Databases fieltered by '${data.snowflake_databases.test.like}' returned ${length(data.snowflake_databases.test.databases)} databases where one was expected"
  }
}
