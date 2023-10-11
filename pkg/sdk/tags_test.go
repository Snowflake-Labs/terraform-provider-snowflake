package sdk

import (
	"testing"
)

func TestTagCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)
	defaultOpts := func() *createTagOptions {
		return &createTagOptions{
			name: id,
		}
	}

	t.Run("create with allowed values", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.AllowedValues = &AllowedValues{
			Values: []AllowedValue{
				{
					Value: "value1",
				},
				{
					Value: "value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TAG %s ALLOWED_VALUES 'value1', 'value2'`, id.FullyQualifiedName())
	})

	t.Run("create with comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TAG %s COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("create with all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(false)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE TAG IF NOT EXISTS %s COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*createTagOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: both AllowedValues and Comment present", func(t *testing.T) {
		opts := defaultOpts()
		opts.AllowedValues = &AllowedValues{
			Values: []AllowedValue{
				{
					Value: "value1",
				},
			},
		}
		opts.Comment = String("comment")
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("Comment", "AllowedValues"))
	})

	t.Run("validation: both ifNotExists and orReplace present", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("OrReplace", "IfNotExists"))
	})

	t.Run("validation: multiple errors", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier, errOneOf("OrReplace", "IfNotExists"))
	})
}

func TestTagDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)
	defaultOpts := func() *dropTagOptions {
		return &dropTagOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropTagOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("drop with name", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP TAG %s`, id.FullyQualifiedName())
	})

	t.Run("drop with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP TAG IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestTagUndrop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)
	defaultOpts := func() *undropTagOptions {
		return &undropTagOptions{
			name: id,
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropTagOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("undrop with name", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `UNDROP TAG %s`, id.FullyQualifiedName())
	})
}

func TestTagShow(t *testing.T) {
	defaultOpts := func() *showTagOptions {
		return &showTagOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *showTagOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, errPatternRequiredForLikeKeyword)
	})

	t.Run("validation: empty in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{}
		assertOptsInvalidJoinedErrors(t, opts, errScopeRequiredForInKeyword)
	})

	t.Run("show with empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW TAGS`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{Pattern: String("test")}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TAGS LIKE 'test'`)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TAGS IN ACCOUNT`)
	})
}

func TestTagAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)
	defaultOpts := func() *alterTagOptions {
		return &alterTagOptions{
			name: id,
		}
	}
	defaultAllowedValues := func() *AllowedValues {
		return &AllowedValues{
			Values: []AllowedValue{
				{
					Value: "value1",
				},
				{
					Value: "value2",
				},
			},
		}
	}
	mp1ID := NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), "policy1")
	mp2ID := NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), "policy2")
	defaultMaskingPolicies := func() []TagMaskingPolicy {
		return []TagMaskingPolicy{
			{
				Name: mp1ID,
			},
			{
				Name: mp2ID,
			},
		}
	}

	t.Run("alter with rename to", func(t *testing.T) {
		opts := defaultOpts()
		opts.Rename = &TagRename{Name: NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), randomStringN(t, 12))}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s RENAME TO %s`, id.FullyQualifiedName(), opts.Rename.Name.FullyQualifiedName())
	})

	t.Run("alter with add", func(t *testing.T) {
		opts := defaultOpts()
		opts.Add = &TagAdd{AllowedValues: defaultAllowedValues()}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s ADD ALLOWED_VALUES 'value1', 'value2'`, id.FullyQualifiedName())
	})

	t.Run("alter with drop", func(t *testing.T) {
		opts := defaultOpts()
		opts.Drop = &TagDrop{AllowedValues: defaultAllowedValues()}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s DROP ALLOWED_VALUES 'value1', 'value2'`, id.FullyQualifiedName())
	})

	t.Run("alter with unset allowed values", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{AllowedValues: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s UNSET ALLOWED_VALUES`, id.FullyQualifiedName())
	})

	t.Run("alter with set masking policies", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{
			MaskingPolicies: &TagSetMaskingPolicies{
				MaskingPolicies: defaultMaskingPolicies(),
				Force:           Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s SET MASKING POLICY %s, MASKING POLICY %s FORCE`, id.FullyQualifiedName(), mp1ID.FullyQualifiedName(), mp2ID.FullyQualifiedName())
	})

	t.Run("alter with unset masking policies", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{
			MaskingPolicies: &TagUnsetMaskingPolicies{
				MaskingPolicies: defaultMaskingPolicies(),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s UNSET MASKING POLICY %s, MASKING POLICY %s`, id.FullyQualifiedName(), mp1ID.FullyQualifiedName(), mp2ID.FullyQualifiedName())
	})

	t.Run("alter with set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{Comment: String("comment")}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s SET COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("alter with unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{Comment: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*createTagOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{
			Comment: String("comment"),
		}
		opts.Unset = &TagUnset{
			AllowedValues: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: invalid new name", func(t *testing.T) {
		opts := defaultOpts()
		opts.Rename = &TagRename{
			Name: NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: new name from different db", func(t *testing.T) {
		newId := NewSchemaObjectIdentifier(id.DatabaseName()+randomStringN(t, 1), randomStringN(t, 12), randomStringN(t, 12))

		opts := defaultOpts()
		opts.Rename = &TagRename{
			Name: newId,
		}
		assertOptsValid(t, opts)
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsAtLeastOneProperty)
	})
}
