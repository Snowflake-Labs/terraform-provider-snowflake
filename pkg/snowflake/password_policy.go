package snowflake

import (
	"fmt"
	"reflect"
)

type PasswordPolicyProps struct {
	Database string
	Schema   string
	Name     string `pos:"name" db:"NAME"`

	// Keywords
	OrReplace     bool `pos:"beforeType" value:"OR REPLACE"`
	OrReplaceOk   bool
	IfNotExists   bool `pos:"afterType" value:"IF NOT EXISTS"`
	IfNotExistsOk bool

	MinLength           int `pos:"parameter" db:"PASSWORD_MIN_LENGTH"`
	MinLengthOk         bool
	MaxLength           int `pos:"parameter" db:"PASSWORD_MAX_LENGTH"`
	MaxLengthOk         bool
	MinUpperCaseChars   int `pos:"parameter" db:"PASSWORD_MIN_UPPER_CASE_CHARS"`
	MinUpperCaseCharsOk bool
	MinLowerCaseChars   int `pos:"parameter" db:"PASSWORD_MIN_LOWER_CASE_CHARS"`
	MinLowerCaseCharsOk bool
	MinNumericChars     int `pos:"parameter" db:"PASSWORD_MIN_NUMERIC_CHARS"`
	MinNumericCharsOk   bool
	MinSpecialChars     int `pos:"parameter" db:"PASSWORD_MIN_SPECIAL_CHARS"`
	MinSpecialCharsOk   bool
	MaxAgeDays          int `pos:"parameter" db:"PASSWORD_MAX_AGE_DAYS"`
	MaxAgeDaysOk        bool
	MaxRetries          int `pos:"parameter" db:"PASSWORD_MAX_RETRIES"`
	MaxRetriesOk        bool
	LockoutTimeMins     int `pos:"parameter" db:"PASSWORD_LOCKOUT_TIME_MINS"`
	LockoutTimeMinsOk   bool

	Comment   string `pos:"parameter" db:"COMMENT"`
	CommentOk bool
}

func (pp *PasswordPolicyProps) QualifiedName() string {
	return fmt.Sprintf("%v.%v.%v", pp.Database, pp.Schema, pp.Name)
}

func (pp *PasswordPolicyProps) ID() string {
	return pp.QualifiedName()
}

func NewPasswordPolicyBuilder() (*NewBuilder, error) {
	return newBuilder(
		"PASSWORD POLICY",
		"PASSWORD POLICIES",
		reflect.TypeOf(PasswordPolicyProps{}),
	)
}
