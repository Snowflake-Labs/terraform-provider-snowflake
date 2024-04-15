package resources

type resource string

const (
	Account        resource = "snowflake_account"
	Alert          resource = "snowflake_alert"
	ApiIntegration resource = "snowflake_api_integration"
	Database       resource = "snowflake_database"
	View           resource = "snowflake_view"
	Schema         resource = "snowflake_schema"
)

type Resource interface {
	xxxProtected()
	String() string
}

func (r resource) xxxProtected() {}

func (r resource) String() string {
	return string(r)
}
