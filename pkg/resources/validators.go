package resources

import (
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func IsDataType() schema.SchemaValidateFunc { //nolint:staticcheck
	return func(value any, key string) (warnings []string, errors []error) {
		stringValue, ok := value.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string, got %T", key, value))
			return warnings, errors
		}

		_, err := sdk.ToDataType(stringValue)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected %s to be one of %T values, got %s", key, sdk.DataTypeString, stringValue))
		}

		return warnings, errors
	}
}

// IsValidIdentifier is a validator that can be used for validating identifiers passed in resources and data sources.
//
// Typically, we expect passed identifiers to be a variation of sdk.ObjectIdentifier.
// For now, we're expecting implementations of sdk.ObjectIdentifier, because we won't support sdk.ExternalObjectIdentifiers.
// The reason behind it is that the functions that parse identifiers are not able to differentiate between
// sdk.ExternalObjectIdentifiers and sdk.DatabaseObjectIdentifier or sdk.SchemaObjectIdentifier.
// That's because sdk.ExternalObjectIdentifiers has varying parts count (2 or 3).
//
// To use this function, pass it as a validation function on identifier field with generic parameter set to the desired sdk.ObjectIdentifier implementation.
func IsValidIdentifier[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier]() schema.SchemaValidateDiagFunc {
	return func(value any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		if _, ok := value.(string); !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid schema identifier type",
				Detail:        fmt.Sprintf("Expected schema string type, but got: %T. This is a provider error please file a report: https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/new/choose", value),
				AttributePath: path,
			})
			return diags
		}

		stringValue := value.(string)
		id, err := helpers.DecodeSnowflakeParameterID(stringValue)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to parse the identifier",
				Detail: fmt.Sprintf(
					"Unable to parse the identifier: %s. Make sure you are using the correct form of the fully qualified name for this field: %s.\nOriginal Error: %s",
					stringValue,
					getExpectedIdentifierRepresentationFromGeneric[T](),
					err.Error(),
				),
				AttributePath: path,
			})
			return diags
		}

		if _, ok := id.(T); !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid identifier type",
				Detail: fmt.Sprintf(
					"Expected %s identifier type, but got: %T. The correct form of the fully qualified name for this field is: %s, but was %s",
					reflect.TypeOf(new(T)).Elem().Name(),
					id,
					getExpectedIdentifierRepresentationFromGeneric[T](),
					getExpectedIdentifierRepresentationFromParam(id),
				),
				AttributePath: path,
			})
		}

		return diags
	}
}

func getExpectedIdentifierRepresentationFromGeneric[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier]() string {
	return getExpectedIdentifierForm(new(T))
}

func getExpectedIdentifierRepresentationFromParam(id sdk.ObjectIdentifier) string {
	return getExpectedIdentifierForm(id)
}

func getExpectedIdentifierForm(id any) string {
	switch id.(type) {
	case sdk.AccountObjectIdentifier, *sdk.AccountObjectIdentifier:
		return "<name>"
	case sdk.DatabaseObjectIdentifier, *sdk.DatabaseObjectIdentifier:
		return "<database_name>.<name>"
	case sdk.SchemaObjectIdentifier, *sdk.SchemaObjectIdentifier:
		return "<database_name>.<schema_name>.<name>"
	case sdk.TableColumnIdentifier, *sdk.TableColumnIdentifier:
		return "<database_name>.<schema_name>.<table_name>.<column_name>"
	}
	return ""
}
