package sdk

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	ErrNilOptions                    = errors.New("options cannot be nil")
	ErrPatternRequiredForLikeKeyword = errors.New("pattern must be specified for like keyword")

	// go-snowflake errors.
	ErrObjectNotExistOrAuthorized = errors.New("object does not exist or not authorized")
	ErrAccountIsEmpty             = errors.New("account is empty")

	// snowflake-sdk errors.
	ErrInvalidObjectIdentifier = errors.New("invalid object identifier")
	ErrDifferentDatabase       = errors.New("database must be the same")
)

func errOneOf(structName string, fieldNames ...string) error {
	return fmt.Errorf("%v fields: %v are incompatible and cannot be set at the same time", structName, fieldNames)
}

func errNotSet(structName string, fieldNames ...string) error {
	return fmt.Errorf("%v fields: %v should be set", structName, fieldNames)
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
		"does not exist or not authorized": ErrObjectNotExistOrAuthorized,
		"account is empty":                 ErrAccountIsEmpty,
	}
	for k, v := range m {
		if strings.Contains(err.Error(), k) {
			return v
		}
	}

	return err
}
