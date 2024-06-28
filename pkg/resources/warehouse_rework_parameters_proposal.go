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
// TODO [SNOW-1348102 - after discussion]: test (unit and acceptance)
// TODO [SNOW-1348102 - after discussion]: more readable errors
// TODO [SNOW-1348102 - after discussion]: handle different types than int
func markChangedParameters(objectParameters []sdk.ObjectParameter, currentParameters []*sdk.Parameter, d *schema.ResourceData, level sdk.ParameterType) error {
	for _, param := range objectParameters {
		currentSnowflakeParameter, err := collections.FindOne(currentParameters, func(p *sdk.Parameter) bool {
			return p.Key == string(param)
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
			intValue, err := strconv.Atoi((*currentSnowflakeParameter).Value)
			if err != nil {
				return err
			}
			if err = d.Set(strings.ToLower(string(param)), intValue); err != nil {
				return err
			}
		}
		// this handles situations in which parameter was unset from the object
		// we can just set the config value to <nil> because:
		// 1. if it was missing in config before, then no drift will be reported
		// 2. if it had a non-empty value, then the drift will be reported and the value will be set during update
		if (*currentSnowflakeParameter).Level != level {
			// TODO [SNOW-1348102 - after discussion]: this is currently set to an artificial default
			if err = d.Set(strings.ToLower(string(param)), IntDefault); err != nil {
				return err
			}
		}
	}
	return nil
}
