package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ProcedureAssert) HasCreatedOnNotEmpty() *ProcedureAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Procedure) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected create_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *ProcedureAssert) HasExternalAccessIntegrationsNil() *ProcedureAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Procedure) error {
		t.Helper()
		if o.ExternalAccessIntegrations != nil {
			return fmt.Errorf("expected external_access_integrations to be nil but was: %v", *o.ExternalAccessIntegrations)
		}
		return nil
	})
	return a
}
