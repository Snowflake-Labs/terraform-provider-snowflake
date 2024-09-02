package objectassert

import (
	"fmt"
	"slices"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (w *ViewAssert) HasCreatedOnNotEmpty() *ViewAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return w
}

func (v *ViewAssert) HasNonEmptyText() *ViewAssert {
	v.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		if o.Text == "" {
			return fmt.Errorf("expected non empty text")
		}
		return nil
	})
	return v
}

func (v *ViewAssert) HasNoRowAccessPolicyReferences() *ViewAssert {
	return v.hasNoPolicyReference(sdk.PolicyKindRowAccessPolicy)
}

func (v *ViewAssert) HasNoAggregationPolicyReferences() *ViewAssert {
	return v.hasNoPolicyReference(sdk.PolicyKindAggregationPolicy)
}

func (v *ViewAssert) HasNoMaskingPolicyReferences() *ViewAssert {
	return v.hasNoPolicyReference(sdk.PolicyKindMaskingPolicy)
}

func (v *ViewAssert) HasNoProjectionPolicyReferences() *ViewAssert {
	return v.hasNoPolicyReference(sdk.PolicyKindProjectionPolicy)
}

func (v *ViewAssert) hasNoPolicyReference(kind sdk.PolicyKind) *ViewAssert {
	v.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		refs, err := acc.TestClient().PolicyReferences.GetPolicyReferences(t, o.ID(), sdk.ObjectTypeView)
		if err != nil {
			return err
		}
		refs = slices.DeleteFunc(refs, func(reference helpers.PolicyReference) bool {
			return reference.PolicyKind != string(kind)
		})
		if len(refs) > 0 {
			return fmt.Errorf("expected no %s policy references; got: %v", kind, refs)
		}
		return nil
	})
	return v
}

func (v *ViewAssert) HasRowAccessPolicyReferences(n int) *ViewAssert {
	return v.hasPolicyReference(sdk.PolicyKindRowAccessPolicy, n)
}

func (v *ViewAssert) HasAggregationPolicyReferences(n int) *ViewAssert {
	return v.hasPolicyReference(sdk.PolicyKindAggregationPolicy, n)
}

func (v *ViewAssert) HasMaskingPolicyReferences(n int) *ViewAssert {
	return v.hasPolicyReference(sdk.PolicyKindMaskingPolicy, n)
}

func (v *ViewAssert) HasProjectionPolicyReferences(n int) *ViewAssert {
	return v.hasPolicyReference(sdk.PolicyKindProjectionPolicy, n)
}

func (v *ViewAssert) hasPolicyReference(kind sdk.PolicyKind, n int) *ViewAssert {
	v.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		refs, err := acc.TestClient().PolicyReferences.GetPolicyReferences(t, o.ID(), sdk.ObjectTypeView)
		if err != nil {
			return err
		}
		refs = slices.DeleteFunc(refs, func(reference helpers.PolicyReference) bool {
			return reference.PolicyKind != string(kind)
		})
		if len(refs) != n {
			return fmt.Errorf("expected %d %s policy references; got: %d, %v", n, kind, len(refs), refs)
		}
		return nil
	})
	return v
}
