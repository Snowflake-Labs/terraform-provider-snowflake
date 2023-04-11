package snowflake

import (
	"database/sql"
	"reflect"
)

type PasswordPolicy struct {
	SchemaObjectIdentifier

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

type PasswordPolicyManager struct {
	BaseManager
}

func NewPasswordPolicyManager() (*PasswordPolicyManager, error) {
	builder, err := newBuilder(
		"PASSWORD POLICY",
		"PASSWORD POLICIES",
		reflect.TypeOf(PasswordPolicyCreateInput{}),
		reflect.TypeOf(PasswordPolicyUpdateInput{}),
		reflect.TypeOf(PasswordPolicyUpdateInput{}),
		reflect.TypeOf(PasswordPolicyDeleteInput{}),
		reflect.TypeOf(PasswordPolicyReadOutput{}),
	)
	if err != nil {
		return nil, err
	}

	return &PasswordPolicyManager{
		BaseManager: BaseManager{
			genericBuilder: *builder,
		},
	}, nil
}

type PasswordPolicyCreateInput struct {
	PasswordPolicy

	OrReplace     bool `pos:"beforeObjectType" value:"OR REPLACE"`
	OrReplaceOk   bool
	IfNotExists   bool `pos:"afterObjectType" value:"IF NOT EXISTS"`
	IfNotExistsOk bool
}

func (m *PasswordPolicyManager) Create(x *PasswordPolicyCreateInput) (string, error) {
	return m.genericBuilder.Create(x)
}

type PasswordPolicyReadInput = SchemaObjectIdentifier
type PasswordPolicyReadOutput = PasswordPolicy

func (m *PasswordPolicyManager) Read(x *PasswordPolicyReadInput) (string, error) {
	return m.genericBuilder.Describe(x)
}

func (m *PasswordPolicyManager) Parse(rows *sql.Rows) (*PasswordPolicyReadOutput, error) {
	output := &PasswordPolicyReadOutput{}
	err := m.genericBuilder.ParseDescribe(rows, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

type PasswordPolicyUpdateInput struct {
	PasswordPolicy

	IfExists   bool `pos:"afterObjectType" value:"IF EXISTS"`
	IfExistsOk bool
}

func (m *PasswordPolicyManager) Update(x *PasswordPolicyUpdateInput) (string, error) {
	return m.genericBuilder.Alter(x)
}
func (m *PasswordPolicyManager) Unset(x *PasswordPolicyUpdateInput) (string, error) {
	return m.genericBuilder.Unset(x)
}

type PasswordPolicyDeleteInput struct {
	SchemaObjectIdentifier

	IfExists   bool `pos:"afterObjectType" value:"IF EXISTS"`
	IfExistsOk bool
}

func (m *PasswordPolicyManager) Delete(x *PasswordPolicyDeleteInput) (string, error) {
	return m.genericBuilder.Drop(x)
}
