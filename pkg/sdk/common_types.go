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

type ExecuteAs string

func ExecuteAsPointer(v ExecuteAs) *ExecuteAs {
	return &v
}

const (
	ExecuteAsCaller ExecuteAs = "EXECUTE AS CALLER"
	ExecuteAsOwner  ExecuteAs = "EXECUTE AS OWNER"
)

type NullInputBehavior string

func NullInputBehaviorPointer(v NullInputBehavior) *NullInputBehavior {
	return &v
}

const (
	NullInputBehaviorCalledOnNullInput NullInputBehavior = "CALLED ON NULL INPUT"
	NullInputBehaviorReturnNullInput   NullInputBehavior = "RETURN NULL ON NULL INPUT"
	NullInputBehaviorStrict            NullInputBehavior = "STRICT"
)

type Secret struct {
	VariableName string `ddl:"keyword,single_quotes"`
	Name         string `ddl:"parameter,no_quotes"`
}
