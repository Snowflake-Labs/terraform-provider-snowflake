package sdk

import (
	"errors"
	"fmt"
)

func printErr(err error) {
	if e, ok := err.(*SDKError); ok {
		fmt.Println(e.errorFileInfoHidden())
	} else if fn, ok := err.(SDKPredefinedError); ok {
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
	//[file:line] Higher level error
	//› › [file:line] Some error
	//› › [file:line] Some predefined error
	//
	//[file:line] Higher level error
	//› › [file:line] Some error
	//› › [file:line] Some predefined error
	//
	//[file:line] Root
	//› › [file:line] Higher level error
	//› › › › [file:line] Some error
	//› › › › [file:line] Some predefined error
	//› › [file:line] Higher level error
	//› › › › [file:line] Some error
	//› › › › [file:line] Some predefined error
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

	// Output:
	//
}
