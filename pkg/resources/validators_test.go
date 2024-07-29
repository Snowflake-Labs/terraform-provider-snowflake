package resources

import (
	"fmt"
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
	externalObjectIdentifierCheck := IsValidIdentifier[sdk.ExternalObjectIdentifier]()

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
			Name:       "validation: incorrect form for external object identifier - less parts than expected",
			Value:      "a",
			Error:      `unexpected number of parts 1 in identifier a, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`,
			CheckingFn: externalObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for external object identifier - more parts than expected",
			Value:      "a.b.c.d",
			Error:      `unexpected number of parts 4 in identifier a.b.c.d, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`,
			CheckingFn: externalObjectIdentifierCheck,
		},
		{
			Name:       "validation: incorrect form for account object identifier - multiple parts",
			Value:      "a.b",
			Error:      `<name>, but was <database_name>.<name>`,
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "correct form for account object identifier",
			Value:      "a",
			CheckingFn: accountObjectIdentifierCheck,
		},
		{
			Name:       "correct form for account object identifier - quoted",
			Value:      "\"a.b\"",
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
		{
			Name:       "correct form for external object identifier",
			Value:      "org.acc.db",
			CheckingFn: externalObjectIdentifierCheck,
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
		{
			Name:     "external object identifier from generic parameter",
			Expected: "<organization_name>.<account_name>.<external_object_name>",
			Actual:   getExpectedIdentifierRepresentationFromGeneric[sdk.ExternalObjectIdentifier](),
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
		{
			Name:              "external object identifier from generic parameter",
			Expected:          "<organization_name>.<account_name>.<external_object_name>",
			Identifier:        sdk.ExternalObjectIdentifier{},
			IdentifierPointer: &sdk.ExternalObjectIdentifier{},
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

func Test_sdkValidation(t *testing.T) {
	genericNormalize := func(value string) (any, error) {
		if value == "ok" {
			return "ok", nil
		} else {
			return nil, fmt.Errorf("incorrect value %s", value)
		}
	}

	t.Run("valid generic normalize", func(t *testing.T) {
		valid := "ok"

		diag := sdkValidation(genericNormalize)(valid, cty.IndexStringPath("path"))

		assert.Empty(t, diag)
	})

	t.Run("invalid generic normalize", func(t *testing.T) {
		invalid := "nok"

		diag := sdkValidation(genericNormalize)(invalid, cty.IndexStringPath("path"))

		assert.Len(t, diag, 1)
		assert.Contains(t, diag[0].Summary, fmt.Sprintf("incorrect value %s", invalid))
	})

	t.Run("valid warehouse size", func(t *testing.T) {
		valid := string(sdk.WarehouseSizeSmall)

		diag := sdkValidation(sdk.ToWarehouseSize)(valid, cty.IndexStringPath("path"))

		assert.Empty(t, diag)
	})

	t.Run("invalid warehouse size", func(t *testing.T) {
		invalid := "SMALLa"

		diag := sdkValidation(sdk.ToWarehouseSize)(invalid, cty.IndexStringPath("path"))

		assert.Len(t, diag, 1)
		assert.Contains(t, diag[0].Summary, fmt.Sprintf("invalid warehouse size: %s", invalid))
	})
}

func Test_IsValidAccountIdentifier(t *testing.T) {
	testCases := []struct {
		Name  string
		Value any
		Error string
	}{
		{
			Name:  "validation: invalid value type",
			Value: 123,
			Error: "Expected schema string type, but got: int",
		},
		{
			Name:  "validation: account locator",
			Value: "ABC12345",
			Error: "Unable to parse the account identifier: ABC12345. Make sure you are using the correct form of the fully qualified account name: <organization_name>.<account_name>.",
		},
		{
			Name:  "validation: identifier too long",
			Value: "a.b.c",
			Error: "Unable to parse the account identifier: a.b.c. Make sure you are using the correct form of the fully qualified account name: <organization_name>.<account_name>.",
		},
		{
			Name:  "correct account object identifier",
			Value: "a.b",
		},
		{
			Name:  "correct account object identifier - quoted",
			Value: `"a"."b"`,
		},
		{
			Name:  "correct account object identifier - mixed quotes",
			Value: `a."b"`,
		},
		{
			Name:  "correct account object identifier - dot inside",
			Value: `a."b.c"`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			diag := IsValidAccountIdentifier()(tt.Value, cty.IndexStringPath("path"))
			if tt.Error != "" {
				assert.Len(t, diag, 1)
				assert.Contains(t, diag[0].Detail, tt.Error)
			} else {
				assert.Len(t, diag, 0)
			}
		})
	}
}

func Test_isNotEqualTo(t *testing.T) {
	testCases := []struct {
		Name             string
		Value            any
		NotExpectedValue string
		ExpectedError    string
	}{
		{
			Name:             "nil value",
			Value:            nil,
			NotExpectedValue: "123",
		},
		{
			Name:             "int value",
			Value:            123,
			NotExpectedValue: "123",
			ExpectedError:    "isNotEqualTo validator: expected string type, got int",
		},
		{
			Name:             "value equal to invalid one",
			Value:            "123",
			NotExpectedValue: "123",
			ExpectedError:    "invalid value (123) set for a field [{{} path}]. error message.",
		},
		{
			Name:             "value equal to valid one",
			Value:            "456",
			NotExpectedValue: "123",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			diag := isNotEqualTo(tt.NotExpectedValue, "error message.")(tt.Value, cty.GetAttrPath("path"))
			if tt.ExpectedError != "" {
				assert.Len(t, diag, 1)
				if diag[0].Detail != "" {
					assert.Contains(t, diag[0].Detail, tt.ExpectedError)
				} else {
					assert.Contains(t, diag[0].Summary, tt.ExpectedError)
				}
			} else {
				assert.Len(t, diag, 0)
			}
		})
	}
}
