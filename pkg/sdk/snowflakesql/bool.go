package snowflakesql

import "database/sql/driver"

// NullBool is inspired by sql.NullBool, but it will handle `"null"` passed as value, too
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the [Scanner] interface.
func (n *NullBool) Scan(value any) error {
	switch value := value.(type) {
	case nil: // untyped nil
		n.Bool, n.Valid = false, false
		return nil
	case string:
		return n.fromString(value)
	case *string:
		if value == nil {
			n.Bool, n.Valid = false, false
			return nil
		}
		return n.fromString(*value)
	}

	return n.convertAny(value)
}

func (n *NullBool) fromString(value string) error {
	if value == "null" {
		// Sadly, we have to do this, as Snowflake can return `"null"` for boolean fields.
		// E.g., `disabled` field in `SHOW USERS` output.
		n.Bool, n.Valid = false, false
		return nil
	}
	return n.convertAny(value)
}

func (n *NullBool) convertAny(value any) error {
	res, err := driver.Bool.ConvertValue(value)
	if err != nil {
		n.Bool, n.Valid = false, false
		return err
	}

	n.Bool, n.Valid = res.(bool)
	return nil
}

// Value implements the [driver.Valuer] interface.
func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}
