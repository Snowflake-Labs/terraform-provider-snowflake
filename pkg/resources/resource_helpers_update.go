package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func intAttributeUpdate(d *schema.ResourceData, key string, setField **int, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = sdk.Int(v.(int))
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func stringAttributeUpdate(d *schema.ResourceData, key string, setField **string, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = sdk.String(v.(string))
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func booleanStringAttributeUpdate(d *schema.ResourceData, key string, setField **bool, unsetField **bool) error {
	if d.HasChange(key) {
		if v := d.Get(key).(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return err
			}
			*setField = sdk.Bool(parsed)
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func booleanStringAttributeUnsetFallbackUpdate(d *schema.ResourceData, key string, setField **bool, fallbackValue bool) error {
	if d.HasChange(key) {
		if v := d.Get(key).(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return err
			}
			*setField = sdk.Bool(parsed)
		} else {
			*setField = sdk.Bool(fallbackValue)
		}
	}
	return nil
}

// TODO: NewAccountObjectIdentifier versus NewAccountObjectIdentifierFromFullyQualifiedName versus one of the new functions?
func accountObjectIdentifierAttributeUpdate(d *schema.ResourceData, key string, setField **sdk.AccountObjectIdentifier, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			*setField = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

// TODO: DecodeSnowflakeParameterID versus one of the new functions?
func objectIdentifierAttributeUpdate(d *schema.ResourceData, key string, setField **sdk.ObjectIdentifier, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			objectIdentifier, err := helpers.DecodeSnowflakeParameterID(v.(string))
			if err != nil {
				return err
			}
			*setField = sdk.Pointer(objectIdentifier)
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}
