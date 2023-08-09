package sdk

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	errNilOptions                    = errors.New("options cannot be nil")
	errPatternRequiredForLikeKeyword = errors.New("pattern must be specified for like keyword")

	// go-snowflake errors.
	errObjectNotExistOrAuthorized = errors.New("object does not exist or not authorized")
	errAccountIsEmpty             = errors.New("account is empty")

	// snowflake-sdk errors.
	errInvalidObjectIdentifier = errors.New("invalid object identifier")
)

func errOneOf(structName string, fieldNames ...string) error {
	return fmt.Errorf("%v fields: %v are incompatible and cannot be set at once", structName, fieldNames)
}

func errNotSet(structName string, fieldName string) error {
	return fmt.Errorf("%v field: %v should be set", structName, fieldName)
}

func errExactlyOneOf(fieldNames ...string) error {
	return fmt.Errorf("exactly one of %v must be set", fieldNames)
}

func errAtLeastOneOf(fieldNames ...string) error {
	return fmt.Errorf("at least one of %v must be set", fieldNames)
}

func decodeDriverError(err error) error {
	if err == nil {
		return nil
	}
	log.Printf("[DEBUG] err: %v\n", err)
	m := map[string]error{
		"does not exist or not authorized": errObjectNotExistOrAuthorized,
		"account is empty":                 errAccountIsEmpty,
	}
	for k, v := range m {
		if strings.Contains(err.Error(), k) {
			return v
		}
	}

	return err
}
