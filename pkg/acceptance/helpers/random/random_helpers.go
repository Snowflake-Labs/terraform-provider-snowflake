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

// AdminName returns admin name acceptable by Snowflake:
// 090088 (22000): ADMIN_NAME can only contain letters, numbers and underscores.
// 090089 (22000): ADMIN_NAME must start with a letter.
func AdminName() string {
	return AlphaN(1) + AlphanumericN(11)
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
