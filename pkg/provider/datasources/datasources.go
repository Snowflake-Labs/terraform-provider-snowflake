package datasources

type datasource string

// TODO: fill out all
const (
	Databases datasource = "snowflake_databases"
)

type Datasource interface {
	xxxProtected()
	String() string
}

func (r datasource) xxxProtected() {}

func (r datasource) String() string {
	return string(r)
}
