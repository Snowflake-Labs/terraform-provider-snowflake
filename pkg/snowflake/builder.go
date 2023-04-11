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

type SQLParameter struct {
	name      string
	paramType ParamType
}

type BuilderConfig struct {
	beforeObjectType map[string]string
	afterObjectType  map[string]string
	parameters       map[string]*SQLParameter
}

type NewBuilder struct {
	objectType       string
	objectTypePlural string
	createConfig     *BuilderConfig
	alterConfig      *BuilderConfig
	unsetConfig      *BuilderConfig
	dropConfig       *BuilderConfig
	readOutputConfig *BuilderConfig
}

type KeywordPosition = string

const (
	BeforeObjectType KeywordPosition = "beforeObjectType"
	AfterObjectType  KeywordPosition = "afterObjectType"
	PosParameter     KeywordPosition = "parameter"
)

func parseConfigFromType(t reflect.Type) (*BuilderConfig, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %v is not a struct", t.Name())
	}

	config := &BuilderConfig{
		beforeObjectType: map[string]string{},
		afterObjectType:  map[string]string{},
		parameters:       map[string]*SQLParameter{},
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// If the field is anonymous, recursively parse the embedded value and merge into main config
		if f.Anonymous {
			subconf, err := parseConfigFromType(f.Type)
			if err != nil {
				return nil, err
			}
			for key := range subconf.beforeObjectType {
				config.beforeObjectType[key] = subconf.beforeObjectType[key]
			}
			for key := range subconf.afterObjectType {
				config.afterObjectType[key] = subconf.afterObjectType[key]
			}
			for key := range subconf.parameters {
				config.parameters[key] = subconf.parameters[key]
			}

			continue
		}

		// Otherwise try to read its "pos" tag
		switch f.Tag.Get("pos") {
		case BeforeObjectType:
			config.beforeObjectType[f.Name] = f.Tag.Get("value")
		case AfterObjectType:
			config.afterObjectType[f.Name] = f.Tag.Get("value")
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

			config.parameters[f.Name] = &SQLParameter{
				name:      f.Tag.Get("db"),
				paramType: paramType,
			}
		}
	}

	return config, nil
}

func newBuilder(
	objectType string,
	objectTypePlural string,
	createInputType, alterInputType, unsetInputType, dropInputType, readOutputType reflect.Type,
) (*NewBuilder, error) {
	createConfig, err := parseConfigFromType(createInputType)
	if err != nil {
		return nil, err
	}
	alterConfig, err := parseConfigFromType(alterInputType)
	if err != nil {
		return nil, err
	}
	unsetConfig, err := parseConfigFromType(unsetInputType)
	if err != nil {
		return nil, err
	}
	dropConfig, err := parseConfigFromType(dropInputType)
	if err != nil {
		return nil, err
	}
	readOutputConfig, err := parseConfigFromType(readOutputType)
	if err != nil {
		return nil, err
	}

	return &NewBuilder{
		objectType:       objectType,
		objectTypePlural: objectTypePlural,
		createConfig:     createConfig,
		alterConfig:      alterConfig,
		unsetConfig:      unsetConfig,
		dropConfig:       dropConfig,
		readOutputConfig: readOutputConfig,
	}, nil
}

func (b *NewBuilder) renderKeywords(obj Identifier, kwConf map[string]string) (string, error) {
	sb := strings.Builder{}

	for key := range kwConf {
		ok, err := getFieldValue(obj, key+"Ok")
		if err != nil {
			return "", err
		}
		val, err := getFieldValue(obj, key)
		if err != nil {
			return "", err
		}
		if ok.Bool() && val.Bool() {
			sb.WriteString(fmt.Sprintf(" %v", kwConf[key]))
		}
	}

	return sb.String(), nil
}

func (b *NewBuilder) renderParameters(obj Identifier, paramConf map[string]*SQLParameter, withValues bool) (string, error) {
	sb := strings.Builder{}

	for key := range paramConf {
		name := paramConf[key].name
		rv, err := getFieldValue(obj, key)
		if err != nil {
			return "", err
		}

		switch paramConf[key].paramType {
		case Integer:
			ok, err := getFieldValue(obj, key+"Ok")
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
			ok, err := getFieldValue(obj, key+"Ok")
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
			ok, err := getFieldValue(obj, key+"Ok")
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

func (b *NewBuilder) Create(obj Identifier) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("CREATE")

	// eg. "OR REPLACE"
	before, err := b.renderKeywords(obj, b.createConfig.beforeObjectType)
	if err != nil {
		return "", err
	}
	sb.WriteString(before)

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.objectType))

	// eg. "IF NOT EXISTS"
	after, err := b.renderKeywords(obj, b.createConfig.afterObjectType)
	if err != nil {
		return "", err
	}
	sb.WriteString(after)

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, (obj).QualifiedName()))

	// eg. `PARAM = "value"`
	params, err := b.renderParameters(obj, b.createConfig.parameters, true)
	if err != nil {
		return "", err
	}
	sb.WriteString(params)

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) Alter(obj Identifier) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("ALTER")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.objectType))

	// eg. "IF EXISTS"
	after, err := b.renderKeywords(obj, b.alterConfig.afterObjectType)
	if err != nil {
		return "", err
	}
	sb.WriteString(after)

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, obj.QualifiedName()))

	sb.WriteString(" SET")

	// eg. `PARAM = "value"`
	params, err := b.renderParameters(obj, b.alterConfig.parameters, true)
	if err != nil {
		return "", err
	}
	sb.WriteString(params)

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) Unset(obj Identifier) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("ALTER")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.objectType))

	// eg. "IF EXISTS"
	after, err := b.renderKeywords(obj, b.unsetConfig.afterObjectType)
	if err != nil {
		return "", err
	}
	sb.WriteString(after)

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, obj.QualifiedName()))

	sb.WriteString(" UNSET")

	// eg. `PARAM = "value"`
	params, err := b.renderParameters(obj, b.unsetConfig.parameters, false)
	if err != nil {
		return "", err
	}
	sb.WriteString(params)

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) Drop(obj Identifier) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("DROP")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.objectType))

	// eg. "IF EXISTS"
	after, err := b.renderKeywords(obj, b.dropConfig.afterObjectType)
	if err != nil {
		return "", err
	}
	sb.WriteString(after)

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, obj.QualifiedName()))

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) Describe(obj Identifier) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("DESCRIBE")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.objectType))

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` %v`, obj.QualifiedName()))

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *NewBuilder) ParseDescribe(rows *sql.Rows, obj Identifier) error {
	var property, value, defaultValue, desc string

	for rows.Next() {
		if err := rows.Scan(&property, &value, &defaultValue, &desc); err != nil {
			return err
		}

		for field := range b.readOutputConfig.parameters {
			if b.readOutputConfig.parameters[field].name == property {
				err := setFieldValue(obj, field, value)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

func getFieldValue(obj Identifier, fieldName string) (*reflect.Value, error) {
	pointToStruct := reflect.ValueOf(obj)
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can't read field value from %v because it is not a struct", obj)
	}
	curField := curStruct.FieldByName(fieldName)
	if !curField.IsValid() {
		return nil, fmt.Errorf("field %v is invalid on %v", fieldName, obj)
	}
	return &curField, nil
}

func setFieldValue(obj Identifier, fieldName string, value string) error {
	pointToStruct := reflect.ValueOf(obj)
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		return fmt.Errorf("can't set field value in %v because it is not a struct", obj)
	}
	curField := curStruct.FieldByName(fieldName)
	if !curField.IsValid() || !curField.CanSet() {
		return fmt.Errorf("can't write to field %v on %v", fieldName, obj)
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
