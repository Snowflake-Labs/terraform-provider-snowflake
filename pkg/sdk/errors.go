package sdk

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

// TODO change to predefined errors
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

const ghIssueBodyTemplate = `
**Provider Version**

The provider version you are using.

**Terraform Version**

The version of Terraform you were using when the bug was encountered.

**Describe the bug**

A clear and concise description of what the bug is.

**Expected behavior**

A clear and concise description of what you expected to happen.

**Code samples and commands**

Please add code examples and commands that were run to cause the problem.

**Additional context**

Add any other context about the problem here.

<!-- 
Please provide additional error messages if we missed any (see the errors below and compare it with your console output)
-->

Errors:
%s
`

// NewTopLevelError wraps an error with "final" error message of sdk.
// It should be placed in the highest place of call stack to catch as much error context as possible.
// TODO Is it possible to call this function somewhere high in the callstack, so it always be top level ? If not maybe this would be overkill / too hard to use
// TODO It's possible to wrap errors multiple times, this function will try to keep only the one with most error context.
func NewTopLevelError(err error) error {
	// TODO if called multiple times, unwrap lower level errors and wrap err to have more context

	return fmt.Errorf(`
Snowflake Terraform Provider error!
If you think you've encountered a bug, please report it with the link below.
If any of the error information is missing in the issue body, please fill it up.
Any additional information (or context what you were trying to achieve) would be helpful
to provide the solution or fix as soon as possible. Thanks :)

https://github.com//Snowflake-Labs/terraform-provider-snowflake/issues/new?labels=bug&title=%s&body=%s

Errors:
%s`,
		url.QueryEscape("New issue"),
		url.QueryEscape(fmt.Sprintf(ghIssueBodyTemplate, err.Error())),
		err.Error(),
	)
}

// TODO: error distinction (go-snowflake, validation, etc.)
// 	-
// 	-

type SDKError struct { //nolint:all
	file         string
	line         int
	message      string
	nestedErrors []error
}

func (e *SDKError) errorFileInfoHidden() string {
	builder := new(strings.Builder)
	writeTree(e, builder, 0, false)
	return builder.String()
}

func (e *SDKError) Error() string {
	builder := new(strings.Builder)
	writeTree(e, builder, 0, true)
	return builder.String()
}

// NewError Creates new sdk.SDKError with information like filename or line number
// depending on where NewError was called
func NewError(message string) error {
	return newSDKError(message, 2)
}

// TODO: new error with error as a parameter

// NewPredefinedError Lets you predefine factory method for given sdk.SDKError which is convenient
// when given error must be returned multiple times
func NewPredefinedError(message string) func() error {
	return func() error {
		return newSDKError(message, 2)
	}
}

type SDKPredefinedError func() error //nolint:all

func (fn SDKPredefinedError) Error() string {
	return fn().Error()
}

// NewPredefinedError2 same as NewPredefinedError, but we're getting filename + line where this error is declared,
// not returned, which may be worse than the NewPredefinedError
// version, but this one is better at replacing current predefined errors like ErrInvalidObjectIdentifier. With NewPredefinedError
// we would have to change every return statement of this error to function call (because NewPredefinedError returns factory method for predefined err)
// and the nature of NewPredefinedError2 is that we can return value as with ErrInvalidObjectIdentifier example.
// The only disadvantage of this is that we'll get file + line from declaration which may cause issues that we have right now
// (not knowing places where errors have been returned from - one function can return ErrInvalidObjectIdentifier in two places
// and we'll know from which place with this approach unfortunately) and this error should solve them.
// We can use both (e.g. for new errors use NewPredefinedError2 and for replacing old use NewPredefinedError)
// or we can just choose one approach over the other. (of course NewPredefinedError2 will be renamed to something that makes more sense)
// TODO this approach to predefined errors, this is lazy evaluated error, so the only filename + line number
// TODO if we're going with this approach, change other methods to handle SDKPredefinedError type
func NewPredefinedError2(message string) SDKPredefinedError {
	err := newSDKError(message, 2)
	return func() error {
		return err
	}
}

func NewErrorOneOf(structPointer any, fields ...any) error {
	structure := reflect.ValueOf(structPointer).Elem()
	var fieldsBuilder strings.Builder

	for fieldIndex := range fields {
		fieldMeta := structure.Type().Field(fieldIndex)
		fieldValue := structure.Field(fieldIndex)
		fieldPointer := reflect.NewAt(fieldMeta.Type, unsafe.Pointer(fieldValue.UnsafeAddr()))
		if fieldMeta.Type.Kind() == reflect.Pointer {
			fieldPointer = fieldPointer.Elem()
		}

		fieldsBuilder.WriteString(fmt.Sprintf("%s %s(%v)", fieldMeta.Name, fieldMeta.Type, fieldPointer.Elem().Interface()))

		if fieldIndex != len(fields)-1 {
			fieldsBuilder.WriteString(", ")
		}
	}

	return newSDKError(
		fmt.Sprintf(
			"fields of struct %s [%s] are incompatible and shouldn't be set at the same time",
			structure.Type().Name(),
			fieldsBuilder.String(),
		),
		2,
	)
}

func NewErrorNotSet(structPointer any, fields ...any) error {
	structure := reflect.ValueOf(structPointer).Elem()
	var fieldsBuilder strings.Builder

	for fieldIndex := range fields {
		fieldMeta := structure.Type().Field(fieldIndex)
		fieldsBuilder.WriteString(fmt.Sprintf("%s %s", fieldMeta.Name, fieldMeta.Type))

		if fieldIndex != len(fields)-1 {
			fieldsBuilder.WriteString(", ")
		}
	}

	return newSDKError(
		fmt.Sprintf(
			"fields of struct %s: [%s] are required and should be set",
			structure.Type().Name(),
			fieldsBuilder.String(),
		),
		2,
	)
}

// WrapErrors wraps errs with wrapper error.
// When wrapper is any other type than SDKError than it'll create new SDKError with wrapper's error message
// and in every case errs will be added to its internal list of errors.
// TODO We can force to use SDKError as wrapper
func WrapErrors(wrapper error, errs ...error) error {
	if err, ok := wrapper.(*SDKError); ok { //nolint:all
		err.nestedErrors = append(err.nestedErrors, errs...)
		return wrapper
	} else {
		return newSDKError(wrapper.Error(), 2, errs...)
	}
}

func newSDKError(message string, skip int, nested ...error) error {
	line, filename := getCallerInfo(skip)
	return &SDKError{
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
	var sdkErr *SDKError
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
