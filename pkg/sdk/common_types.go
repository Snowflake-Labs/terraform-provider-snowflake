package sdk

type In struct {
	Account  *bool                   `ddl:"keyword" db:"ACCOUNT"`
	Database AccountObjectIdentifier `ddl:"identifier" db:"DATABASE"`
	Schema   SchemaIdentifier        `ddl:"identifier" db:"SCHEMA"`
}

type Like struct {
	Pattern *string `ddl:"keyword,single_quotes"`
}

type TagAssociation struct {
	Name  ObjectIdentifier
	Value string
}

type TableColumnSignature struct {
	Name string   `ddl:"keyword,double_quotes"`
	Type DataType `ddl:"keyword"`
}

type StringProperty struct {
	Value        string
	DefaultValue string
	Description  string
}

type IntProperty struct {
	Value        int
	DefaultValue int
	Description  string
}

type propertyRow struct {
	Property     string `db:"property"`
	Value        string `db:"value"`
	DefaultValue string `db:"default"`
	Description  string `db:"description"`
}

func (row *propertyRow) toStringProperty() *StringProperty {
	if row.Value == "null" {
		row.Value = ""
	}
	if row.DefaultValue == "null" {
		row.DefaultValue = ""
	}
	return &StringProperty{
		Value:        row.Value,
		DefaultValue: row.DefaultValue,
		Description:  row.Description,
	}
}

func (row *propertyRow) toIntProperty() *IntProperty {
	return &IntProperty{
		Value:        toInt(row.Value),
		DefaultValue: toInt(row.DefaultValue),
		Description:  row.Description,
	}
}
