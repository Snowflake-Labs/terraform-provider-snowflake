package random

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

func Uuid(t *testing.T) string {
	t.Helper()
	v, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return v
}

func Comment(t *testing.T) string {
	t.Helper()
	return gofakeit.Sentence(10)
}

func Bool(t *testing.T) bool {
	t.Helper()
	return gofakeit.Bool()
}

func String(t *testing.T) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, 28)
}

func StringN(t *testing.T, num int) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, num)
}

func AlphanumericN(t *testing.T, num int) string {
	t.Helper()
	return gofakeit.Password(true, true, true, false, false, num)
}

func StringRange(t *testing.T, min, max int) string {
	t.Helper()
	if min > max {
		t.Errorf("min %d is greater than max %d", min, max)
	}
	return gofakeit.Password(true, true, true, true, false, IntRange(t, min, max))
}

func IntRange(t *testing.T, min, max int) int {
	t.Helper()
	if min > max {
		t.Errorf("min %d is greater than max %d", min, max)
	}
	return gofakeit.IntRange(min, max)
}
