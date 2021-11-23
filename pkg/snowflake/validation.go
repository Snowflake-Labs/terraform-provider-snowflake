package snowflake

import "github.com/pkg/errors"

// ValidateIdentifier implements a strict definition of valid identifiers from
// https://docs.snowflake.net/manuals/sql-reference/identifiers-syntax.html
func ValidateIdentifier(val interface{}) (warns []string, errs []error) {
	name, ok := val.(string)
	if !ok {
		errs = append(errs, errors.Errorf("Unable to assert identifier as string type."))
		return
	}

	if len(name) == 0 {
		errs = append(errs, errors.Errorf("Identifier must be at least 1 character."))
		return
	}

	if len(name) > 256 {
		errs = append(errs, errors.Errorf("Identifier must be <= 256 characters."))
		return
	}

	// TODO handle quoted identifiers
	for k, r := range name {
		if k == 0 && !isInitialIdentifierRune(r) {
			errs = append(errs, errors.Errorf("'%s' can not start an identifier.", string(r)))
			continue
		}

		if !isIdentifierRune(r) {
			errs = append(errs, errors.Errorf("'%s' is not a valid identifier character.", string(r)))
		}
	}
	return

}

func isIdentifierRune(r rune) bool {
	return isInitialIdentifierRune(r) || r == '$' || (r >= '0' && r <= '9')
}

func isInitialIdentifierRune(r rune) bool {
	return (r == '_' ||
		(r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z'))
}
