package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestIsDataType(t *testing.T) {
	isDataType := IsDataType()
	key := "tag"

	testCases := []struct {
		Name  string
		Value any
		Error string
	}{
		{
			Name:  "validation: correct DataType value",
			Value: "NUMBER",
		},
		{
			Name:  "validation: correct DataType value in lowercase",
			Value: "number",
		},
		{
			Name:  "validation: incorrect DataType value",
			Value: "invalid data type",
			Error: "expected tag to be one of",
		},
		{
			Name:  "validation: incorrect value type",
			Value: 123,
			Error: "expected type of tag to be string",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			_, errors := isDataType(tt.Value, key)
			if tt.Error != "" {
				assert.Len(t, errors, 1)
				assert.ErrorContains(t, errors[0], tt.Error)
			} else {
				assert.Len(t, errors, 0)
			}
		})
	}
}

func TestIsValidIdentifier(t *testing.T) {
	accountObjectIdentifierCheck := IsValidIdentifier[sdk.AccountObjectIdentifier]()
	databaseObjectIdentifierCheck := IsValidIdentifier[sdk.DatabaseObjectIdentifier]()
	schemaObjectIdentifierCheck := IsValidIdentifier[sdk.SchemaObjectIdentifier]()
	tableColumnIdentifierCheck := IsValidIdentifier[sdk.TableColumnIdentifier]()

	testCases := []struct {
		Name       string
		Value      any
		Error      string
		CheckingFn schema.SchemaValidateDiagFunc
	}{
		{
			Name:       "validation: invalid value type",
			Value:      123,
			Error:      "Expected schema string type, but got: int",
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "validation: invalid identifier representation",
			Value:      "",
			Error:      "Unable to parse the identifier: ",
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for account object identifier",
			Value:      "a.b",
			Error:      "<name>, but was <database_name>.<name>",
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for database object identifier",
			Value:      "a.b.c",
			Error:      "<database_name>.<name>, but was <database_name>.<schema_name>.<name>",
			CheckingFn: databaseObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for schema object identifier",
			Value:      "a.b.c.d",
			Error:      "<database_name>.<schema_name>.<name>, but was <database_name>.<schema_name>.<table_name>.<column_name>",
			CheckingFn: schemaObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for table column identifier",
			Value:      "a",
			Error:      "<database_name>.<schema_name>.<table_name>.<column_name>, but was <name>",
			CheckingFn: tableColumnIdentifierCheck,
		},
		{
			Name:       "correct form for account object identifier",
			Value:      "a",
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "correct form for database object identifier",
			Value:      "a.b",
			CheckingFn: databaseObjectIdentifierCheck,
		},
		{
			Name:       "correct form for schema object identifier",
			Value:      "a.b.c",
			CheckingFn: schemaObjectIdentifierCheck,
		},
		{
			Name:       "correct form for table column identifier",
			Value:      "a.b.c.d",
			CheckingFn: tableColumnIdentifierCheck,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			diag := tt.CheckingFn(tt.Value, cty.IndexStringPath("path"))
			if tt.Error != "" {
				assert.Len(t, diag, 1)
				assert.Contains(t, diag[0].Detail, tt.Error)
			} else {
				assert.Len(t, diag, 0)
			}
		})
	}
}

func TestGetExpectedIdentifierFormGeneric(t *testing.T) {
	testCases := []struct {
		Name     string
		Expected string
		Actual   string
	}{
		{
			Name:     "correct account object identifier from generic parameter",
			Expected: "<name>",
			Actual:   getExpectedIdentifierRepresentationFromGeneric[sdk.AccountObjectIdentifier](),
		},
		{
			Name:     "correct database object identifier from generic parameter",
			Expected: "<database_name>.<name>",
			Actual:   getExpectedIdentifierRepresentationFromGeneric[sdk.DatabaseObjectIdentifier](),
		},
		{
			Name:     "correct schema object identifier from generic parameter",
			Expected: "<database_name>.<schema_name>.<name>",
			Actual:   getExpectedIdentifierRepresentationFromGeneric[sdk.SchemaObjectIdentifier](),
		},
		{
			Name:     "correct table column identifier from generic parameter",
			Expected: "<database_name>.<schema_name>.<table_name>.<column_name>",
			Actual:   getExpectedIdentifierRepresentationFromGeneric[sdk.TableColumnIdentifier](),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Actual)
		})
	}
}

func TestGetExpectedIdentifierFormParam(t *testing.T) {
	testCases := []struct {
		Name              string
		Expected          string
		Identifier        sdk.ObjectIdentifier
		IdentifierPointer sdk.ObjectIdentifier
	}{
		{
			Name:              "correct account object identifier from function argument",
			Expected:          "<name>",
			Identifier:        sdk.AccountObjectIdentifier{},
			IdentifierPointer: &sdk.AccountObjectIdentifier{},
		},
		{
			Name:              "correct database object identifier from function argument",
			Expected:          "<database_name>.<name>",
			Identifier:        sdk.DatabaseObjectIdentifier{},
			IdentifierPointer: &sdk.DatabaseObjectIdentifier{},
		},
		{
			Name:              "correct schema object identifier from function argument",
			Expected:          "<database_name>.<schema_name>.<name>",
			Identifier:        sdk.SchemaObjectIdentifier{},
			IdentifierPointer: &sdk.SchemaObjectIdentifier{},
		},
		{
			Name:              "correct table column identifier from function argument",
			Expected:          "<database_name>.<schema_name>.<table_name>.<column_name>",
			Identifier:        sdk.TableColumnIdentifier{},
			IdentifierPointer: &sdk.TableColumnIdentifier{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name+" - non-pointer", func(t *testing.T) {
			assert.Equal(t, tt.Expected, getExpectedIdentifierRepresentationFromParam(tt.Identifier))
		})

		t.Run(tt.Name+" - pointer", func(t *testing.T) {
			assert.Equal(t, tt.Expected, getExpectedIdentifierRepresentationFromParam(tt.IdentifierPointer))
		})
	}
}
