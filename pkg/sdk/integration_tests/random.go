package sdk_integration_tests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/brianvoe/gofakeit/v6"
)

func randomString(t *testing.T) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, 28)
}

func randomStringRange(t *testing.T, min, max int) string {
	t.Helper()
	if min > max {
		t.Errorf("min %d is greater than max %d", min, max)
	}
	return gofakeit.Password(true, true, true, true, false, randomIntRange(t, min, max))
}

func randomIntRange(t *testing.T, min, max int) int {
	t.Helper()
	if min > max {
		t.Errorf("min %d is greater than max %d", min, max)
	}
	return gofakeit.IntRange(min, max)
}

func randomAlphanumericN(t *testing.T, num int) string {
	t.Helper()
	return gofakeit.Password(true, true, true, false, false, num)
}

func alphanumericSchemaIdentifier(t *testing.T) sdk.SchemaIdentifier {
	t.Helper()
	return sdk.NewSchemaIdentifier(randomAlphanumericN(t, 12), randomAlphanumericN(t, 12))
}

func randomComment(t *testing.T) string {
	t.Helper()
	return gofakeit.Sentence(10)
}
