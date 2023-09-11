package generator2

// Operation defines a single operation for given object or objects family (e.g. CREATE DATABASE ROLE)
type Operation struct {
	// Name is the operation's name, e.g. "Create"
	Name string
	// ObjectInterface points to the containing interface
	ObjectInterface *Interface
	// Doc is the URL for the doc used to create given operation, e.g. https://docs.snowflake.com/en/sql-reference/sql/create-database-role
	Doc string
	// Options TODO
	Options *Struct
	// UtilStructs used to form more complex queries in Options struct
	UtilStructs []*Struct
}

func NewOperation(opName string, doc string, optionsStruct *Struct, utilStructs ...*Struct) *Operation {
	return &Operation{
		Name:        opName,
		Doc:         doc,
		Options:     optionsStruct,
		UtilStructs: utilStructs,
	}
}
