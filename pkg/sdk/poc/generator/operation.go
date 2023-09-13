package generator

// Operation defines a single operation for given object or objects family (e.g. CREATE DATABASE ROLE)
type Operation struct {
	// Name is the operation's name, e.g. "Create"
	Name string
	// ObjectInterface points to the containing interface
	ObjectInterface *Interface
	// Doc is the URL for the doc used to create given operation, e.g. https://docs.snowflake.com/en/sql-reference/sql/create-database-role
	Doc string
	// OptsField defines opts used to create SQL for given operation
	OptsField *Field
}

func NewOperation(opName string, doc string) *Operation {
	return &Operation{
		Name: opName,
		Doc:  doc,
	}
}

func (s *Operation) WithOptionsStruct(optsField *Field) *Operation {
	s.OptsField = optsField
	return s
}
