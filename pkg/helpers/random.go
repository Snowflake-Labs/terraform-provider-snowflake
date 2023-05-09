package helpers

import (
	"github.com/brianvoe/gofakeit/v6"
)

func RandomBool() bool {
	return gofakeit.Bool()
}

func RandomString() string {
	return gofakeit.Password(true, true, true, true, false, 28)
}

func RandomStringRange(min, max int) string {
	if min > max {
		return ""
	}
	return gofakeit.Password(true, true, true, true, false, RandomIntRange(min, max))
}

func RandomIntRange(min, max int) int {
	if min > max {
		return 0
	}
	return gofakeit.IntRange(min, max)
}
