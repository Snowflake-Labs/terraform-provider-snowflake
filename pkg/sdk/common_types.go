package sdk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	var errs []error
	if everyValueSet(v.At, v.Before) {
		errs = append(errs, errors.New("only one of AT or BEFORE can be set"))
	}
	if valueSet(v.At) {
		errs = append(errs, v.At.validate())
	}
	if valueSet(v.Before) {
		errs = append(errs, v.Before.validate())
	}
	return errors.Join(errs...)
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

type ExtendedIn struct {
	In
	Application        AccountObjectIdentifier `ddl:"identifier" sql:"APPLICATION"`
	ApplicationPackage AccountObjectIdentifier `ddl:"identifier" sql:"APPLICATION PACKAGE"`
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

// Format in database is `(column <data_type>)`
// TODO(SNOW-1596962): Fully support VECTOR data type
// TODO(SNOW-1660588): Use ParseFunctionArgumentsFromString
func ParseTableColumnSignature(signature string) ([]TableColumnSignature, error) {
	plainSignature := strings.ReplaceAll(signature, "(", "")
	plainSignature = strings.ReplaceAll(plainSignature, ")", "")
	signatureParts := strings.Split(plainSignature, ", ")
	arguments := make([]TableColumnSignature, len(signatureParts))

	for i, elem := range signatureParts {
		parts := strings.Split(elem, " ")
		if len(parts) < 2 {
			return []TableColumnSignature{}, fmt.Errorf("expected argument name and type, got %s", elem)
		}
		dataType, err := ToDataType(parts[len(parts)-1])
		if err != nil {
			return []TableColumnSignature{}, err
		}
		arguments[i] = TableColumnSignature{
			Name: strings.Join(parts[:len(parts)-1], " "),
			Type: dataType,
		}
	}
	return arguments, nil
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

type FloatProperty struct {
	Value        *float64
	DefaultValue *float64
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

func (row *propertyRow) toFloatProperty() *FloatProperty {
	var value *float64
	var defaultValue *float64
	v, err := strconv.ParseFloat(row.Value, 64)
	if err == nil {
		value = &v
	} else {
		value = nil
	}
	dv, err := strconv.ParseFloat(row.DefaultValue, 64)
	if err == nil {
		defaultValue = &dv
	} else {
		defaultValue = nil
	}
	return &FloatProperty{
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

type ReturnResultsBehavior string

var (
	ReturnResultsBehaviorVolatile  ReturnResultsBehavior = "VOLATILE"
	ReturnResultsBehaviorImmutable ReturnResultsBehavior = "IMMUTABLE"
)

func ReturnResultsBehaviorPointer(v ReturnResultsBehavior) *ReturnResultsBehavior {
	return &v
}

type ReturnNullValues string

var (
	ReturnNullValuesNull    ReturnNullValues = "NULL"
	ReturnNullValuesNotNull ReturnNullValues = "NOT NULL"
)

func ReturnNullValuesPointer(v ReturnNullValues) *ReturnNullValues {
	return &v
}

type Secret struct {
	VariableName string `ddl:"keyword,single_quotes"`
	Name         string `ddl:"parameter,no_quotes"`
}

type ValuesBehavior string

var (
	ValuesBehaviorOrder   ValuesBehavior = "ORDER"
	ValuesBehaviorNoOrder ValuesBehavior = "NOORDER"
)

func ValuesBehaviorPointer(v ValuesBehavior) *ValuesBehavior {
	return &v
}

type Distribution string

var (
	DistributionInternal Distribution = "INTERNAL"
	DistributionExternal Distribution = "EXTERNAL"
)

func DistributionPointer(v Distribution) *Distribution {
	return &v
}

type LogLevel string

const (
	LogLevelTrace LogLevel = "TRACE"
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelFatal LogLevel = "FATAL"
	LogLevelOff   LogLevel = "OFF"
)

func ToLogLevel(value string) (LogLevel, error) {
	switch strings.ToUpper(value) {
	case string(LogLevelTrace):
		return LogLevelTrace, nil
	case string(LogLevelDebug):
		return LogLevelDebug, nil
	case string(LogLevelInfo):
		return LogLevelInfo, nil
	case string(LogLevelWarn):
		return LogLevelWarn, nil
	case string(LogLevelError):
		return LogLevelError, nil
	case string(LogLevelFatal):
		return LogLevelFatal, nil
	case string(LogLevelOff):
		return LogLevelOff, nil
	default:
		return "", fmt.Errorf("unknown log level: %s", value)
	}
}

var AllLogLevels = []LogLevel{
	LogLevelTrace,
	LogLevelDebug,
	LogLevelInfo,
	LogLevelWarn,
	LogLevelError,
	LogLevelFatal,
	LogLevelOff,
}

type TraceLevel string

const (
	TraceLevelAlways  TraceLevel = "ALWAYS"
	TraceLevelOnEvent TraceLevel = "ON_EVENT"
	TraceLevelOff     TraceLevel = "OFF"
)

func ToTraceLevel(value string) (TraceLevel, error) {
	switch strings.ToUpper(value) {
	case string(TraceLevelAlways):
		return TraceLevelAlways, nil
	case string(TraceLevelOnEvent):
		return TraceLevelOnEvent, nil
	case string(TraceLevelOff):
		return TraceLevelOff, nil
	default:
		return "", fmt.Errorf("unknown trace level: %s", value)
	}
}

var AllTraceLevels = []TraceLevel{
	TraceLevelAlways,
	TraceLevelOnEvent,
	TraceLevelOff,
}

// StringAllowEmpty is a wrapper on string to allow using empty strings in SQL.
type StringAllowEmpty struct {
	Value string `ddl:"keyword,single_quotes"`
}
