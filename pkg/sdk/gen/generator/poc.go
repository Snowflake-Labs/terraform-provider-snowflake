package generator

import (
	"fmt"
	"strings"
)

// Interface groups operations for particular object or objects family (e.g. DATABASE ROLE)
type Interface struct {
	// Name is the interface's name, e.g. "DatabaseRoles"
	Name string
	// nameSingular is the prefix/suffix which can be used to create other structs and methods, e.g. "DatabaseRole"
	nameSingular string
	// Operations contains all operations for given interface
	Operations []*Operation
}

// Operation defines a single operation for given object or objects family (e.g. CREATE DATABASE ROLE)
type Operation struct {
	// Name is the operation's name, e.g. "Create"
	Name string
	// ObjectInterface points to the containing interface
	ObjectInterface *Interface
	// Doc is the URL for the doc used to create given operation, e.g. https://docs.snowflake.com/en/sql-reference/sql/create-database-role
	Doc string
	// OptsStructFields defines opts used to create SQL for given operation
	OptsStructFields []*Field
}

// OptsName should create a name for opts in a form of OperationObjectOptions where:
// - Operation is e.g. Create
// - Object is e.g. DatabaseRole (singular)
// which together makes CreateDatabaseRoleOptions
func (o *Operation) OptsName() string {
	return fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.nameSingular)
}

type Field struct {
	parent *Field
	Fields []*Field

	Name string
	Kind string
	tags map[string][]string
}

func (f *Field) TagsPrintable() string {
	var tagNames = []string{"ddl", "sql"}
	var tagParts []string
	for _, tagName := range tagNames {
		var v, ok = f.tags[tagName]
		if ok {
			tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, tagName, strings.Join(v, ",")))
		}
	}
	return fmt.Sprintf("`%s`", strings.Join(tagParts, " "))
}

func (f *Field) KindNoPtr() string {
	kindWithoutPtr, _ := strings.CutPrefix(f.Kind, "*")
	return kindWithoutPtr
}
