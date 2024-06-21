package resources

import (
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const parametersAttributeName = "parameters"

// markChangedParameters assumes that the snowflake parameter name is mirrored in schema (as lower-cased name)
// TODO [after discussion/next PR]: test (unit and acceptance)
// TODO [after discussion/next PR]: more readable errors
// TODO [after discussion/next PR]: handle different types than int
func markChangedParameters(objectParameters []sdk.ParameterWithType, currentParameters []*sdk.Parameter, d *schema.ResourceData, level sdk.ParameterLevel) error {
	for _, param := range objectParameters {
		currentSnowflakeParameter, err := collections.FindOne(currentParameters, func(p *sdk.Parameter) bool {
			return p.Key == string(param.Key)
		})
		if err != nil {
			return err
		}
		// this handles situations in which parameter was set on object externally (so either the value or the level was changed)
		// we can just set the config value to the current Snowflake value because:
		// 1. if it did not change, then no drift will be reported
		// 2. if it had different non-empty value, then the drift will be reported and the value will be set during update
		// 3. if it had empty value, then the drift will be reported and the value will be unset during update
		if (*currentSnowflakeParameter).Level == level {
			switch (*currentSnowflakeParameter).Type {
			case sdk.ParameterTypeString:
				if err = d.Set(strings.ToLower(string(param.Key)), (*currentSnowflakeParameter).Value); err != nil {
					return err
				}
			case sdk.ParameterTypeNumber:
				intValue, err := strconv.Atoi((*currentSnowflakeParameter).Value)
				if err != nil {
					return err
				}
				if err = d.Set(strings.ToLower(string(param.Key)), intValue); err != nil {
					return err
				}
			case sdk.ParameterTypeBoolean:
				boolValue, err := strconv.ParseBool((*currentSnowflakeParameter).Value)
				if err != nil {
					return err
				}
				if err = d.Set(strings.ToLower(string(param.Key)), strconv.FormatBool(boolValue)); err != nil {
					return err
				}
			}
		}
		// this handles situations in which parameter was unset from the object
		// we can just set the config value to <nil> because:
		// 1. if it was missing in config before, then no drift will be reported
		// 2. if it had a non-empty value, then the drift will be reported and the value will be set during update
		if (*currentSnowflakeParameter).Level != level {
			// TODO [after discussion/next PR]: this is currently set to an artificial default
			switch (*currentSnowflakeParameter).Type {
			case sdk.ParameterTypeString:
				if err = d.Set(strings.ToLower(string(param.Key)), "unknown"); err != nil {
					return err
				}
			case sdk.ParameterTypeNumber:
				if err = d.Set(strings.ToLower(string(param.Key)), -1); err != nil {
					return err
				}
			case sdk.ParameterTypeBoolean:
				if err = d.Set(strings.ToLower(string(param.Key)), "unknown"); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
