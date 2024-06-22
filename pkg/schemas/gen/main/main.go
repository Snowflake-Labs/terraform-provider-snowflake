//go:build exclude

package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var showStructs = []any{
	sdk.Account{},
	sdk.Alert{},
	sdk.ApiIntegration{},
	sdk.ApplicationPackage{},
	sdk.ApplicationRole{},
	sdk.Application{},
	sdk.DatabaseRole{},
	sdk.Database{},
	sdk.DynamicTable{},
	sdk.EventTable{},
	sdk.ExternalFunction{},
	sdk.ExternalTable{},
	sdk.FailoverGroup{},
	sdk.FileFormat{},
	sdk.Function{},
	sdk.Grant{},
	sdk.ManagedAccount{},
	sdk.MaskingPolicy{},
	sdk.MaterializedView{},
	sdk.NetworkPolicy{},
	sdk.NetworkRule{},
	sdk.NotificationIntegration{},
	sdk.Parameter{},
	sdk.PasswordPolicy{},
	sdk.Pipe{},
	sdk.PolicyReference{},
	sdk.Procedure{},
	sdk.ReplicationAccount{},
	sdk.ReplicationDatabase{},
	sdk.Region{},
	sdk.ResourceMonitor{},
	sdk.Role{},
	sdk.RowAccessPolicy{},
	sdk.Schema{},
	sdk.SecurityIntegration{},
	sdk.Sequence{},
	sdk.SessionPolicy{},
	sdk.Share{},
	sdk.Stage{},
	sdk.StorageIntegration{},
	sdk.Streamlit{},
	sdk.Stream{},
	sdk.Table{},
	sdk.Tag{},
	sdk.Task{},
	sdk.User{},
	sdk.View{},
	sdk.Warehouse{},
}

func main() {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	for _, s := range showStructs {
		printFields(s)
	}
}

// TODO: test completely new struct with:
//   - basic type fields (string, int, float, bool)
//   - pointer to basic types fields
//   - time.Time, *time.Time
//   - enum (string and int)
//   - slice (string, enum)
//   - identifier (each type)
//   - slice (identifier)
//   - (?) slice of pointers
//   - (?) pointer to slice
//   - (?) struct
func printFields(s any) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	fmt.Println("===========================")
	fmt.Printf("%s\n", v.Type().String())
	fmt.Println("===========================")

	for i := 0; i < v.NumField(); i++ {
		currentField := v.Field(i)
		currentName := v.Type().Field(i).Name
		currentType := v.Type().Field(i).Type.String()
		//currentValue := currentField.Interface()

		var kind reflect.Kind
		var isPtr bool

		if currentField.Kind() == reflect.Pointer {
			isPtr = true
			kind = currentField.Type().Elem().Kind()
		} else {
			kind = currentField.Kind()
		}

		var underlyingType string
		if isPtr {
			underlyingType = "*"
		}
		underlyingType += kind.String()

		gen.TabularOutput(40, currentName, currentType, underlyingType)
	}
	fmt.Println()
}
