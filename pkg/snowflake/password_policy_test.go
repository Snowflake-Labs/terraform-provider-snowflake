package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestCreatePasswordPolicy(t *testing.T) {
	r := require.New(t)

	props := &snowflake.PasswordPolicyProps{
		Database:      "testdb",
		Schema:        "testschema",
		Name:          "testres",
		OrReplace:     true,
		OrReplaceOk:   true,
		IfNotExists:   false,
		IfNotExistsOk: true,
		MinLength:     10,
		MinLengthOk:   true,
	}

	mb, err := snowflake.NewPasswordPolicyBuilder()
	r.Nil(err)
	createStmt, err := mb.Create(props)
	r.Nil(err)
	r.Equal(`CREATE OR REPLACE PASSWORD POLICY testdb.testschema.testres PASSWORD_MIN_LENGTH = 10;`, createStmt)
}

func TestAlterPasswordPolicy(t *testing.T) {
	r := require.New(t)

	props := &snowflake.PasswordPolicyProps{
		Database:          "testdb",
		Schema:            "testschema",
		Name:              "passpol",
		MinNumericChars:   16,
		MinNumericCharsOk: true,
		LockoutTimeMins:   50,
		LockoutTimeMinsOk: true,
	}

	mb, err := snowflake.NewPasswordPolicyBuilder()
	r.Nil(err)
	alterStmt, err := mb.Alter(props)
	r.Nil(err)
	// Order of parameters is not guaranteed
	r.Contains(
		[2]string{
			`ALTER PASSWORD POLICY testdb.testschema.passpol SET PASSWORD_MIN_NUMERIC_CHARS = 16 PASSWORD_LOCKOUT_TIME_MINS = 50;`,
			`ALTER PASSWORD POLICY testdb.testschema.passpol SET PASSWORD_LOCKOUT_TIME_MINS = 50 PASSWORD_MIN_NUMERIC_CHARS = 16;`,
		},
		alterStmt,
	)
}

func TestUnsetPasswordPolicy(t *testing.T) {
	r := require.New(t)

	props := &snowflake.PasswordPolicyProps{
		Database:          "testdb",
		Schema:            "testschema",
		Name:              "passpol",
		MinNumericCharsOk: true,
		CommentOk:         true,
	}

	mb, err := snowflake.NewPasswordPolicyBuilder()
	r.Nil(err)
	unsetStmt, err := mb.Unset(props)
	r.Nil(err)
	// Order of parameters is not guaranteed
	r.Contains(
		[2]string{
			`ALTER PASSWORD POLICY testdb.testschema.passpol UNSET PASSWORD_MIN_NUMERIC_CHARS COMMENT;`,
			`ALTER PASSWORD POLICY testdb.testschema.passpol UNSET COMMENT PASSWORD_MIN_NUMERIC_CHARS;`,
		},
		unsetStmt,
	)
}

func TestDropPasswordPolicy(t *testing.T) {
	r := require.New(t)

	props := &snowflake.PasswordPolicyProps{
		Database: "testdb",
		Schema:   "testschema",
		Name:     "passpol",
	}

	mb, err := snowflake.NewPasswordPolicyBuilder()
	r.Nil(err)
	dropStmt, err := mb.Drop(props)
	r.Nil(err)
	r.Equal(`DROP PASSWORD POLICY testdb.testschema.passpol;`, dropStmt)
}

func TestDescribePasswordPolicy(t *testing.T) {
	r := require.New(t)

	props := &snowflake.PasswordPolicyProps{
		Database: "testdb",
		Schema:   "testschema",
		Name:     "passpol",
	}

	mb, err := snowflake.NewPasswordPolicyBuilder()
	r.Nil(err)
	describeStmt, err := mb.Describe(props)
	r.Nil(err)
	r.Equal(`DESCRIBE PASSWORD POLICY testdb.testschema.passpol;`, describeStmt)
}
