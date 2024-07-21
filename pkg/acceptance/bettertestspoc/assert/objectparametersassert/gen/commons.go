package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type SnowflakeObjectParameters struct {
	Name       string
	IdType     string
	Level      sdk.ParameterType
	Parameters []SnowflakeParameter
}

func (p SnowflakeObjectParameters) ObjectName() string {
	return p.Name
}

type SnowflakeParameter struct {
	ParameterName string
	ParameterType string
	DefaultValue  string
	DefaultLevel  string
}
