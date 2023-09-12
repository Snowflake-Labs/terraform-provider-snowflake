package generator

func NewInterface(name string, nameSingular string, identifierKind string, operations ...*Operation) *Interface {
	return &Interface{
		Name:           name,
		NameSingular:   nameSingular,
		IdentifierKind: identifierKind,
		Operations:     operations,
	}
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

// NewOptionsStruct factory method for creating top level fields (Option Structs)
func NewOptionsStruct() *Field {
	return NewField("", "", nil)
}

func NewField(name string, kind string, tags *TagBuilder) *Field {
	var tagsResult map[string][]string
	if tags != nil {
		tagsResult = tags.Build()
	}
	return &Field{
		Name: name,
		Kind: kind,
		Tags: tagsResult,
	}
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

type TagBuilder struct {
	ddl []string
	sql []string
}

func Tags() *TagBuilder {
	return &TagBuilder{
		ddl: make([]string, 0),
		sql: make([]string, 0),
	}
}

func (v *TagBuilder) Static() *TagBuilder {
	v.ddl = append(v.ddl, "static")
	return v
}

func (v *TagBuilder) Keyword() *TagBuilder {
	v.ddl = append(v.ddl, "keyword")
	return v
}

func (v *TagBuilder) Parameter() *TagBuilder {
	v.ddl = append(v.ddl, "parameter")
	return v
}

func (v *TagBuilder) Identifier() *TagBuilder {
	v.ddl = append(v.ddl, "identifier")
	return v
}

func (v *TagBuilder) List() *TagBuilder {
	v.ddl = append(v.ddl, "list")
	return v
}

func (v *TagBuilder) NoParentheses() *TagBuilder {
	v.ddl = append(v.ddl, "no_parentheses")
	return v
}

func (v *TagBuilder) DDL(ddl ...string) *TagBuilder {
	v.ddl = append(v.ddl, ddl...)
	return v
}

func (v *TagBuilder) SQL(sql ...string) *TagBuilder {
	v.sql = append(v.sql, sql...)
	return v
}

func (v *TagBuilder) Build() map[string][]string {
	return map[string][]string{
		"ddl": v.ddl,
		"sql": v.sql,
	}
}

// Static / SQL

func SQL(sql string) *Field {
	return NewField(sqlToFieldName(sql, false), "bool", Tags().Static().SQL(sql))
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
	return NewField(sqlToFieldName(sql, true), "*bool", Tags().Keyword().SQL(sql))
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

func Text(name string, tags *TagBuilder) *Field {
	return NewField(name, "string", tags)
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

func Assignment(sqlPrefix string, kind string, options *parameterOptions) *Field {
	if options != nil {
		return NewField(sqlToFieldName(sqlPrefix, true), kind, Tags().Parameter().DDL(options.toOptions()...).SQL(sqlPrefix))
	}
	return NewField(sqlToFieldName(sqlPrefix, true), kind, Tags().Parameter().SQL(sqlPrefix))
}

func TextAssignment(sqlPrefix string, paramOptions *parameterOptions) *Field {
	return Assignment(sqlPrefix, "string", paramOptions)
}

func OptionalTextAssignment(sqlPrefix string, paramOptions *parameterOptions) *Field {
	return Assignment(sqlPrefix, "*string", paramOptions)
}

// Identifier

func AccountObjectIdentifier(fieldName string) *Field {
	return NewField(fieldName, "AccountObjectIdentifier", Tags().Identifier()).WithRequired(true)
}

func DatabaseObjectIdentifier(fieldName string) *Field {
	return NewField(fieldName, "DatabaseObjectIdentifier", Tags().Identifier()).WithRequired(true)
}
