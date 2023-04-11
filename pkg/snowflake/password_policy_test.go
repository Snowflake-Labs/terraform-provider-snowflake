package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestCreatePasswordPolicy(t *testing.T) {
	r := require.New(t)

	input := &snowflake.PasswordPolicyCreateInput{
		PasswordPolicy: snowflake.PasswordPolicy{
			SchemaObjectIdentifier: snowflake.SchemaObjectIdentifier{
				Database:   "testdb",
				Schema:     "testschema",
				ObjectName: "testres",
			},
			MinLength:   10,
			MinLengthOk: true,
		},
		OrReplace:     true,
		OrReplaceOk:   true,
		IfNotExists:   false,
		IfNotExistsOk: true,
	}

	mb, err := snowflake.NewPasswordPolicyManager()
	r.Nil(err)
	createStmt, err := mb.Create(input)
	r.Nil(err)
	r.Equal(`CREATE OR REPLACE PASSWORD POLICY testdb.testschema.testres PASSWORD_MIN_LENGTH = 10;`, createStmt)
}

func TestAlterPasswordPolicy(t *testing.T) {
	r := require.New(t)

	input := &snowflake.PasswordPolicyUpdateInput{
		PasswordPolicy: snowflake.PasswordPolicy{
			SchemaObjectIdentifier: snowflake.SchemaObjectIdentifier{
				Database:   "testdb",
				Schema:     "testschema",
				ObjectName: "passpol",
			},
			MinNumericChars:   16,
			MinNumericCharsOk: true,
			LockoutTimeMins:   50,
			LockoutTimeMinsOk: true,
		},
	}

	mb, err := snowflake.NewPasswordPolicyManager()
	r.Nil(err)
	alterStmt, err := mb.Update(input)
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

	input := &snowflake.PasswordPolicyUpdateInput{
		PasswordPolicy: snowflake.PasswordPolicy{
			SchemaObjectIdentifier: snowflake.SchemaObjectIdentifier{
				Database:   "testdb",
				Schema:     "testschema",
				ObjectName: "passpol",
			},
			MinNumericCharsOk: true,
			CommentOk:         true,
		},
	}

	mb, err := snowflake.NewPasswordPolicyManager()
	r.Nil(err)
	unsetStmt, err := mb.Unset(input)
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

func TestDeletePasswordPolicy(t *testing.T) {
	r := require.New(t)

	input := &snowflake.PasswordPolicyDeleteInput{
		SchemaObjectIdentifier: snowflake.SchemaObjectIdentifier{
			Database:   "testdb",
			Schema:     "testschema",
			ObjectName: "passpol",
		},
	}

	mb, err := snowflake.NewPasswordPolicyManager()
	r.Nil(err)
	dropStmt, err := mb.Delete(input)
	r.Nil(err)
	r.Equal(`DROP PASSWORD POLICY testdb.testschema.passpol;`, dropStmt)
}

func TestReadPasswordPolicy(t *testing.T) {
	r := require.New(t)

	input := &snowflake.PasswordPolicyReadInput{
		Database:   "testdb",
		Schema:     "testschema",
		ObjectName: "passpol",
	}

	mb, err := snowflake.NewPasswordPolicyManager()
	r.Nil(err)
	describeStmt, err := mb.Read(input)
	r.Nil(err)
	r.Equal(`DESCRIBE PASSWORD POLICY testdb.testschema.passpol;`, describeStmt)
}
