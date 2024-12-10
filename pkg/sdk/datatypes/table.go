package datatypes

// TableDataType is based on TODO [this PR]
// It does not have synonyms.
// It consists of a list of column name + column type; may be empty.
// TODO [this PR]: test and improve
type TableDataType struct {
	columns        []TableDataTypeColumn
	underlyingType string
}

type TableDataTypeColumn struct {
	name     string
	dataType DataType
}

func (c *TableDataTypeColumn) ColumnName() string {
	return c.name
}

func (c *TableDataTypeColumn) ColumnType() DataType {
	return c.dataType
}

func (t *TableDataType) ToSql() string {
	return t.underlyingType
}

func (t *TableDataType) ToLegacyDataTypeSql() string {
	return TableLegacyDataType
}

func (t *TableDataType) Canonical() string {
	return TableLegacyDataType
}

func (t *TableDataType) Columns() []TableDataTypeColumn {
	return t.columns
}
