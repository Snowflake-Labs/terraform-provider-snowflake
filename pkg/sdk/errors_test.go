package sdk

import (
	"errors"
	"fmt"
)

func printErr(err error) {
	if e, ok := err.(*SDKError); ok { //nolint:all
		fmt.Println(e.errorFileInfoHidden())
	} else if fn, ok := err.(SDKPredefinedError); ok { //nolint:all
		printErr(fn())
	} else {
		fmt.Println(err)
	}
}

func ExampleNewError() {
	printErr(NewError("New SDK SDKError"))

	// Output:
	// [file:line] New SDK SDKError
}

func ExampleNewPredefinedError() {
	predefinedErr := NewPredefinedError("Some predefined error")
	errorFunc := func() error {
		return predefinedErr()
	}

	// Filename and line number will be extracted from the point of calling factory method (line: return predefinedErr() in this case)
	printErr(errorFunc())

	// Output:
	// [file:line] Some predefined error
}

func ExampleNewPredefinedError2() {
	predefinedErr := NewPredefinedError2("Some predefined error")
	errorFunc := func() error {
		return predefinedErr
	}

	// Filename and line number will be extracted from the point of declaring error (line: predefinedErr := sdk.NewPredefinedError2("Some predefined error") in this case)
	printErr(errorFunc())

	// Output:
	// [file:line] Some predefined error
}

func ExampleWrapErrors() {
	sdkErr := NewError("Some error")
	predefinedErr := NewPredefinedError("Some predefined error")
	joinedErr := errors.Join(sdkErr, predefinedErr())

	err := WrapErrors(NewError("Higher level error"), joinedErr)
	// will be same as
	err2 := WrapErrors(NewError("Higher level error"), sdkErr, predefinedErr())

	// a more complicated structure could look something like this, see Output
	err3 := WrapErrors(NewError("Root"), err, err2)

	printErr(err)
	fmt.Println()
	printErr(err2)
	fmt.Println()
	printErr(err3)

	// Output:
	// [file:line] Higher level error
	// › › [file:line] Some error
	// › › [file:line] Some predefined error
	//
	// [file:line] Higher level error
	// › › [file:line] Some error
	// › › [file:line] Some predefined error
	//
	// [file:line] Root
	// › › [file:line] Higher level error
	// › › › › [file:line] Some error
	// › › › › [file:line] Some predefined error
	// › › [file:line] Higher level error
	// › › › › [file:line] Some error
	// › › › › [file:line] Some predefined error
}

func ExampleNewErrorOneOf() {
	type Example struct {
		ExportedField          string
		unexportedField        int
		unexportedFieldPointer *string
		ExportedFieldSlice     []int
	}

	e := Example{
		ExportedField:          "str value",
		unexportedField:        134,
		unexportedFieldPointer: String("str pointer"),
		ExportedFieldSlice:     []int{123, 321},
	}

	printErr(NewErrorOneOf(&e, e.ExportedField, e.unexportedField, e.unexportedFieldPointer, e.ExportedFieldSlice))

	// Output:
	// [file:line] fields of struct Example [ExportedField string(str value), unexportedField int(134), unexportedFieldPointer *string(str pointer), ExportedFieldSlice []int([123 321])] are incompatible and shouldn't be set at the same time
}

func ExampleNewErrorNotSet() {
	type Example struct {
		ExportedField          string
		unexportedField        int
		unexportedFieldPointer *string
		ExportedFieldSlice     []int
	}

	e := Example{}

	printErr(NewErrorNotSet(&e, e.ExportedField, e.unexportedField, e.unexportedFieldPointer, e.ExportedFieldSlice))

	// Output:
	// [file:line] fields of struct Example: [ExportedField string, unexportedField int, unexportedFieldPointer *string, ExportedFieldSlice []int] are required and should be set
}

func ExampleNewTopLevelError() {
	sdkErr := NewError("Some error")
	predefinedErr := NewPredefinedError("Some predefined error")
	joinedErr := errors.Join(sdkErr, predefinedErr())
	err := WrapErrors(NewError("Higher level error"), joinedErr)
	err2 := WrapErrors(NewError("Higher level error"), sdkErr, predefinedErr())
	err3 := WrapErrors(NewError("Root"), err, err2)
	gigaErr := errors.Join(err, err2, err3)
	printErr(NewTopLevelError(gigaErr))

	// Prints:
	// Snowflake Terraform Provider error!
	// If you think you've encountered a bug, please report it with the link below.
	// If any of the error information is missing in the issue body, please fill it up.
	// Any additional information (or context what you were trying to achieve) would be helpful
	// to provide the solution or fix as soon as possible. Thanks :)
	//
	// https://github.com//Snowflake-Labs/terraform-provider-snowflake/issues/new?labels=bug&title=New+issue&body=%0A%3C%21--+%0A%2A%2AProvider+Version%2A%2A%0A%0AThe+provider+version+you+are+using.%0A%0A%2A%2ATerraform+Version%2A%2A%0A%0AThe+version+of+Terraform+you+were+using+when+the+bug+was+encountered.%0A%0A%2A%2ADescribe+the+bug%2A%2A%0A%0AA+clear+and+concise+description+of+what+the+bug+is.%0A%0A%2A%2AExpected+behavior%2A%2A%0A%0AA+clear+and+concise+description+of+what+you+expected+to+happen.%0A%0A%2A%2ACode+samples+and+commands%2A%2A%0A%0APlease+add+code+examples+and+commands+that+were+run+to+cause+the+problem.%0A%0A%2A%2AAdditional+context%2A%2A%0A%0AAdd+any+other+context+about+the+problem+here.%0A%0A%3C%21--+%0APlease+provide+additional+error+messages+if+we+missed+any+%28see+the+errors+below+and+compare+it+with+your+console+output%29%0A--%3E%0A%0AErrors%3A%0A%5Berrors_test.go%3A128%5D+Higher+level+error%0A%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A125%5D+Some+error%0A%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A127%5D+Some+predefined+error%0A%5Berrors_test.go%3A129%5D+Higher+level+error%0A%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A125%5D+Some+error%0A%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A129%5D+Some+predefined+error%0A%5Berrors_test.go%3A130%5D+Root%0A%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A128%5D+Higher+level+error%0A%E2%80%BA+%E2%80%BA+%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A125%5D+Some+error%0A%E2%80%BA+%E2%80%BA+%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A127%5D+Some+predefined+error%0A%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A129%5D+Higher+level+error%0A%E2%80%BA+%E2%80%BA+%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A125%5D+Some+error%0A%E2%80%BA+%E2%80%BA+%E2%80%BA+%E2%80%BA+%5Berrors_test.go%3A129%5D+Some+predefined+error%0A
	//
	// Errors:
	// [errors_test.go:128] Higher level error
	// › › [errors_test.go:125] Some error
	// › › [errors_test.go:127] Some predefined error
	// [errors_test.go:129] Higher level error
	// › › [errors_test.go:125] Some error
	// › › [errors_test.go:129] Some predefined error
	// [errors_test.go:130] Root
	// › › [errors_test.go:128] Higher level error
	// › › › › [errors_test.go:125] Some error
	// › › › › [errors_test.go:127] Some predefined error
	// › › [errors_test.go:129] Higher level error
	// › › › › [errors_test.go:125] Some error
	// › › › › [errors_test.go:129] Some predefined error
}
