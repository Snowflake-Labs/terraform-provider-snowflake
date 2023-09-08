package sdk_test

import (
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func ExampleNewError() {
	fmt.Println(sdk.NewError("New SDK Error"))

	// Output:
	// [errors_test.go:10] New SDK Error
}

func ExampleNewPredefinedError() {
	predefinedErr := sdk.NewPredefinedError("Some predefined error")
	// Filename and line number will be extracted from the point of calling factory method (line 19 in this case)
	fmt.Println(predefinedErr())

	// Output:
	// [errors_test.go:19] Some predefined error
}

func ExampleWrapErrors() {
	sdkErr := sdk.NewError("Some error")
	predefinedErr := sdk.NewPredefinedError("Some predefined error")
	joinedErr := errors.Join(sdkErr, predefinedErr())

	err := sdk.WrapErrors(sdk.NewError("Higher level error"), joinedErr)
	// will be same as
	err2 := sdk.WrapErrors(sdk.NewError("Higher level error"), sdkErr, predefinedErr())

	// a more complicated structure could look something like this, see Output
	err3 := sdk.WrapErrors(sdk.NewError("Root"), err, err2)

	fmt.Printf("%v\n\n", err)
	fmt.Printf("%v\n\n", err2)
	fmt.Printf("%v\n\n", err3)

	// Output:
	// [errors_test.go:30] Higher level error
	// › › [errors_test.go:26] Some error
	// › › [errors_test.go:28] Some predefined error
	//
	// [errors_test.go:32] Higher level error
	// › › [errors_test.go:26] Some error
	// › › [errors_test.go:32] Some predefined error
	//
	// [errors_test.go:35] Root
	// › › [errors_test.go:30] Higher level error
	// › › › › [errors_test.go:26] Some error
	// › › › › [errors_test.go:28] Some predefined error
	// › › [errors_test.go:32] Higher level error
	// › › › › [errors_test.go:26] Some error
	// › › › › [errors_test.go:32] Some predefined error
}

func ExampleNewErrorOneOf() {
	type Example struct {
		FieldOne   string
		fieldTwo   int
		fieldThree *string
		FieldFour  []int
	}

	e := Example{
		FieldOne:   "str value",
		fieldTwo:   134,
		fieldThree: sdk.String("str pointer"),
		FieldFour:  []int{123, 321},
	}

	fmt.Println(sdk.NewErrorOneOf(&e, e.FieldOne, e.fieldTwo, e.fieldThree, e.FieldFour))

	// Output:
	// [errors_test.go:74] fields of struct Example [FieldOne string(str value), fieldTwo int(134), fieldThree *string(str pointer), FieldFour []int([123 321])] are incompatible and shouldn't be set at the same time
}

func ExampleNewErrorNotSet() {
	type Example struct {
		FieldOne   string
		fieldTwo   int
		fieldThree *string
		FieldFour  []int
	}

	e := Example{}

	fmt.Println(sdk.NewErrorNotSet(&e, e.FieldOne, e.fieldTwo, e.fieldThree, e.FieldFour))

	// Output:
	// [errors_test.go:90] fields of struct Example: [FieldOne string, fieldTwo int, fieldThree *string, FieldFour []int] are required and should be set
}

func ExampleNewTopLevelError() {
	sdkErr := sdk.NewError("Some error")
	predefinedErr := sdk.NewPredefinedError("Some predefined error")
	joinedErr := errors.Join(sdkErr, predefinedErr())
	err := sdk.WrapErrors(sdk.NewError("Higher level error"), joinedErr)
	err2 := sdk.WrapErrors(sdk.NewError("Higher level error"), sdkErr, predefinedErr())
	err3 := sdk.WrapErrors(sdk.NewError("Root"), err, err2)
	gigaErr := errors.Join(err, err2, err3)
	fmt.Println(sdk.NewTopLevelError(gigaErr))

	// Output:
	//
}
