package generator

func NewInterface(name string, nameSingular string, identifierKind string, operations ...*Operation) *Interface {
	return &Interface{
		Name:           name,
		NameSingular:   nameSingular,
		IdentifierKind: identifierKind,
		Operations:     operations,
	}
}

func (i *Interface) WithOperations(operations ...*Operation) *Interface {
	i.Operations = operations
	return i
}

func NewOperation(opName string, doc string) *Operation {
	return &Operation{
		Name: opName,
		Doc:  doc,
	}
}

// TODO Do we need that ?
func (s *Operation) WithObjectInterface(objectInterface *Interface) *Operation {
	s.ObjectInterface = objectInterface
	return s
}

func (s *Operation) WithOptsField(optsField *Field) *Operation {
	s.OptsField = optsField
	return s
}

func NewField(name string, kind string, tags map[string][]string) *Field {
	return &Field{
		Name: name,
		Kind: kind,
		Tags: tags,
	}
}

func (field *Field) WithParent(parent *Field) *Field {
	field.Parent = parent
	return field
}

func (field *Field) WithFields(fields ...*Field) *Field {
	field.Fields = fields
	return field
}

func (field *Field) WithValidations(validations ...*Validation) *Field {
	field.Validations = validations
	return field
}

func (field *Field) WithRequired(required bool) *Field {
	field.Required = required
	return field
}

func NewValidation(validationType ValidationType, fieldNames ...string) *Validation {
	return &Validation{
		Type:       validationType,
		FieldNames: fieldNames,
	}
}

/// new
// Static / SQL

func queryTags(ddlTags []string, sqlTags []string) map[string][]string {
	tags := make(map[string][]string)
	if len(ddlTags) > 0 {
		tags["ddl"] = ddlTags
	}
	if len(sqlTags) > 0 {
		tags["sql"] = sqlTags
	}
	return tags
}

func SQL(sql string) *Field {
	return NewField(sqlToFieldName(sql, false), "bool", queryTags([]string{"static"}, []string{sql}))
}

func Create() *Field {
	return SQL("CREATE")
}

func Alter() *Field {
	return SQL("ALTER")
}

func Drop() *Field {
	return SQL("DROP")
}

func Show() *Field {
	return SQL("SHOW")
}

func Describe() *Field {
	return SQL("DESCRIBE")
}

// Keyword / Value

func OptionalSQL(sql string) *Field {
	return NewField(sqlToFieldName(sql, false), "bool", queryTags([]string{"keyword"}, []string{sql}))
}

func OrReplace() *Field {
	return OptionalSQL("OR REPLACE")
}

func IfNotExists() *Field {
	return OptionalSQL("IF NOT EXISTS")
}

func IfExists() *Field {
	return OptionalSQL("IF EXISTS")
}

// Parameters

type parameterOptions struct {
	singleQuotes bool
}

func ParameterOptions() *parameterOptions {
	return &parameterOptions{}
}

func (v *parameterOptions) SingleQuotes(value bool) *parameterOptions {
	v.singleQuotes = value
	return v
}

func (v *parameterOptions) toOptions() []string {
	opts := make([]string, 0)
	if v.singleQuotes {
		opts = append(opts, "single_quotes")
	}
	return opts
}

func OptionalTextAssignment(sqlPrefix string, paramOptions *parameterOptions) *Field {
	if paramOptions != nil {
		return NewField(sqlToFieldName(sqlPrefix, true), "*string", queryTags(append([]string{"parameter"}, paramOptions.toOptions()...), []string{sqlPrefix}))
	}
	return NewField(sqlToFieldName(sqlPrefix, true), "*string", queryTags([]string{"parameter"}, []string{sqlPrefix}))
}

// Identifier

func DatabaseObjectIdentifier(fieldName string) *Field {
	return NewField(fieldName, "DatabaseObjectIdentifier", queryTags([]string{"identifier"}, nil)).WithRequired(true)
}
