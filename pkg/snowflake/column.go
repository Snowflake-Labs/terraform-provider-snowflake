package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ColumnDataType string

//"VARCHAR(16777216)", "VARCHAR", "text", "string", "NVARCHAR", "NVARCHAR2", "CHAR VARYING", "NCHAR VARYING"

const (
	ColumnDataTypeVarchar16777216 ColumnDataType = "VARCHAR(16777216)"
	ColumnDataTypeVarchar         ColumnDataType = "VARCHAR"
	ColumnDataTypeText            ColumnDataType = "TEXT"
	ColumnDataTypeString          ColumnDataType = "STRING"
	ColumnDataTypeNvarchar        ColumnDataType = "NVARCHAR"
	ColumnDataTypeNvarchar2       ColumnDataType = "NVARCHAR2"
	ColumnDataTypeCHARVARYING     ColumnDataType = "CHAR VARYING"
	ColumnDataTypeNCHARVARYING    ColumnDataType = "NCHAR VARYING"
)

type Identity struct {
	startNum int
	stepNum  int
}

type TableColumnBuilder struct {
	name            string
	tableIdentifier *ObjectIdentifier
	dataType        *ColumnDataType
	nullable        bool
	default_        *ColumnDefault
	identity        *Identity
	comment         string
	db              *sql.DB
}

func NewCreateTableColumnBuilder(name string, tableIdentifier *ObjectIdentifier, dataType *ColumnDataType, db *sql.DB) *TableColumnBuilder {
	return &TableColumnBuilder{
		name:            name,
		tableIdentifier: tableIdentifier,
		dataType:        dataType,
		db:              db,
	}
}

func (v *TableColumnBuilder) WithNullable(nullable bool) *TableColumnBuilder {
	v.nullable = nullable
	return v
}

func (v *TableColumnBuilder) WithDefault(default_ *ColumnDefault) *TableColumnBuilder {
	v.default_ = default_
	return v
}

func (v *TableColumnBuilder) WithIdentity(identity *Identity) *TableColumnBuilder {
	v.identity = identity
	return v
}

func (v *TableColumnBuilder) WithComment(comment string) *TableColumnBuilder {
	v.comment = comment
	return v
}

//func (v *TableColumnBuilder) WithDefaultExpression(defaultExpression string) *TableColumnBuilder {
//	v.defaultExpression = defaultExpression
//	return v
//}
//
//func (v *TableColumnBuilder) WithDefaultSequence(defaultSequence *ObjectIdentifier) *TableColumnBuilder {
//	v.defaultSequence = defaultSequence
//	return v
//}

type TableColumn struct {
	TableName     sql.NullString `db:"table_name"`
	SchemaName    sql.NullString `db:"schema_name"`
	ColumnName    sql.NullString `db:"column_name"`
	DataType      sql.NullString `db:"data_type"`
	Null          sql.NullBool   `db:"null?"`
	Default_      sql.NullString `db:"default"`
	Kind          sql.NullString `db:"kind"`
	Expression    sql.NullString `db:"expression"`
	Comment       sql.NullString `db:"comment"`
	DatabaseName  sql.NullString `db:"database_name"`
	Autoincrement sql.NullString `db:"autoincrement"`
}

func (v *TableColumnBuilder) Show() (*TableColumn, error) {
	var value TableColumn

	stmt := fmt.Sprintf("SHOW COLUMNS LIKE '%s' IN TABLE '%s.%s.%s'", v.name, v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name)
	rows, err := v.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tableColumns []TableColumn
	if err := sqlx.StructScan(rows, &tableColumns); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	value = tableColumns[0]

	return &value, nil
}

func (v *TableColumnBuilder) Drop() error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' DROP COLUMN '%s'", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name)
	_, err := v.db.Exec(stmt)
	return err
}

func (v *TableColumnBuilder) Rename(newName string) error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' RENAME COLUMN '%s' TO '%s'", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, newName)
	_, err := v.db.Exec(stmt)
	return err
}

func (v *TableColumnBuilder) SetComment() error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' COMMENT '%s'", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, v.comment)
	_, err := v.db.Exec(stmt)
	return err
}

func (v *TableColumnBuilder) UnsetComment() error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' UNSET COMMENT", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name)
	_, err := v.db.Exec(stmt)
	return err
}

//func (v *TableColumnBuilder) SetMaskingPolicy(maskingPolicy string) error {
//	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' MASKING POLICY '%s'", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, maskingPolicy)
//	_, err := v.db.Exec(stmt)
//	return err
//}
//
//func (v *TableColumnBuilder) UnsetMaskingPolicy() error {
//	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' MASKING POLICY", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name)
//	_, err := v.db.Exec(stmt)
//	return err
//}

func (v *TableColumnBuilder) SetNotNull() error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' SET NOT NULL", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, v.default_.String(fmt.Sprintf("%s", *v.dataType)))
	_, err := v.db.Exec(stmt)
	return err
}

func (v *TableColumnBuilder) DropNotNull() error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' DROP NOT NULL", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, v.default_.String(fmt.Sprintf("%s", *v.dataType)))
	_, err := v.db.Exec(stmt)
	return err
}

func (v *TableColumnBuilder) SetDefault() error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' SET DEFAULT '%s'", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, v.default_.String(fmt.Sprintf("%s", *v.dataType)))
	_, err := v.db.Exec(stmt)
	return err
}

//func (v *TableColumnBuilder) SetDefaultExpression(sequenceIdentifier *ObjectIdentifier) error {
//	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' SET DEFAULT '%s.%s.%s'.NEXTVAL ", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, sequenceIdentifier.Database, sequenceIdentifier.Schema, sequenceIdentifier.Name)
//	_, err := v.db.Exec(stmt)
//	return err
//}
//
//func (v *TableColumnBuilder) SetDefaultSequence(sequenceIdentifier *ObjectIdentifier) error {
//	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' SET DEFAULT '%s.%s.%s'.NEXTVAL ", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, sequenceIdentifier.Database, sequenceIdentifier.Schema, sequenceIdentifier.Name)
//	_, err := v.db.Exec(stmt)
//	return err
//}

func (v *TableColumnBuilder) DropDefault() error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' DROP DEFAULT", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name)
	_, err := v.db.Exec(stmt)
	return err
}

func (v *TableColumnBuilder) SetIdentity() error {
	stmt := fmt.Sprintf("ALTER TABLE '%s.%s.%s' COLUMN '%s' SET IDENTITY(%d, %d)", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, v.identity.startNum, v.identity.stepNum)
	_, err := v.db.Exec(stmt)
	return err
}

func (v *TableColumnBuilder) DropIdentity() error {
	return v.DropDefault()
}

func (v *TableColumnBuilder) Create() error {
	q := strings.Builder{}

	q.WriteString(fmt.Sprintf("ALTER TABLE '%s.%s.%s' ADD COLUMN '%s' '%s'", v.tableIdentifier.Database, v.tableIdentifier.Schema, v.tableIdentifier.Name, v.name, *v.dataType))

	stmt := q.String()
	_, err := v.db.Exec(stmt)
	if err != nil {
		return err
	}

	if !v.nullable {
		err = v.SetNotNull()
		if err != nil {
			return err
		}
	}

	if v.default_ != nil {
		err = v.SetDefault()
		if err != nil {
			return err
		}
	}

	if v.identity != nil {
		err = v.SetIdentity()
		if err != nil {
			return err
		}
	}

	err = v.SetComment()

	return err
}
