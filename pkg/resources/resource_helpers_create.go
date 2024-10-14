package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringAttributeCreate(d *schema.ResourceData, key string, createField **string) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.String(v.(string))
	}
	return nil
}

func intAttributeCreate(d *schema.ResourceData, key string, createField **int) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Int(v.(int))
	}
	return nil
}

func intAttributeWithSpecialDefaultCreate(d *schema.ResourceData, key string, createField **int) error {
	if v := d.Get(key).(int); v != IntDefault {
		*createField = sdk.Int(v)
	}
	return nil
}

func booleanStringAttributeCreate(d *schema.ResourceData, key string, createField **bool) error {
	if v := d.Get(key).(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return err
		}
		*createField = sdk.Bool(parsed)
	}
	return nil
}

func accountObjectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.AccountObjectIdentifier) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	return nil
}

func objectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.ObjectIdentifier) error {
	if v, ok := d.GetOk(key); ok {
		objectIdentifier, err := sdk.ParseObjectIdentifierString(v.(string))
		if err != nil {
			return err
		}
		*createField = sdk.Pointer(objectIdentifier)
	}
	return nil
}

func attributeDirectValueCreate[T any](d *schema.ResourceData, key string, createField **T, value *T) error {
	if _, ok := d.GetOk(key); ok {
		*createField = value
	}
	return nil
}

func copyGrantsAttributeCreate(d *schema.ResourceData, isOrReplace bool, orReplaceField, copyGrantsField **bool) error {
	if isOrReplace {
		*orReplaceField = sdk.Bool(true)
		if d.Get("copy_grants").(bool) {
			*copyGrantsField = sdk.Bool(true)
		}
	}
	return nil
}
