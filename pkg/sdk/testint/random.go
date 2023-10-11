package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal"
)

func alphanumericDatabaseObjectIdentifier(t *testing.T) sdk.DatabaseObjectIdentifier {
	t.Helper()
	return sdk.NewDatabaseObjectIdentifier(internal.RandomAlphanumericN(t, 12), internal.RandomAlphanumericN(t, 12))
}

func randomAccountObjectIdentifier(t *testing.T) sdk.AccountObjectIdentifier {
	t.Helper()
	return sdk.RandomAccountObjectIdentifier(t)
}
