package generator

func NewInterface(
	Name string,
	NameSingular string,
	IdentifierKind string,
) *Interface {
	s := Interface{}
	s.Name = Name
	s.NameSingular = NameSingular
	s.IdentifierKind = IdentifierKind
	return &s
}

func (i *Interface) WithOperations(Operations []*Operation) *Interface {
	i.Operations = Operations
	return i
}

func NewOperation(
	Name string,
	Doc string,
) *Operation {
	s := Operation{}
	s.Name = Name
	s.Doc = Doc
	return &s
}

func (s *Operation) WithObjectInterface(ObjectInterface *Interface) *Operation {
	s.ObjectInterface = ObjectInterface
	return s
}

func (s *Operation) WithOptsField(OptsField *Field) *Operation {
	s.OptsField = OptsField
	return s
}

func NewField(
	Name string,
	Kind string,
	Tags map[string][]string,
) *Field {
	s := Field{}
	s.Name = Name
	s.Kind = Kind
	s.Tags = Tags
	return &s
}

func (field *Field) WithParent(Parent *Field) *Field {
	field.Parent = Parent
	return field
}

func (field *Field) WithFields(Fields []*Field) *Field {
	field.Fields = Fields
	return field
}

func (field *Field) WithValidations(Validations []*Validation) *Field {
	field.Validations = Validations
	return field
}

func (field *Field) WithRequired(Required bool) *Field {
	field.Required = Required
	return field
}

func NewValidation(
	Type ValidationType,
	FieldNames []string,
) *Validation {
	s := Validation{}
	s.Type = Type
	s.FieldNames = FieldNames
	return &s
}
