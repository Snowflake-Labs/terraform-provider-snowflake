package objectassert

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO [SNOW-1501905]: this file should be fully regenerated when adding and option to assert the results of describe
type FunctionDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.FunctionDetails, sdk.SchemaObjectIdentifierWithArguments]
}

func FunctionDetails(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) *FunctionDetailsAssert {
	t.Helper()
	return &FunctionDetailsAssert{
		assert.NewSnowflakeObjectAssertWithProvider(sdk.ObjectType("FUNCTION_DETAILS"), id, acc.TestClient().Function.DescribeDetails),
	}
}

func (f *FunctionDetailsAssert) HasSignature(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Signature != expected {
			return fmt.Errorf("expected signature: %v; got: %v", expected, o.Signature)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasReturns(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Returns != expected {
			return fmt.Errorf("expected returns: %v; got: %v", expected, o.Returns)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasLanguage(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Language != expected {
			return fmt.Errorf("expected language: %v; got: %v", expected, o.Language)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasBody(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Body == nil {
			return fmt.Errorf("expected body to have value; got: nil")
		}
		if *o.Body != expected {
			return fmt.Errorf("expected body: %v; got: %v", expected, *o.Body)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasNullHandling(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.NullHandling == nil {
			return fmt.Errorf("expected null handling to have value; got: nil")
		}
		if *o.NullHandling != expected {
			return fmt.Errorf("expected null handling: %v; got: %v", expected, *o.NullHandling)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasVolatility(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Volatility == nil {
			return fmt.Errorf("expected volatility to have value; got: nil")
		}
		if *o.Volatility != expected {
			return fmt.Errorf("expected volatility: %v; got: %v", expected, *o.Volatility)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasExternalAccessIntegrations(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.ExternalAccessIntegrations == nil {
			return fmt.Errorf("expected external access integrations to have value; got: nil")
		}
		if *o.ExternalAccessIntegrations != expected {
			return fmt.Errorf("expected external access integrations: %v; got: %v", expected, *o.ExternalAccessIntegrations)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasSecrets(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Secrets == nil {
			return fmt.Errorf("expected secrets to have value; got: nil")
		}
		if *o.Secrets != expected {
			return fmt.Errorf("expected secrets: %v; got: %v", expected, *o.Secrets)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasImports(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Imports == nil {
			return fmt.Errorf("expected imports to have value; got: nil")
		}
		if *o.Imports != expected {
			return fmt.Errorf("expected imports: %v; got: %v", expected, *o.Imports)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasHandler(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Handler == nil {
			return fmt.Errorf("expected handler to have value; got: nil")
		}
		if *o.Handler != expected {
			return fmt.Errorf("expected handler: %v; got: %v", expected, *o.Handler)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasRuntimeVersion(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.RuntimeVersion == nil {
			return fmt.Errorf("expected runtime version to have value; got: nil")
		}
		if *o.RuntimeVersion != expected {
			return fmt.Errorf("expected runtime version: %v; got: %v", expected, *o.RuntimeVersion)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasPackages(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Packages == nil {
			return fmt.Errorf("expected packages to have value; got: nil")
		}
		if *o.Packages != expected {
			return fmt.Errorf("expected packages: %v; got: %v", expected, *o.Packages)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasTargetPath(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.TargetPath == nil {
			return fmt.Errorf("expected target path to have value; got: nil")
		}
		if *o.TargetPath != expected {
			return fmt.Errorf("expected target path: %v; got: %v", expected, *o.TargetPath)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasInstalledPackages(expected string) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.InstalledPackages == nil {
			return fmt.Errorf("expected installed packages to have value; got: nil")
		}
		if *o.InstalledPackages != expected {
			return fmt.Errorf("expected installed packages: %v; got: %v", expected, *o.InstalledPackages)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasIsAggregate(expected bool) *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.IsAggregate == nil {
			return fmt.Errorf("expected is aggregate to have value; got: nil")
		}
		if *o.IsAggregate != expected {
			return fmt.Errorf("expected is aggregate: %v; got: %v", expected, *o.IsAggregate)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasBodyNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Body != nil {
			return fmt.Errorf("expected body to be nil, was %v", *o.Body)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasNullHandlingNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.NullHandling != nil {
			return fmt.Errorf("expected null handling to be nil, was %v", *o.NullHandling)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasVolatilityNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Volatility != nil {
			return fmt.Errorf("expected volatility to be nil, was %v", *o.Volatility)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasExternalAccessIntegrationsNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.ExternalAccessIntegrations != nil {
			return fmt.Errorf("expected external access integrations to be nil, was %v", *o.ExternalAccessIntegrations)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasSecretsNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Secrets != nil {
			return fmt.Errorf("expected secrets to be nil, was %v", *o.Secrets)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasImportsNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Imports != nil {
			return fmt.Errorf("expected imports to be nil, was %v", *o.Imports)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasHandlerNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Handler != nil {
			return fmt.Errorf("expected handler to be nil, was %v", *o.Handler)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasRuntimeVersionNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.RuntimeVersion != nil {
			return fmt.Errorf("expected runtime version to be nil, was %v", *o.RuntimeVersion)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasPackagesNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.Packages != nil {
			return fmt.Errorf("expected packages to be nil, was %v", *o.Packages)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasTargetPathNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.TargetPath != nil {
			return fmt.Errorf("expected target path to be nil, was %v", *o.TargetPath)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasInstalledPackagesNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.InstalledPackages != nil {
			return fmt.Errorf("expected installed packages to be nil, was %v", *o.InstalledPackages)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasIsAggregateNil() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.IsAggregate != nil {
			return fmt.Errorf("expected is aggregate to be nil, was %v", *o.IsAggregate)
		}
		return nil
	})
	return f
}

func (f *FunctionDetailsAssert) HasInstalledPackagesNotEmpty() *FunctionDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FunctionDetails) error {
		t.Helper()
		if o.InstalledPackages == nil {
			return fmt.Errorf("expected installed packages to not be nil")
		}
		if *o.InstalledPackages == "" {
			return fmt.Errorf("expected installed packages to not be empty")
		}
		return nil
	})
	return f
}
