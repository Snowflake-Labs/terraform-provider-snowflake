package resources

import (
	"fmt"
	"strings"

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

// constructWithFallbacks constructs sdk.SchemaObjectIdentifier by accepting sdk.ObjectIdentifier interface and figuring out which id components
// are already there and which one we could replace with default values (databaseName, schemaName).
func constructWithFallbacks(databaseName string, schemaName string, identifier sdk.ObjectIdentifier) sdk.SchemaObjectIdentifier {
	switch id := identifier.(type) {
	case sdk.AccountObjectIdentifier:
		return sdk.NewSchemaObjectIdentifier(databaseName, schemaName, id.Name())
	case sdk.DatabaseObjectIdentifier:
		return sdk.NewSchemaObjectIdentifier(id.DatabaseName(), schemaName, id.Name())
	case sdk.SchemaObjectIdentifier:
		return id
	default:
		return sdk.NewSchemaObjectIdentifier(databaseName, schemaName, id.Name())
	}
}
