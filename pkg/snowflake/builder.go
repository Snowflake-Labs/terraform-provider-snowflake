package snowflake

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ParamType string

var (
	Integer    ParamType = "int"
	String     ParamType = "string"
	StringList ParamType = "stringlist"
)

type Param struct {
	name      string
	paramType ParamType
}

type BuilderConfig struct {
	beforeEntityType map[string]string
	afterEntityType  map[string]string
	parameters       map[string]*Param
}

type NewBuilder struct {
	entityType       string
	entityTypePlural string
	config           BuilderConfig
}

type Props interface {
	QualifiedName() string
	ID() string
}

type StatementPosition = string

const (
	BeforeEntityType StatementPosition = "beforeType"
	AfterEntityType  StatementPosition = "afterType"
	PosParameter     StatementPosition = "parameter"
)

func parseConfigFromType(t reflect.Type) (*BuilderConfig, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %v is not a struct", t.Name())
	}

	config := &BuilderConfig{
		beforeEntityType: map[string]string{},
		afterEntityType:  map[string]string{},
		parameters:       map[string]*Param{},
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		switch f.Tag.Get("pos") {
		case BeforeEntityType:
			config.beforeEntityType[f.Name] = f.Tag.Get("value")
		case AfterEntityType:
			config.afterEntityType[f.Name] = f.Tag.Get("value")
		case PosParameter:
			var paramType ParamType

			switch f.Type {
			case reflect.TypeOf(0):
				paramType = Integer
			case reflect.TypeOf(""):
				paramType = String
			case reflect.SliceOf(reflect.TypeOf("")):
				paramType = StringList
			}

			config.parameters[f.Name] = &Param{
				name:      f.Tag.Get("db"),
				paramType: paramType,
			}
		}
	}

	return config, nil
}

func newBuilder(entityType string, entityTypePlural string, t reflect.Type) (*NewBuilder, error) {
	config, err := parseConfigFromType(t)
	if err != nil {
		return nil, err
	}

	return &NewBuilder{
		entityType:       entityType,
		entityTypePlural: entityTypePlural,
		config:           *config,
	}, nil
}

func (b *NewBuilder) renderKeywords(props Props, kwConf map[string]string) (string, error) {
	sb := strings.Builder{}

	for key := range kwConf {
		ok, err := getFieldValue(props, key+"Ok")
		if err != nil {
			return "", err
		}
		val, err := getFieldValue(props, key)
		if err != nil {
			return "", err
		}
		if ok.Bool() && val.Bool() {
			sb.WriteString(fmt.Sprintf(" %v", kwConf[key]))
		}
	}

	return sb.String(), nil
}

func (b *NewBuilder) renderParameters(props Props, withValues bool) (string, error) {
	sb := strings.Builder{}

	for key := range b.config.parameters {
		name := b.config.parameters[key].name
		rv, err := getFieldValue(props, key)
		if err != nil {
			return "", err
		}

		switch b.config.parameters[key].paramType {
		case Integer:
			ok, err := getFieldValue(props, key+"Ok")
			if err != nil {
				return "", err
			}
			if ok.Bool() {
				if withValues {
					sb.WriteString(fmt.Sprintf(` %v = %v`, name, rv.Int()))
				} else {
					sb.WriteString(fmt.Sprintf(` %v`, name))
				}
			}
		case String:
			ok, err := getFieldValue(props, key+"Ok")
			if err != nil {
				return "", err
			}
			if ok.Bool() {
				if withValues {
					sb.WriteString(fmt.Sprintf(` %v = '%v'`, name, rv.String()))
				} else {
					sb.WriteString(fmt.Sprintf(` %v`, name))
				}
			}
		case StringList:
			ok, err := getFieldValue(props, key+"Ok")
			if err != nil {
				return "", err
			}
			if ok.Bool() {
				if withValues {
					slice, _ := rv.Interface().([]string)
					sb.WriteString(fmt.Sprintf(` %v = ('%v')`, name, strings.Join(slice, "', '")))
				} else {
					sb.WriteString(fmt.Sprintf(` %v`, name))
				}
			}
		}
	}

	return sb.String(), nil
}

func (b *NewBuilder) Create(props Props) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("CREATE")

	// eg. "OR REPLACE"
	before, err := b.renderKeywords(props, b.config.beforeEntityType)
	if err != nil {
		return "", err
	}
	sb.WriteString(before)

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.entityType))

	// eg. "IF NOT EXISTS"
	after, err := b.renderKeywords(props, b.config.afterEntityType)
	if err != nil {
		return "", err
	}
	sb.WriteString(after)

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, (props).QualifiedName()))

	// eg. `PARAM = "value"`
	params, err := b.renderParameters(props, true)
	if err != nil {
		return "", err
	}
	sb.WriteString(params)

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) Alter(props Props) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("ALTER")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.entityType))

	// eg. "IF EXISTS"
	after, err := b.renderKeywords(props, b.config.afterEntityType)
	if err != nil {
		return "", err
	}
	sb.WriteString(after)

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, props.QualifiedName()))

	sb.WriteString(" SET")

	// eg. `PARAM = "value"`
	params, err := b.renderParameters(props, true)
	if err != nil {
		return "", err
	}
	sb.WriteString(params)

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) Unset(props Props) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("ALTER")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.entityType))

	// eg. "IF EXISTS"
	after, err := b.renderKeywords(props, b.config.afterEntityType)
	if err != nil {
		return "", err
	}
	sb.WriteString(after)

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, props.QualifiedName()))

	sb.WriteString(" UNSET")

	// eg. `PARAM = "value"`
	params, err := b.renderParameters(props, false)
	if err != nil {
		return "", err
	}
	sb.WriteString(params)

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) Drop(props Props) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("DROP")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.entityType))

	// eg. "IF EXISTS"
	after, err := b.renderKeywords(props, b.config.afterEntityType)
	if err != nil {
		return "", err
	}
	sb.WriteString(after)

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, props.QualifiedName()))

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) Describe(props Props) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("DESCRIBE")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.entityType))

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, props.QualifiedName()))

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) ParseDescribe(rows *sql.Rows, props Props) error {
	var property, value, defaultValue, desc string

	for rows.Next() {
		if err := rows.Scan(&property, &value, &defaultValue, &desc); err != nil {
			return err
		}

		for field := range b.config.parameters {
			if b.config.parameters[field].name == property {
				err := setFieldValue(props, field, value)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

func (b *NewBuilder) Ok(_ interface{}, ok bool) bool {
	return ok
}

func getFieldValue(props Props, fieldName string) (*reflect.Value, error) {
	pointToStruct := reflect.ValueOf(props)
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can't read field value from %v because it is not a struct", props)
	}
	curField := curStruct.FieldByName(fieldName)
	if !curField.IsValid() {
		return nil, fmt.Errorf("field %v is invalid on %v", fieldName, props)
	}
	return &curField, nil
}

func setFieldValue(props Props, fieldName string, value string) error {
	pointToStruct := reflect.ValueOf(props)
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		return fmt.Errorf("can't set field value in %v because it is not a struct", props)
	}
	curField := curStruct.FieldByName(fieldName)
	if !curField.IsValid() || !curField.CanSet() {
		return fmt.Errorf("can't write to field %v on %v", fieldName, props)
	}
	switch curField.Kind() {
	case reflect.Int:
		val, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("couldn't parse %v to int", value)
		}
		curField.SetInt(int64(val))
	case reflect.String:
		if value == "null" {
			curField.SetString("")
		} else {
			curField.SetString(value)
		}
	}
	return nil
}
