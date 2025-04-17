package resources

import (
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
	return dataType, nil
}

// sqlNew is temporary as not all the data types has the temporary method implemented
func sqlNew(dt datatypes.DataType) string {
	switch v := dt.(type) {
	case *datatypes.NumberDataType:
		return v.ToSqlNew()
	case *datatypes.TextDataType:
		return v.ToSqlNew()
	default:
		return v.ToSql()
	}
}
