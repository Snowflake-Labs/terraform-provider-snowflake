package generator

func NewInterface(
	name string,
	nameSingular string,
	identifierKind string,
) *Interface {
	s := Interface{}
	s.Name = name
	s.NameSingular = nameSingular
	s.IdentifierKind = identifierKind
	return &s
}

func (i *Interface) WithOperations(operations []*Operation) *Interface {
	i.Operations = operations
	return i
}

func NewOperation(
	name string,
	doc string,
) *Operation {
	s := Operation{}
	s.Name = name
	s.Doc = doc
	return &s
}

func (s *Operation) WithObjectInterface(objectInterface *Interface) *Operation {
	s.ObjectInterface = objectInterface
	return s
}

func (s *Operation) WithOptsField(optsField *Field) *Operation {
	s.OptsField = optsField
	return s
}

func NewField(
	name string,
	kind string,
	tags map[string][]string,
) *Field {
	s := Field{}
	s.Name = name
	s.Kind = kind
	s.Tags = tags
	return &s
}

func (field *Field) WithParent(parent *Field) *Field {
	field.Parent = parent
	return field
}

func (field *Field) WithFields(fields []*Field) *Field {
	field.Fields = fields
	return field
}

func (field *Field) WithValidations(validations []*Validation) *Field {
	field.Validations = validations
	return field
}

func (field *Field) WithRequired(required bool) *Field {
	field.Required = required
	return field
}

func NewValidation(
	vType ValidationType,
	fieldNames []string,
) *Validation {
	s := Validation{}
	s.Type = vType
	s.FieldNames = fieldNames
	return &s
}
