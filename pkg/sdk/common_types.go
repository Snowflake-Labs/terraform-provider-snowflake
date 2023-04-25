package sdk

type In struct {
	Account  *bool   `ddl:"keyword"`
	Database *string `ddl:"command,double_quotes" db:"DATABASE"`
	Schema   *string `ddl:"command,double_quotes" db:"SCHEMA"`
}

type Like struct {
	Pattern *string `ddl:"keyword,single_quotes"`
}

type DescribeStringProperty struct {
	Value        string
	DefaultValue string
	Description  string
}

type DescribeIntProperty struct {
	Value        int
	DefaultValue int
	Description  string
}

type describePropertyRow struct {
	Property     string `db:"property"`
	Value        string `db:"value"`
	DefaultValue string `db:"default"`
	Description  string `db:"description"`
}

func (row *describePropertyRow) toDescribeStringProperty() *DescribeStringProperty {
	return &DescribeStringProperty{
		Value:        row.Value,
		DefaultValue: row.DefaultValue,
		Description:  row.Description,
	}
}

func (row *describePropertyRow) toDescribeIntProperty() *DescribeIntProperty {
	return &DescribeIntProperty{
		Value:        toInt(row.Value),
		DefaultValue: toInt(row.DefaultValue),
		Description:  row.Description,
	}
}
