package objectparametersassert

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func AccountParametersExt(t *testing.T) *AccountParametersAssert {
	t.Helper()
	return &AccountParametersAssert{
		assert.NewSnowflakeParametersAssertWithProvider(sdk.NewAccountObjectIdentifier(""), sdk.ObjectTypeAccount, func(t *testing.T, identifier sdk.AccountObjectIdentifier) []*sdk.Parameter {
			t.Helper()
			return acc.TestClient().Parameter.ShowAccountParameters(t)
		}),
	}
}

func AccountParametersPrefetchedExt(t *testing.T, parameters []*sdk.Parameter) *AccountParametersAssert {
	t.Helper()
	return &AccountParametersAssert{
		assert.NewSnowflakeParametersAssertWithParameters(sdk.NewAccountObjectIdentifier(""), sdk.ObjectTypeAccount, parameters),
	}
}
