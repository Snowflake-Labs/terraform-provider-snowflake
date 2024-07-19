package assert

import (
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO: make assertions naming consistent (resource paramaters vs snowflake parameters)
type UserParametersAssert struct {
	*SnowflakeParametersAssert[sdk.AccountObjectIdentifier]
}

func UserParameters(t *testing.T, id sdk.AccountObjectIdentifier) *UserParametersAssert {
	t.Helper()
	return &UserParametersAssert{
		NewSnowflakeParametersAssertWithProvider(id, sdk.ObjectTypeUser, acc.TestClient().Parameter.ShowUserParameters),
	}
}

func UserParametersPrefetched(t *testing.T, id sdk.AccountObjectIdentifier, parameters []*sdk.Parameter) *UserParametersAssert {
	t.Helper()
	return &UserParametersAssert{
		NewSnowflakeParametersAssertWithParameters(id, sdk.ObjectTypeUser, parameters),
	}
}

func (w *UserParametersAssert) HasAbortDetachedQuery(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterAbortDetachedQuery, strconv.FormatBool(expected)))
	return w
}

func (w *UserParametersAssert) HasBinaryInputFormat(expected sdk.BinaryInputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterBinaryInputFormat, string(expected)))
	return w
}

func (w *UserParametersAssert) HasClientMemoryLimit(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterClientMemoryLimit, strconv.Itoa(expected)))
	return w
}

func (w *UserParametersAssert) HasDateOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterDateOutputFormat, expected))
	return w
}
