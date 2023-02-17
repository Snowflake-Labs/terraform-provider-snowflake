package helpers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/xid"
)

func RandomSnowflakeID() string {
	guid := xid.New()
	return fmt.Sprintf("snow-%s", guid.String())
}

func DecodeSnowflakeImportID(id string, v interface{}) (interface{}, error) {
	attributes := make(map[string]string)
	parts := strings.Split(id, "|")
	for _, part := range parts {
		if !strings.Contains(part, "=") {
			return nil, fmt.Errorf("invalid import ID format: %s, attributes must be defined as key=value format", id)
		}
		key := strings.TrimSpace(strings.Split(part, "=")[0])
		value := strings.TrimSpace(strings.Split(part, "=")[1])
		attributes[key] = value
	}
	for k, v := range attributes {
		fmt.Printf("[DEBUG] %s=%s\n", k, v)
	}

	// w is the interface{}
	w := reflect.ValueOf(&v).Elem()

	// Allocate a temporary variable with type of the struct.
	//    v.Elem() is the value contained in the interface.
	tmp := reflect.New(w.Elem().Type()).Elem()

	// Copy the struct value contained in interface to
	// the temporary variable.
	tmp.Set(w.Elem())

	t := tmp.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tf")
		importValue := attributes[tag]
		f := tmp.FieldByName(field.Name)
		switch f.Kind() {
		case reflect.String:
			f.SetString(importValue)
		case reflect.Int:
			intVal, err := strconv.Atoi(importValue)
			if err != nil {
				return nil, err
			}
			f.SetInt(int64(intVal))
		case reflect.Bool:
			f.SetBool(importValue == "true")
		case reflect.Slice:
			p := strings.Split(importValue, ",")
			for _, v := range p {
				v := strings.Trim(v, "\"")
				f.Set(reflect.Append(f, reflect.ValueOf(v)))
			}
		}
	}
	// Set the interface to the modified struct value.
	w.Set(tmp)
	fmt.Printf("[DEBUG] v=%+v\n", v)
	return v, nil
}
