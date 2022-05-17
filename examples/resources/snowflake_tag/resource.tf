resource "snowflake_tag" "test_tag" {
    // Required
    name = "tag_name"
    database = "test_db"
    schema = "test_schema"

    // Optionals
    comment = "test comment"
    allowed_values = ["foo", "bar"]
}
