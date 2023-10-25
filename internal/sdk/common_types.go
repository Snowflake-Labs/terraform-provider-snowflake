// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"errors"
	"strconv"
	"time"
)

var (
	_ validatable = new(TimeTravel)
	_ validatable = new(Clone)
)

type TimeTravel struct {
	Timestamp *time.Time `ddl:"parameter,single_quotes,arrow_equals" sql:"TIMESTAMP"`
	Offset    *int       `ddl:"parameter,arrow_equals" sql:"OFFSET"`
	Statement *string    `ddl:"parameter,single_quotes,arrow_equals" sql:"STATEMENT"`
}

func (v *TimeTravel) validate() error {
	if !exactlyOneValueSet(v.Timestamp, v.Offset, v.Statement) {
		return errors.New("exactly one of TIMESTAMP, OFFSET or STATEMENT can be set")
	}
	return nil
}

type Clone struct {
	SourceObject ObjectIdentifier `ddl:"identifier" sql:"CLONE"`
	At           *TimeTravel      `ddl:"list,parentheses,no_comma" sql:"AT"`
	Before       *TimeTravel      `ddl:"list,parentheses,no_comma" sql:"BEFORE"`
}

func (v *Clone) validate() error {
	if everyValueSet(v.At, v.Before) {
		return errors.New("only one of AT or BEFORE can be set")
	}
	if valueSet(v.At) {
		return v.At.validate()
	}
	if valueSet(v.Before) {
		return v.Before.validate()
	}
	return nil
}

type LimitFrom struct {
	Rows *int    `ddl:"keyword"`
	From *string `ddl:"parameter,no_equals,single_quotes" sql:"FROM"`
}

type In struct {
	Account  *bool                    `ddl:"keyword" sql:"ACCOUNT"`
	Database AccountObjectIdentifier  `ddl:"identifier" sql:"DATABASE"`
	Schema   DatabaseObjectIdentifier `ddl:"identifier" sql:"SCHEMA"`
}

type Like struct {
	Pattern *string `ddl:"keyword,single_quotes"`
}

type TagAssociation struct {
	Name  ObjectIdentifier `ddl:"identifier"`
	Value string           `ddl:"parameter,single_quotes"`
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
	Value        *int
	DefaultValue *int
	Description  string
}

type BoolProperty struct {
	Value        bool
	DefaultValue bool
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
	var value *int
	var defaultValue *int
	v, err := strconv.Atoi(row.Value)
	if err == nil {
		value = &v
	} else {
		value = nil
	}
	dv, err := strconv.Atoi(row.DefaultValue)
	if err == nil {
		defaultValue = &dv
	} else {
		defaultValue = nil
	}
	return &IntProperty{
		Value:        value,
		DefaultValue: defaultValue,
		Description:  row.Description,
	}
}

type RowAccessPolicy struct {
	rowAccessPolicy bool                   `ddl:"static" sql:"ROW ACCESS POLICY"`
	Name            SchemaObjectIdentifier `ddl:"identifier"`
	On              []string               `ddl:"keyword,parentheses" sql:"ON"`
}

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

type ColumnConstraintType string

var (
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

func (v *InlineForeignKey) validate() error {
	return nil
}

type MatchType string

var (
	FullMatchType    MatchType = "FULL"
	SimpleMatchType  MatchType = "SIMPLE"
	PartialMatchType MatchType = "PARTIAL"
)

type ForeignKeyOnAction struct {
	OnUpdate *bool `ddl:"parameter,no_equals" sql:"ON UPDATE"`
	OnDelete *bool `ddl:"parameter,no_equals" sql:"ON DELETE"`
}

func (row *propertyRow) toBoolProperty() *BoolProperty {
	var value bool
	if row.Value != "" && row.Value != "null" {
		value = ToBool(row.Value)
	} else {
		value = false
	}
	var defaultValue bool
	if row.DefaultValue != "" && row.Value != "null" {
		defaultValue = ToBool(row.DefaultValue)
	} else {
		defaultValue = false
	}
	return &BoolProperty{
		Value:        value,
		DefaultValue: defaultValue,
		Description:  row.Description,
	}
}
