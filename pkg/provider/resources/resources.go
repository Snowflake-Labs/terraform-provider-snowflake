package resources

type resource string

const (
	View   resource = "snowflake_view"
	Schema resource = "snowflake_schema"
)

type Resource interface {
	xxxProtected()
	String() string
}

func (r resource) xxxProtected() {}

func (r resource) String() string {
	return string(r)
}
