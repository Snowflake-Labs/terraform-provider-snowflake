package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"
)

// Upper case
// case a"b in a single field and in e.g. grant field where it's encoded

func TestInt_IdentifierParsing(t *testing.T) {
	testCases := []struct {
		InputName string
	}{}

	testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, sdk.NewCreateNetworkPolicyRequest())

}
