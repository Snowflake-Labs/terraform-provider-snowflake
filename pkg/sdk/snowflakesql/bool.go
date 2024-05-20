package snowflakesql

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// Bool is inspired by sql.NullBool, but it will handle `"null"` passed as value, too
type Bool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the [Scanner] interface.
func (n *Bool) Scan(value any) error {
	switch value := value.(type) {
	case nil: // untyped nil
		n.Bool, n.Valid = false, false
		return nil
	case bool:
		return n.fromBool(&value)
	case *bool:
		return n.fromBool(value)
	case string:
		return n.fromString(&value)
	case *string:
		return n.fromString(value)
	default:
		return n.convertAny(value)
	}
}

func (n *Bool) fromBool(value *bool) error {
	if n.Valid = value != nil; n.Valid {
		n.Bool = *value
	} else {
		n.Bool = false
	}
	return nil
}

func (n *Bool) fromString(value *string) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}

	str := *value
	if str == "null" {
		// Sadly, we have to do this, as Snowflake can return `"null"` for boolean fields.
		// E.g., `disabled` field in `SHOW USERS` output.
		n.Bool, n.Valid = false, false
		return nil
	}

	return n.convertAny(str)
}

func (n *Bool) convertAny(value any) error {
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			// nil pointer to some value
			n.Bool, n.Valid = false, false
			return nil
		}
		v = v.Elem()
	}

	if !v.CanInterface() {
		// shouldn't be here, but fail without panic
		n.Bool, n.Valid = false, false
		return fmt.Errorf("can't convert %v (%T) into bool", value, value)
	}

	res, err := driver.Bool.ConvertValue(v.Interface())
	if err != nil {
		n.Bool, n.Valid = false, false
		return err
	}

	n.Bool, n.Valid = res.(bool)
	return nil
}

// Value implements the [driver.Valuer] interface.
func (n Bool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

// BoolValue returns either the default bool (false) if the Bool.Valid != true, of the underlying Bool.Value.
func (n Bool) BoolValue() bool {
	return n.Valid && n.Bool
}
