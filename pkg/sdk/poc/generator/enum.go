package generator

// Enum defines const values
type Enum struct {
	// Kind is fields type (e.g. WarehouseType, LanguageType)
	Kind string
	// Values is a list of possible values (e.g. "XSmall", "Large")
	Values []string
}

func NewEnum(kind string, values []string) *Enum {
	return &Enum{
		Kind:   kind,
		Values: values,
	}
}

func (i *Interface) WithEnums(enums ...Enum) *Interface {
	i.Enums = enums
	return i
}
