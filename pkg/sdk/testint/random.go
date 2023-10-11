package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// All functions in this file exist also in SDK (they are used also there).
// They were copied to allow easier extraction of all integration tests to separate package.
// This will be dealt with in subsequent PRs.

func randomUUID(t *testing.T) string {
	return sdk.RandomUUID(t)
}

func randomComment(t *testing.T) string {
	return sdk.RandomComment(t)
}

func randomBool(t *testing.T) bool {
	return sdk.RandomBool(t)
}

func randomString(t *testing.T) string {
	return sdk.RandomString(t)
}

func randomStringN(t *testing.T, num int) string {
	return sdk.RandomStringN(t, num)
}

func randomAlphanumericN(t *testing.T, num int) string {
	return sdk.RandomAlphanumericN(t, num)
}

func randomStringRange(t *testing.T, min, max int) string {
	return sdk.RandomStringRange(t, min, max)
}

func randomIntRange(t *testing.T, min, max int) int {
	return sdk.RandomIntRange(t, min, max)
}

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
