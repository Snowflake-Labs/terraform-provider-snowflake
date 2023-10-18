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

type Error struct { //nolint:all
	file         string
	line         int
	message      string
	nestedErrors []error
}

func (e *Error) Error() string {
	builder := new(strings.Builder)
	writeTree(e, builder, 0, true)
	return builder.String()
}

// NewError Creates new sdk.Error with information like filename or line number (depending on where NewError was called)
func NewError(message string) error {
	return newSDKError(message, 2)
}

// JoinErrors returns an error that wraps the given errors.
// Any nil error values are discarded.
// JoinErrors returns nil if errs contains no non-nil values, otherwise returns sdk.Error with nested errors
func JoinErrors(errs ...error) error {
	notNilErrs := make([]error, 0)
	for _, err := range errs {
		if err != nil {
			notNilErrs = append(notNilErrs, err)
		}
	}
	if len(notNilErrs) == 0 {
		return nil
	}
	err := newSDKError("joined error", 2)
	err.nestedErrors = notNilErrs
	return err
}

func newSDKError(message string, skip int, nested ...error) *Error {
	line, filename := getCallerInfo(skip)
	return &Error{
		file:         filename,
		line:         line,
		message:      message,
		nestedErrors: nested,
	}
}

func getCallerInfo(skip int) (int, string) {
	_, file, line, _ := runtime.Caller(skip + 1)
	fileSplit := strings.Split(file, "/")
	var filename string
	if len(fileSplit) > 1 {
		filename = fileSplit[len(fileSplit)-1]
	} else {
		filename = fileSplit[0]
	}
	return line, filename
}

func writeTree(e error, builder *strings.Builder, indent int, fileInfo bool) {
	var sdkErr *Error
	if joinedErr, ok := e.(interface{ Unwrap() []error }); ok { //nolint:all
		errs := joinedErr.Unwrap()
		for i, err := range errs {
			if i > 0 {
				builder.WriteByte('\n')
			}
			writeTree(err, builder, indent, fileInfo)
		}
	} else if errors.As(e, &sdkErr) {
		if fileInfo {
			builder.WriteString(strings.Repeat("› ", indent) + fmt.Sprintf("[%s:%d] %s", sdkErr.file, sdkErr.line, sdkErr.message))
		} else {
			builder.WriteString(strings.Repeat("› ", indent) + fmt.Sprintf("[file:line] %s", sdkErr.message))
		}
		for _, err := range sdkErr.nestedErrors {
			builder.WriteByte('\n')
			writeTree(err, builder, indent+2, fileInfo)
		}
	} else {
		builder.WriteString(strings.Repeat("› ", indent) + e.Error())
	}
}
