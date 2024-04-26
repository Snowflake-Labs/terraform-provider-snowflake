package resources

import (
	"fmt"
	"strings"

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
		if _, ok := value.(string); !ok {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Invalid schema identifier type",
					Detail:        fmt.Sprintf("Expected schema string type, but got: %T. This is a provider error please file a report: https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/new/choose", value),
					AttributePath: path,
				},
			}
		}

		_, err := helpers.SafelyDecodeSnowflakeID[T](value.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

// StringInSlice has the same implementation as validation.StringInSlice, but adapted to schema.SchemaValidateDiagFunc
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %v to be string", path)
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.EqualFold(v, str)) {
				return nil
			}
		}

		return diag.Errorf("expected %v to be one of %q, got %s", path, valid, v)
	}
}
