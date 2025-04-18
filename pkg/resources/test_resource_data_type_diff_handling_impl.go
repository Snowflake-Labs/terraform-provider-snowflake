package resources

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// contents of this file will be used as common functions if approved

func readDatatypeCommon(d *schema.ResourceData, key string) (datatypes.DataType, error) {
	dataTypeRawConfig := d.Get(key).(string)
	dataType, err := datatypes.ParseDataType(dataTypeRawConfig)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] correctly parsed data type %v", dataType)
	return dataType, nil
}

func readChangedDatatypeCommon(d *schema.ResourceData, key string) (datatypes.DataType, error) {
	_, n := d.GetChange(key)
	dataType, err := datatypes.ParseDataType(n.(string))
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] correctly parsed data type %v", dataType)
	return dataType, nil
}

// AsFullyKnown is temporary as not all the data types has the temporary method implemented
func AsFullyKnown(dt datatypes.DataType) datatypes.DataType {
	switch v := dt.(type) {
	case *datatypes.NumberDataType:
		return v.AsFullyKnown()
	case *datatypes.TextDataType:
		return v.AsFullyKnown()
	default:
		return v
	}
}
