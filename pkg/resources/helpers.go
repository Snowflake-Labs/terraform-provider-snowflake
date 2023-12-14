package resources

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func isOk(_ interface{}, ok bool) bool {
	return ok
}

func dataTypeValidateFunc(val interface{}, _ string) (warns []string, errs []error) {
	if ok := sdk.IsValidDataType(val.(string)); !ok {
		errs = append(errs, fmt.Errorf("%v is not a valid data type", val))
	}
	return
}

func dataTypeDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	oldDT, err := sdk.ToDataType(old)
	if err != nil {
		return false
	}
	newDT, err := sdk.ToDataType(new)
	if err != nil {
		return false
	}
	return oldDT == newDT
}

func ignoreTrimSpaceSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func setIntProperty(d *schema.ResourceData, key string, property *sdk.IntProperty) error {
	if property != nil && property.Value != nil {
		if err := d.Set(key, *property.Value); err != nil {
			return err
		}
	}
	return nil
}

func setStringProperty(d *schema.ResourceData, key string, property *sdk.StringProperty) error {
	if property != nil {
		if err := d.Set(key, property.Value); err != nil {
			return err
		}
	}
	return nil
}

func setBoolProperty(d *schema.ResourceData, key string, property *sdk.BoolProperty) error {
	if property != nil {
		if err := d.Set(key, property.Value); err != nil {
			return err
		}
	}
	return nil
}

func getTagObjectIdentifier(v map[string]any) sdk.ObjectIdentifier {
	if _, ok := v["database"]; ok {
		if _, ok := v["schema"]; ok {
			return sdk.NewSchemaObjectIdentifier(v["database"].(string), v["schema"].(string), v["name"].(string))
		}
		return sdk.NewDatabaseObjectIdentifier(v["database"].(string), v["name"].(string))
	}
	return sdk.NewAccountObjectIdentifier(v["name"].(string))
}

func getPropertyTags(d *schema.ResourceData, key string) []sdk.TagAssociation {
	if from, ok := d.GetOk(key); ok {
		tags := from.([]any)
		to := make([]sdk.TagAssociation, len(tags))
		for i, t := range tags {
			v := t.(map[string]any)
			to[i] = sdk.TagAssociation{
				Name:  getTagObjectIdentifier(v),
				Value: v["value"].(string),
			}
		}
		return to
	}
	return nil
}

func GetTagsDiff(d *schema.ResourceData, key string) (unsetTags []sdk.ObjectIdentifier, setTags []sdk.TagAssociation) {
	o, n := d.GetChange(key)
	removed, added, changed := getTags(o).diffs(getTags(n))

	unsetTags = make([]sdk.ObjectIdentifier, len(removed))
	for i, t := range removed {
		unsetTags[i] = sdk.NewDatabaseObjectIdentifier(t.database, t.name)
	}

	setTags = make([]sdk.TagAssociation, len(added)+len(changed))
	for i, t := range added {
		setTags[i] = sdk.TagAssociation{
			Name:  sdk.NewSchemaObjectIdentifier(t.database, t.schema, t.name),
			Value: t.value,
		}
	}
	for i, t := range changed {
		setTags[len(added)+i] = sdk.TagAssociation{
			Name:  sdk.NewSchemaObjectIdentifier(t.database, t.schema, t.name),
			Value: t.value,
		}
	}

	return unsetTags, setTags
}

func GetPropertyAsPointer[T any](d *schema.ResourceData, property string) *T {
	value, ok := d.GetOk(property)
	if !ok {
		return nil
	}
	typedValue, ok := value.(T)
	if !ok {
		return nil
	}
	return &typedValue
}

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
// Typically, we expect passed identifiers to be a variation of sdk.ObjectIdentifier. To use this function, pass it as
// a validation function on identifier field with generic parameter set to the desired sdk.ObjectIdentifier implementation.
func IsValidIdentifier[T sdk.ObjectIdentifier]() schema.SchemaValidateDiagFunc {
	return func(value any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		// For now, we won't support sdk.ExternalObjectIdentifiers. The reason behind it is that the functions that parse identifiers are not
		// able to differentiate between sdk.ExternalObjectIdentifiers and sdk.DatabaseObjectIdentifier or sdk.SchemaObjectIdentifier,
		// because sdk.ExternalObjectIdentifiers has varying parts count (2 or 3).
		if _, ok := any(sdk.ExternalObjectIdentifier{}).(T); ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid schema identifier type",
				Detail:        "Identifier validation is not available for sdk.ExternalObjectIdentifier type. This is a provider error please file a report: https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/new/choose",
				AttributePath: path,
			})
			return diags
		}

		if stringValue, ok := value.(string); ok {
			id, err := helpers.DecodeSnowflakeParameterID(stringValue)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to parse the identifier",
					Detail: fmt.Sprintf(
						"Unable to parse the identifier: %s. Make sure you are using the correct form of the fully qualified name for this field: %s",
						stringValue,
						getExpectedIdentifierForm[T](nil),
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
						getExpectedIdentifierForm[T](nil),
						getExpectedIdentifierForm(&id),
					),
					AttributePath: path,
				})
			}
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid schema identifier type",
				Detail:        fmt.Sprintf("Expected schema string type, but got: %T. This is a provider error please file a report: https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/new/choose", value),
				AttributePath: path,
			})
		}

		return diags
	}
}

// getExpectedIdentifierForm will choose the type either from the objectIdentifier parameter if it's present. If it's not,
// then it will create a new identifier based on the generic type parameter T, then it will return the proper structure
// we are expecting for the given sdk.ObjectIdentifier type.
func getExpectedIdentifierForm[T sdk.ObjectIdentifier](objectIdentifier *T) string {
	if objectIdentifier != nil {
		return (*objectIdentifier).Representation()
	}
	id := new(T)
	return sdk.GetIdentifierRepresentation(*id)
}
