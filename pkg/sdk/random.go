package sdk

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

func randomUUID(t *testing.T) string {
	t.Helper()
	v, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return v
}

func randomComment(t *testing.T) string {
	t.Helper()
	return gofakeit.Sentence(10)
}

func randomBool(t *testing.T) bool {
	t.Helper()
	return gofakeit.Bool()
}

func randomString(t *testing.T) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, 28)
}

func randomStringN(t *testing.T, num int) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, num)
}

func randomAlphanumericN(t *testing.T, num int) string {
	t.Helper()
	return gofakeit.Password(true, true, true, false, false, num)
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

func randomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
	t.Helper()
	return NewSchemaObjectIdentifier(randomStringN(t, 12), randomStringN(t, 12), randomStringN(t, 12))
}

func randomDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
	t.Helper()
	return NewDatabaseObjectIdentifier(randomStringN(t, 12), randomStringN(t, 12))
}

func alphanumericDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
	t.Helper()
	return NewDatabaseObjectIdentifier(randomAlphanumericN(t, 12), randomAlphanumericN(t, 12))
}

func randomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
	t.Helper()
	return NewAccountObjectIdentifier(randomStringN(t, 12))
}
