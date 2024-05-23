package random

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/go-uuid"
)

func UUID() string {
	v, _ := uuid.GenerateUUID()
	return v
}

func Comment() string {
	return gofakeit.Sentence(10)
}

func Password() string {
	return StringN(12)
}

func Bool() bool {
	return gofakeit.Bool()
}

func String() string {
	return gofakeit.Password(true, true, true, true, false, 28)
}

func StringN(num int) string {
	return gofakeit.Password(true, true, true, true, false, num)
}

func AlphanumericN(num int) string {
	return gofakeit.Password(true, true, true, false, false, num)
}

func AlphaN(num int) string {
	return gofakeit.Password(true, true, false, false, false, num)
}

func StringRange(min, max int) string {
	return gofakeit.Password(true, true, true, true, false, IntRange(min, max))
}

func IntRange(min, max int) int {
	return gofakeit.IntRange(min, max)
}
