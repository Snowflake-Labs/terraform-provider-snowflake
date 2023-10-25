// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ParamType string

var (
	Integer    ParamType = "int"
	String     ParamType = "string"
	NullString ParamType = "nullstring"
	StringList ParamType = "stringlist"
	Bool       ParamType = "bool"
)

type SQLParameter struct {
	structName string
	sqlName    string
	paramType  ParamType
}

type SQLBuilderConfig struct {
	beforeObjectType map[string]string
	afterObjectType  map[string]string
	parameters       []*SQLParameter
}

type SQLBuilder struct {
	objectType       string
	objectTypePlural string
	createConfig     *SQLBuilderConfig
	alterConfig      *SQLBuilderConfig
	unsetConfig      *SQLBuilderConfig
	dropConfig       *SQLBuilderConfig
	readOutputConfig *SQLBuilderConfig
}

type KeywordPosition = string

const (
	BeforeObjectType KeywordPosition = "beforeObjectType"
	AfterObjectType  KeywordPosition = "afterObjectType"
	PosParameter     KeywordPosition = "parameter"
)

func parseConfigFromType(t reflect.Type) (*SQLBuilderConfig, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %v is not a struct", t.Name())
	}

	config := &SQLBuilderConfig{
		beforeObjectType: map[string]string{},
		afterObjectType:  map[string]string{},
		parameters:       []*SQLParameter{},
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
			config.parameters = append(config.parameters, subconf.parameters...)

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
			case reflect.TypeOf(sql.NullString{}):
				paramType = NullString
			case reflect.SliceOf(reflect.TypeOf("")):
				paramType = StringList
			case reflect.TypeOf(true):
				paramType = Bool
			default:
				switch f.Type.Kind() {
				case reflect.String:
					paramType = String
				default:
					return nil, fmt.Errorf("unsupported field type: %v", f.Type)
				}
			}

			config.parameters = append(config.parameters, &SQLParameter{
				structName: f.Name,
				sqlName:    f.Tag.Get("db"),
				paramType:  paramType,
			})
		}
	}

	return config, nil
}

func newSQLBuilder(
	objectType string,
	objectTypePlural string,
	createInputType, alterInputType, unsetInputType, dropInputType, readOutputType reflect.Type,
) (*SQLBuilder, error) {
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

	return &SQLBuilder{
		objectType:       objectType,
		objectTypePlural: objectTypePlural,
		createConfig:     createConfig,
		alterConfig:      alterConfig,
		unsetConfig:      unsetConfig,
		dropConfig:       dropConfig,
		readOutputConfig: readOutputConfig,
	}, nil
}

func (b *SQLBuilder) renderKeywords(obj Identifier, kwConf map[string]string) (string, error) {
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

func (b *SQLBuilder) renderParameters(obj Identifier, paramConf []*SQLParameter, withValues bool) (string, error) {
	sb := strings.Builder{}

	for i := range paramConf {
		param := paramConf[i]

		rv, err := getFieldValue(obj, param.structName)
		if err != nil {
			return "", err
		}

		switch paramConf[i].paramType {
		case Bool:
			ok, err := getFieldValue(obj, param.structName+"Ok")
			if err != nil {
				return "", err
			}
			if ok.Bool() {
				if withValues {
					sb.WriteString(fmt.Sprintf(` %v = %t`, param.sqlName, rv.Bool()))
				} else {
					sb.WriteString(fmt.Sprintf(` %v`, param.sqlName))
				}
			}
		case Integer:
			ok, err := getFieldValue(obj, param.structName+"Ok")
			if err != nil {
				return "", err
			}
			if ok.Bool() {
				if withValues {
					sb.WriteString(fmt.Sprintf(` %v = %v`, param.sqlName, rv.Int()))
				} else {
					sb.WriteString(fmt.Sprintf(` %v`, param.sqlName))
				}
			}
		case String:
			ok, err := getFieldValue(obj, param.structName+"Ok")
			if err != nil {
				return "", err
			}
			if ok.Bool() {
				if withValues {
					sb.WriteString(fmt.Sprintf(` %v = '%v'`, param.sqlName, rv.String()))
				} else {
					sb.WriteString(fmt.Sprintf(` %v`, param.sqlName))
				}
			}
		case NullString:
			ok, err := getFieldValue(obj, param.structName+"Ok")
			if err != nil {
				return "", err
			}
			if ok.Bool() {
				if withValues {
					ns, ok := rv.Interface().(sql.NullString)
					if !ok {
						return "", fmt.Errorf("Cannot convert %v to NullString", rv)
					}
					sb.WriteString(fmt.Sprintf(` %v = '%v'`, param.sqlName, ns.String))
				} else {
					sb.WriteString(fmt.Sprintf(` %v`, param.sqlName))
				}
			}
		case StringList:
			ok, err := getFieldValue(obj, param.structName+"Ok")
			if err != nil {
				return "", err
			}
			if ok.Bool() {
				if withValues {
					slice, _ := rv.Interface().([]string)
					sb.WriteString(fmt.Sprintf(` %v = ('%v')`, param.sqlName, strings.Join(slice, "', '")))
				} else {
					sb.WriteString(fmt.Sprintf(` %v`, param.sqlName))
				}
			}
		}
	}

	return sb.String(), nil
}

func (b *SQLBuilder) Create(obj Identifier) (string, error) {
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
	sb.WriteString(fmt.Sprintf(` "%v"`, (obj).QualifiedName()))

	// eg. `PARAM = "value"`
	params, err := b.renderParameters(obj, b.createConfig.parameters, true)
	if err != nil {
		return "", err
	}
	sb.WriteString(params)

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *SQLBuilder) Alter(obj Identifier) (string, error) {
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
	sb.WriteString(fmt.Sprintf(` "%v"`, obj.QualifiedName()))

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

func (b *SQLBuilder) Unset(obj Identifier) (string, error) {
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
	sb.WriteString(fmt.Sprintf(` "%v"`, obj.QualifiedName()))

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

func (b *SQLBuilder) Drop(obj Identifier) (string, error) {
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
	sb.WriteString(fmt.Sprintf(` "%v"`, obj.QualifiedName()))

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *SQLBuilder) ShowLike(obj Identifier) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("SHOW")

	// eg. "TABLES"
	sb.WriteString(fmt.Sprintf(" %v", b.objectTypePlural))

	sb.WriteString(" LIKE")

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` '%v'`, obj.QualifiedName()))

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *SQLBuilder) Describe(obj Identifier) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("DESCRIBE")

	// eg. "TABLE"
	sb.WriteString(fmt.Sprintf(" %v", b.objectType))

	// eg. "my_table"
	sb.WriteString(fmt.Sprintf(` "%v"`, obj.QualifiedName()))

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *SQLBuilder) ParseDescribe(rows *sql.Rows, obj Identifier) error {
	var property, value, defaultValue, devnull string

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	vars := []interface{}{}
	for i := range cols {
		switch cols[i] {
		case "property":
			vars = append(vars, &property)
		case "value", "property_value":
			vars = append(vars, &value)
		case "default", "property_default":
			vars = append(vars, &defaultValue)
		default:
			vars = append(vars, &devnull)
		}
	}

	for rows.Next() {
		if err := rows.Scan(vars...); err != nil {
			return err
		}

		for i := range b.readOutputConfig.parameters {
			param := b.readOutputConfig.parameters[i]

			if param.sqlName == property {
				err := setFieldValue(obj, param.structName, value)
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
	case reflect.Bool:
		curField.SetBool(strings.ToLower(value) == "true")
	case reflect.Slice:
		if len(value) == 0 {
			curField.Set(reflect.ValueOf(make([]string, 0)))
		} else if matched, err := regexp.MatchString(`^\[('.+'(,\s*'.+')*)?\]$`, value); err == nil && matched { // parse object types
			trimmed := strings.Trim(value, "[]")
			split := strings.Split(trimmed, ",")
			values := make([]string, 0)
			for i := range split {
				values = append(values, strings.Trim(strings.Trim(split[i], " "), "'"))
			}
			curField.Set(reflect.ValueOf(values))
		} else {
			split := strings.Split(value, ",")
			curField.Set(reflect.ValueOf(split))
		}
	}
	return nil
}
