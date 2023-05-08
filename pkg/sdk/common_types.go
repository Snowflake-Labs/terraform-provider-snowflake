package sdk

// TimeTravel is an enum for the time travel options. AT | BEFORE are supported
type TimeTravel struct {
	Timestamp string `ddl:"command" db:"TIMESTAMP =>"`
	Offset    string `ddl:"command" db:"OFFSET =>"`
	Statement string `ddl:"command,single_quotes" db:"STATEMENT =>"`
}

type Clone struct {
	SourceObject Identifier  `ddl:"identifier"`
	At           *TimeTravel `ddl:"AT"`
	Before       *TimeTravel `ddl:"BEFORE"`
}

type In struct {
	Account  *bool                  `ddl:"keyword" db:"ACCOUNT"`
	Database AccountLevelIdentifier `ddl:"identifier" db:"DATABASE"`
	Schema   SchemaIdentifier       `ddl:"identifier" db:"SCHEMA"`
}

type Like struct {
	Pattern *string `ddl:"keyword,single_quotes"`
}

type TagAssociation struct {
	Name  Identifier `ddl:"identifier"`
	eq    bool       `ddl:"static" db:"="` //lint:ignore U1000 This is used in the ddl tag
	Value string     `ddl:"keyword,single_quotes"`
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
