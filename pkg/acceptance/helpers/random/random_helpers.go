package random

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/go-uuid"
)

// generatedRandomValue is used to mask random values in GitHub Action logs.
// It always starts with a letter and contains only letters and numbers.
var generatedRandomValue string

func init() {
	generatedRandomValue = os.Getenv(string(testenvs.GeneratedRandomValue))
	requireGeneratedRandomValue := os.Getenv(string(testenvs.RequireGeneratedRandomValue))
	if requireGeneratedRandomValue != "" && generatedRandomValue == "" {
		log.Println("generated random value is required for tests to run")
		os.Exit(1)
	}
}

func UUID() string {
	v, _ := uuid.GenerateUUID()
	return v
}

func Comment() string {
	return gofakeit.Sentence(10)
}

// AdminName returns admin name acceptable by Snowflake:
// 090088 (22000): ADMIN_NAME can only contain letters, numbers and underscores.
// 090089 (22000): ADMIN_NAME must start with a letter.
func AdminName() string {
	return SensitiveAlphanumeric()
}

func Email() string {
	return SensitiveAlphanumeric() + gofakeit.Email()
}

func Password() string {
	return SensitiveString()
}

// SensitiveString returns a random string prefixed with a generated random value that is masked in GitHub Action logs.
// The string returned by SensitiveString always starts with a letter and contains only letters, numbers, and symbols.
func SensitiveString() string {
	return generatedRandomValue + StringN(30)
}

// SensitiveAlphanumeric returns a random string prefixed with a generated random value that is masked in GitHub Action logs.
// The string returned by SensitiveAlphanumeric always starts with a letter and contains only letters and numbers.
func SensitiveAlphanumeric() string {
	return generatedRandomValue + AlphanumericN(30)
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

func AlphaLowerN(num int) string {
	return gofakeit.Password(true, false, false, false, false, num)
}

func Bytes() []byte {
	return []byte(AlphaN(10))
}
