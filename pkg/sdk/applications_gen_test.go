package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
)

func TestApplications_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()
	pid := RandomAccountObjectIdentifier()

	defaultOpts := func() *CreateApplicationOptions {
		return &CreateApplicationOptions{
			name:        id,
			PackageName: pid,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateApplicationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Version = &ApplicationVersion{
			VersionAndPatch: &VersionAndPatch{
				Version: "1.0",
				Patch:   Int(1),
			},
			VersionDirectory: String("@test"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateApplicationOptions.Version", "VersionDirectory", "VersionAndPatch"))
	})

	t.Run("validation: version must be set when debug mode is set", func(t *testing.T) {
		opts := defaultOpts()
		opts.DebugMode = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, NewError("CreateApplicationOptions.DebugMode can be set only when CreateApplicationOptions.Version is set"))
	})

	t.Run("all options", func(t *testing.T) {
		tid := NewSchemaObjectIdentifier(random.StringN(4), random.StringN(4), random.StringN(4))

		opts := defaultOpts()
		opts.Comment = String("test")
		opts.Tag = []TagAssociation{
			{
				Name:  tid,
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE APPLICATION %s FROM APPLICATION PACKAGE %s COMMENT = 'test' TAG (%s = 'v1')`, id.FullyQualifiedName(), pid.FullyQualifiedName(), tid.FullyQualifiedName())

		opts = defaultOpts()
		opts.Comment = String("test")
		opts.Version = &ApplicationVersion{
			VersionDirectory: String("@test"),
		}
		opts.DebugMode = Bool(true)
		opts.Tag = []TagAssociation{
			{
				Name:  tid,
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE APPLICATION %s FROM APPLICATION PACKAGE %s USING '@test' DEBUG_MODE = true COMMENT = 'test' TAG (%s = 'v1')`, id.FullyQualifiedName(), pid.FullyQualifiedName(), tid.FullyQualifiedName())

		opts = defaultOpts()
		opts.Comment = String("test")
		opts.Version = &ApplicationVersion{
			VersionAndPatch: &VersionAndPatch{
				Version: "V001",
				Patch:   Int(1),
			},
		}
		opts.DebugMode = Bool(true)
		opts.Tag = []TagAssociation{
			{
				Name:  tid,
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE APPLICATION %s FROM APPLICATION PACKAGE %s USING VERSION V001 PATCH 1 DEBUG_MODE = true COMMENT = 'test' TAG (%s = 'v1')`, id.FullyQualifiedName(), pid.FullyQualifiedName(), tid.FullyQualifiedName())
	})
}

func TestApplications_Alter(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	defaultOpts := func() *AlterApplicationOptions {
		return &AlterApplicationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterApplicationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApplicationOptions", "Set", "Unset", "Upgrade", "UpgradeVersion", "UnsetReferences", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Upgrade = Bool(true)
		opts.Unset = &ApplicationUnset{
			Comment: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApplicationOptions", "Set", "Unset", "Upgrade", "UpgradeVersion", "UnsetReferences", "SetTags", "UnsetTags"))
	})

	t.Run("validation: if exits can be set only when set or unset is set", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, NewError("AlterApplicationOptions.IfExists can be set only when AlterApplicationOptions.Set or AlterApplicationOptions.Unset is set"))
	})

	t.Run("alter: set options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &ApplicationSet{
			ShareEventsWithProvider: Bool(true),
			DebugMode:               Bool(true),
			Comment:                 String("test"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION IF EXISTS %s SET COMMENT = 'test' SHARE_EVENTS_WITH_PROVIDER = true DEBUG_MODE = true`, id.FullyQualifiedName())
	})

	t.Run("alter: unset options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &ApplicationUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &ApplicationUnset{
			ShareEventsWithProvider: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION IF EXISTS %s UNSET SHARE_EVENTS_WITH_PROVIDER`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &ApplicationUnset{
			DebugMode: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION IF EXISTS %s UNSET DEBUG_MODE`, id.FullyQualifiedName())
	})

	t.Run("alter: set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION %s SET TAG "tag1" = 'value1'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})

	t.Run("alter: upgrade", func(t *testing.T) {
		opts := defaultOpts()
		opts.Upgrade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION %s UPGRADE`, id.FullyQualifiedName())
	})

	t.Run("alter: upgrade version", func(t *testing.T) {
		opts := defaultOpts()
		opts.UpgradeVersion = &ApplicationVersion{
			VersionDirectory: String("@test"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION %s UPGRADE USING '@test'`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.UpgradeVersion = &ApplicationVersion{
			VersionAndPatch: &VersionAndPatch{
				Version: "V001",
				Patch:   Int(1),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION %s UPGRADE USING VERSION V001 PATCH 1`, id.FullyQualifiedName())
	})

	t.Run("alter: unset references", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetReferences = &ApplicationReferences{}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION %s UNSET REFERENCES`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.UnsetReferences = &ApplicationReferences{
			References: []ApplicationReference{
				{
					Reference: "ref1",
				},
				{
					Reference: "ref2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION %s UNSET REFERENCES ('ref1', 'ref2')`, id.FullyQualifiedName())
	})
}

func TestApplications_Drop(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	defaultOpts := func() *DropApplicationOptions {
		return &DropApplicationOptions{
			name: id,
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropApplicationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Cascade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP APPLICATION IF EXISTS %s CASCADE`, id.FullyQualifiedName())
	})
}

func TestApplications_Describe(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	defaultOpts := func() *DescribeApplicationOptions {
		return &DescribeApplicationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeApplicationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE APPLICATION %s`, id.FullyQualifiedName())
	})
}

func TestApplications_Show(t *testing.T) {
	defaultOpts := func() *ShowApplicationOptions {
		return &ShowApplicationOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowApplicationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW APPLICATIONS`)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		opts.StartsWith = String("A")
		opts.Limit = &LimitFrom{
			Rows: Int(1),
			From: String("B"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW APPLICATIONS LIKE 'pattern' STARTS WITH 'A' LIMIT 1 FROM 'B'`)
	})
}
