package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1348101 - next PR]: rename all these methods to setFromXxxProperty
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

func setStringPropertyIfNotEmpty(d *schema.ResourceData, key string, property *sdk.StringProperty) error {
	if property != nil && property.Value != "" {
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

func setBooleanStringFromBoolProperty(d *schema.ResourceData, key string, property *sdk.BoolProperty) error {
	if property != nil {
		if err := d.Set(key, booleanStringFromBool(property.Value)); err != nil {
			return err
		}
	}
	return nil
}
