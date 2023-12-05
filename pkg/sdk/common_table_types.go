package sdk

import "errors"

type RowAccessPolicy struct {
	rowAccessPolicy bool                   `ddl:"static" sql:"ROW ACCESS POLICY"`
	Name            SchemaObjectIdentifier `ddl:"identifier"`
	On              []string               `ddl:"keyword,parentheses" sql:"ON"`
}

// ColumnInlineConstraint is based on https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint#inline-unique-primary-foreign-key.
type ColumnInlineConstraint struct {
	NotNull    *bool                 `ddl:"keyword" sql:"NOT NULL"`
	Name       *string               `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       *ColumnConstraintType `ddl:"keyword"`
	ForeignKey *InlineForeignKey     `ddl:"keyword" sql:"FOREIGN KEY"`

	// optional
	Enforced           *bool `ddl:"keyword" sql:"ENFORCED"`
	NotEnforced        *bool `ddl:"keyword" sql:"NOT ENFORCED"`
	Deferrable         *bool `ddl:"keyword" sql:"DEFERRABLE"`
	NotDeferrable      *bool `ddl:"keyword" sql:"NOT DEFERRABLE"`
	InitiallyDeferred  *bool `ddl:"keyword" sql:"INITIALLY DEFERRED"`
	InitiallyImmediate *bool `ddl:"keyword" sql:"INITIALLY IMMEDIATE"`
	Enable             *bool `ddl:"keyword" sql:"ENABLE"`
	Disable            *bool `ddl:"keyword" sql:"DISABLE"`
	Validate           *bool `ddl:"keyword" sql:"VALIDATE"`
	NoValidate         *bool `ddl:"keyword" sql:"NOVALIDATE"`
	Rely               *bool `ddl:"keyword" sql:"RELY"`
	NoRely             *bool `ddl:"keyword" sql:"NORELY"`
}

func (v *ColumnInlineConstraint) validate() error {
	// TODO[SNOW-934647]: type required
	var errs []error
	if *v.Type == ColumnConstraintTypeForeignKey {
		if !valueSet(v.ForeignKey) {
			errs = append(errs, errNotSet("ColumnInlineConstraint", "ForeignKey"))
		}
	} else {
		if valueSet(v.ForeignKey) {
			errs = append(errs, errSet("ColumnInlineConstraint", "ForeignKey"))
		}
	}
	if moreThanOneValueSet(v.Enforced, v.NotEnforced) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Enforced", "NotEnforced"))
	}
	if moreThanOneValueSet(v.Deferrable, v.NotDeferrable) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Deferrable", "NotDeferrable"))
	}
	if moreThanOneValueSet(v.InitiallyDeferred, v.InitiallyImmediate) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "InitiallyDeferred", "InitiallyImmediate"))
	}
	if moreThanOneValueSet(v.Enable, v.Disable) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Enable", "Disable"))
	}
	if moreThanOneValueSet(v.Validate, v.NoValidate) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Validate", "Novalidate"))
	}
	if moreThanOneValueSet(v.Rely, v.NoRely) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Rely", "Norely"))
	}
	return errors.Join(errs...)
}

type ColumnConstraintType string

const (
	ColumnConstraintTypeUnique     ColumnConstraintType = "UNIQUE"
	ColumnConstraintTypePrimaryKey ColumnConstraintType = "PRIMARY KEY"
	ColumnConstraintTypeForeignKey ColumnConstraintType = "FOREIGN KEY"
)

type InlineForeignKey struct {
	TableName  string              `ddl:"keyword" sql:"REFERENCES"`
	ColumnName []string            `ddl:"keyword,parentheses"`
	Match      *MatchType          `ddl:"keyword" sql:"MATCH"`
	On         *ForeignKeyOnAction `ddl:"keyword" sql:"ON"`
}

// REFERENCES <ref_table_name> [ ( <ref_col_name> ) ]
// [ MATCH { FULL | SIMPLE | PARTIAL } ]
// [ ON [ UPDATE { CASCADE | SET NULL | SET DEFAULT | RESTRICT | NO ACTION } ]
// [ DELETE { CASCADE | SET NULL | SET DEFAULT | RESTRICT | NO ACTION } ] ]
// TODO [SNOW-934647]: validate
func (v *InlineForeignKey) validate() error {
	// table name required (not empty)
	// at least one column
	return nil
}

type MatchType string

var (
	FullMatchType    MatchType = "FULL"
	SimpleMatchType  MatchType = "SIMPLE"
	PartialMatchType MatchType = "PARTIAL"
)

type ForeignKeyOnAction struct {
	OnUpdate *ForeignKeyAction `ddl:"parameter,no_equals" sql:"ON UPDATE"`
	OnDelete *ForeignKeyAction `ddl:"parameter,no_equals" sql:"ON DELETE"`
}

type ForeignKeyAction string

const (
	ForeignKeyCascadeAction    ForeignKeyAction = "CASCADE"
	ForeignKeySetNullAction    ForeignKeyAction = "SET NULL"
	ForeignKeySetDefaultAction ForeignKeyAction = "SET DEFAULT"
	ForeignKeyRestrictAction   ForeignKeyAction = "RESTRICT"
	ForeignKeyNoAction         ForeignKeyAction = "NO ACTION"
)
