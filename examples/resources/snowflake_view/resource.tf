resource "snowflake_view" "view" {
  database = "database"
  schema   = "schema"
  name     = "view"

  comment = "comment"

  statement  = <<-SQL
    select * from foo;
SQL
  or_replace = false
  is_secure  = false
}
