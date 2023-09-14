package generator

type OperationKind string

const (
	OperationKindCreate   OperationKind = "Create"
	OperationKindAlter    OperationKind = "Alter"
	OperationKindDrop     OperationKind = "Drop"
	OperationKindShow     OperationKind = "Show"
	OperationKindShowByID OperationKind = "ShowByID"
	OperationKindDescribe OperationKind = "Describe"
)

// Operation defines a single operation for given object or objects family (e.g. CREATE DATABASE ROLE)
type Operation struct {
	// Name is the operation's name, e.g. "Create"
	Name OperationKind
	// ObjectInterface points to the containing interface
	ObjectInterface *Interface
	// Doc is the URL for the doc used to create given operation, e.g. https://docs.snowflake.com/en/sql-reference/sql/create-database-role
	Doc string
	// OptsField defines opts used to create SQL for given operation
	OptsField *Field
	// HelperStructs are struct definitions that are not tied to OptsField, but tied to the Operation itself, e.g. Show() return type
	HelperStructs []*Field
	ShowMapping   *Mapping
	//CustomMappings []*Mapping
}

type Mapping struct {
	MappingFuncName string
	From            *Field
	To              *Field
}

func NewOperation(kind OperationKind, doc string) *Operation {
	return &Operation{
		Name:          kind,
		Doc:           doc,
		HelperStructs: make([]*Field, 0),
	}
}

func NewMapping(mappingFuncName string, from, to *Field) *Mapping {
	return &Mapping{
		MappingFuncName: mappingFuncName,
		From:            from,
		To:              to,
	}
}

func (s *Operation) withOptionsStruct(optsField *Field) *Operation {
	s.OptsField = optsField
	return s
}

func (s *Operation) withHelperStruct(helperStruct *Field) *Operation {
	s.HelperStructs = append(s.HelperStructs, helperStruct)
	return s
}

func (s *Operation) withShowMapping(from, to *Field) *Operation {
	s.ShowMapping = NewMapping("convert", from, to)
	return s
}

//func (s *Operation) withMapping(mappingFuncName string, from, to *Field) *Operation {
//	s.CustomMappings = append(s.CustomMappings, NewMapping(mappingFuncName, from, to))
//	return s
//}

func CreateOperation(doc string, queryStruct *Field) *Operation {
	return NewOperation(OperationKindCreate, doc).withOptionsStruct(queryStruct)
}

func ShowOperation(doc string, dbRepresentation *dbStruct, resourceRepresentation *plainStruct, queryStruct *Field) *Operation {
	db := dbRepresentation.IntoField()
	res := resourceRepresentation.IntoField()
	return NewOperation(OperationKindShow, doc).
		withHelperStruct(db).
		withHelperStruct(res).
		withShowMapping(db, res).
		withOptionsStruct(queryStruct)
}
