package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRowAccessPolicies_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid CreateRowAccessPolicyOptions
	defaultOpts := func() *CreateRowAccessPolicyOptions {
		return &CreateRowAccessPolicyOptions{
			name: id,
			args: []CreateRowAccessPolicyArgs{{
				Name: "n",
				Type: DataTypeVARCHAR,
			}},
			body: "true",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateRowAccessPolicyOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: [opts.args] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.args = []CreateRowAccessPolicyArgs{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateRowAccessPolicyOptions", "args"))
	})

	t.Run("validation: [opts.body] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.body = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateRowAccessPolicyOptions", "body"))
	})

	t.Run("one parameter", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE ROW ACCESS POLICY %s AS ("n" VARCHAR) RETURNS BOOLEAN -> true`, id.FullyQualifiedName())
	})

	t.Run("two parameters", func(t *testing.T) {
		opts := defaultOpts()
		opts.args = []CreateRowAccessPolicyArgs{{
			Name: "n",
			Type: DataTypeVARCHAR,
		}, {
			Name: "h",
			Type: DataTypeVARCHAR,
		}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE ROW ACCESS POLICY %s AS ("n" VARCHAR, "h" VARCHAR) RETURNS BOOLEAN -> true`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE ROW ACCESS POLICY %s AS ("n" VARCHAR) RETURNS BOOLEAN -> true COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestRowAccessPolicies_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid AlterRowAccessPolicyOptions
	defaultOpts := func() *AlterRowAccessPolicyOptions {
		return &AlterRowAccessPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetBody opts.SetTags opts.UnsetTags opts.SetComment opts.UnsetComment] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterRowAccessPolicyOptions", "RenameTo", "SetBody", "SetTags", "UnsetTags", "SetComment", "UnsetComment"))
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetBody opts.SetTags opts.UnsetTags opts.SetComment opts.UnsetComment] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("comment")
		opts.UnsetComment = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterRowAccessPolicyOptions", "RenameTo", "SetBody", "SetTags", "UnsetTags", "SetComment", "UnsetComment"))
	})

	t.Run("rename", func(t *testing.T) {
		newId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, "ALTER ROW ACCESS POLICY %s RENAME TO %s", id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set body", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetBody = String("true")
		assertOptsValidAndSQLEquals(t, opts, "ALTER ROW ACCESS POLICY %s SET BODY -> true", id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "ALTER ROW ACCESS POLICY %s SET COMMENT = 'comment'", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER ROW ACCESS POLICY %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
			{
				Name:  NewAccountObjectIdentifier("tag2"),
				Value: "value2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROW ACCESS POLICY %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROW ACCESS POLICY %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestRowAccessPolicies_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DropRowAccessPolicyOptions
	defaultOpts := func() *DropRowAccessPolicyOptions {
		return &DropRowAccessPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP ROW ACCESS POLICY %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP ROW ACCESS POLICY IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestRowAccessPolicies_Show(t *testing.T) {
	// Minimal valid ShowRowAccessPolicyOptions
	defaultOpts := func() *ShowRowAccessPolicyOptions {
		return &ShowRowAccessPolicyOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW ROW ACCESS POLICIES")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("myaccount"),
		}
		opts.In = &ExtendedIn{
			In: In{
				Account: Bool(true),
			},
		}
		opts.Limit = &LimitFrom{
			Rows: Pointer(10),
			From: Pointer("foo"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ROW ACCESS POLICIES LIKE 'myaccount' IN ACCOUNT LIMIT 10 FROM 'foo'")
	})
}

func TestRowAccessPolicies_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DescribeRowAccessPolicyOptions
	defaultOpts := func() *DescribeRowAccessPolicyOptions {
		return &DescribeRowAccessPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeRowAccessPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE ROW ACCESS POLICY %s", id.FullyQualifiedName())
	})
}

func TestRowAccessPolicyDescription_Arguments(t *testing.T) {
	tests := []struct {
		name      string
		signature string
		want      []RowAccessPolicyArgument
	}{
		{
			name:      "signature with 1 arg",
			signature: "(A VARCHAR)",
			want: []RowAccessPolicyArgument{
				{
					Name: "A",
					Type: "VARCHAR",
				},
			},
		},
		{
			name:      "signature with multiple args",
			signature: "(A VARCHAR, B BOOLEAN)",
			want: []RowAccessPolicyArgument{
				{
					Name: "A",
					Type: "VARCHAR",
				},
				{
					Name: "B",
					Type: "BOOLEAN",
				},
			},
		},
		{
			name:      "signature with complex name",
			signature: "(a B VARCHAR)",
			want: []RowAccessPolicyArgument{
				{
					Name: "a B",
					Type: "VARCHAR",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &describeRowAccessPolicyDBRow{
				Signature: tt.signature,
			}
			got := d.convert()
			require.Equal(t, tt.want, got.Signature)
		})
	}
}
