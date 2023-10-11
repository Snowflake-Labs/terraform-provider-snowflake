package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// All functions in this file exist also in SDK (they are used also there).
// They were copied to allow easier extraction of all integration tests to separate package.
// This will be dealt with in subsequent PRs.

func randomSchemaObjectIdentifier(t *testing.T) sdk.SchemaObjectIdentifier {
	t.Helper()
	return sdk.RandomSchemaObjectIdentifier(t)
}

func randomDatabaseObjectIdentifier(t *testing.T) sdk.DatabaseObjectIdentifier {
	t.Helper()
	return sdk.RandomDatabaseObjectIdentifier(t)
}

func alphanumericDatabaseObjectIdentifier(t *testing.T) sdk.DatabaseObjectIdentifier {
	t.Helper()
	return sdk.AlphanumericDatabaseObjectIdentifier(t)
}

func randomAccountObjectIdentifier(t *testing.T) sdk.AccountObjectIdentifier {
	t.Helper()
	return sdk.RandomAccountObjectIdentifier(t)
}
