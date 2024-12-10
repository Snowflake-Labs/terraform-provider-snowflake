package testdatatypes

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"

var (
	DataTypeNumber_36_2, _ = datatypes.ParseDataType("NUMBER(36, 2)")
	DataTypeVarchar_100, _ = datatypes.ParseDataType("VARCHAR(100)")
	DataTypeFloat, _       = datatypes.ParseDataType("FLOAT")
	DataTypeVariant, _     = datatypes.ParseDataType("VARIANT")
)
