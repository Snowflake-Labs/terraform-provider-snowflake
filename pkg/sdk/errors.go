package sdk

import (
	"errors"
	"fmt"
	"log"
	"runtime"
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

type IntErrType string

const (
	IntErrEqual          IntErrType = "equal to"
	IntErrGreaterOrEqual IntErrType = "greater than or equal to"
	IntErrGreater        IntErrType = "greater than"
	IntErrLessOrEqual    IntErrType = "less than or equal to"
	IntErrLess           IntErrType = "less than"
)

func errIntValue(structName string, fieldName string, intErrType IntErrType, limit int) error {
	return fmt.Errorf("%s field: %s must be %s %d", structName, fieldName, string(intErrType), limit)
}

func errIntBetween(structName string, fieldName string, from int, to int) error {
	return fmt.Errorf("%s field: %s must be between %d and %d", structName, fieldName, from, to)
}

func errInvalidIdentifier(structName string, identifierField string) error {
	return fmt.Errorf("invalid object identifier of %s field: %s", structName, identifierField)
}

func errOneOf(structName string, fieldNames ...string) error {
	return fmt.Errorf("%v fields: %v are incompatible and cannot be set at the same time", structName, fieldNames)
}

func errNotSet(structName string, fieldNames ...string) error {
	return fmt.Errorf("%v fields: %v should be set", structName, fieldNames)
}

func errExactlyOneOf(structName string, fieldNames ...string) error {
	return fmt.Errorf("exactly one of %s fileds %v must be set", structName, fieldNames)
}

func errAtLeastOneOf(structName string, fieldNames ...string) error {
	return fmt.Errorf("at least one of %s fields %v must be set", structName, fieldNames)
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

type errorSDK struct {
	file         string
	line         int
	message      string
	nestedErrors []error
}

func NewError(message string) error {
	_, file, line, _ := runtime.Caller(1)
	return &errorSDK{
		file:         file,
		line:         line,
		message:      message,
		nestedErrors: make([]error, 0),
	}
}

func (e *errorSDK) Error() string {
	errorMessage := fmt.Sprintf("[%s: line %d] %s", e.file, e.line, e.message)
	if len(e.nestedErrors) > 0 {
		var b []byte
		for i, err := range e.nestedErrors {
			if i > 0 {
				b = append(b, '\n')
			}
			b = append(b, err.Error()...)
		}
		return fmt.Sprintf("%s\n%s", errorMessage, string(b))
	}
	return errorMessage
}
