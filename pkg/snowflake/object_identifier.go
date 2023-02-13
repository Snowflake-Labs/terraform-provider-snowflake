package snowflake

// ObjectIdentifier is a helper struct for identifying objects
type ObjectIdentifier struct {
	Name     string
	Database string
	Schema   string
}
