package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func TestEventTables_Create(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateEventTableOptions {
		return &CreateEventTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateEventTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.ClusterBy = []string{"a", "b"}
		opts.DataRetentionTimeInDays = Int(1)
		opts.MaxDataExtensionTimeInDays = Int(2)
		opts.ChangeTracking = Bool(true)
		opts.DefaultDdlCollation = String("en_US")
		opts.CopyGrants = Bool(true)
		opts.Comment = String("test")
		pn := NewSchemaObjectIdentifier(random.StringN(4), random.StringN(4), random.StringN(4))
		opts.RowAccessPolicy = &TableRowAccessPolicy{
			Name: pn,
			On:   []string{"c1", "c2"},
		}
		tn := NewSchemaObjectIdentifier(random.StringN(4), random.StringN(4), random.StringN(4))
		opts.Tag = []TagAssociation{
			{
				Name:  tn,
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EVENT TABLE %s CLUSTER BY (a, b) DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 2 CHANGE_TRACKING = true DEFAULT_DDL_COLLATION = 'en_US' COPY GRANTS COMMENT = 'test' ROW ACCESS POLICY %s ON (c1, c2) TAG (%s = 'v1')`, id.FullyQualifiedName(), pn.FullyQualifiedName(), tn.FullyQualifiedName())
	})
}

func TestEventTables_Show(t *testing.T) {
	id := RandomSchemaObjectIdentifier()
	defaultOpts := func() *ShowEventTableOptions {
		return &ShowEventTableOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowEventTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW EVENT TABLES IN DATABASE "database"`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW EVENT TABLES LIKE '%s'`, id.Name())
	})

	t.Run("show with like and in", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE EVENT TABLES LIKE '%s' IN DATABASE "database"`, id.Name())
	})
}

func TestEventTables_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *DescribeEventTableOptions {
		return &DescribeEventTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeEventTableOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE EVENT TABLE %s`, id.FullyQualifiedName())
	})
}

func TestEventTables_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *AlterEventTableOptions {
		return &AlterEventTableOptions{
			name:        id,
			IfNotExists: Bool(true),
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterEventTableOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterEventTableOptions", "RenameTo", "Set", "Unset", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "ClusteringAction", "SearchOptimizationAction"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAllRowAccessPolicies = Bool(true)
		opts.Set = &EventTableSet{
			DataRetentionTimeInDays: Int(1),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterEventTableOptions", "RenameTo", "Set", "Unset", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "ClusteringAction", "SearchOptimizationAction"))
	})

	t.Run("alter: rename to", func(t *testing.T) {
		opts := defaultOpts()
		target := NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), random.StringN(12))
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), target.FullyQualifiedName())
	})

	t.Run("alter: clustering action", func(t *testing.T) {
		opts := defaultOpts()
		cluster := []string{"a", "b"}
		opts.ClusteringAction = &EventTableClusteringAction{
			ClusterBy: &cluster,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s CLUSTER BY (a, b)`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.ClusteringAction = &EventTableClusteringAction{
			SuspendRecluster: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s SUSPEND RECLUSTER`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.ClusteringAction = &EventTableClusteringAction{
			ResumeRecluster: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s RESUME RECLUSTER`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.ClusteringAction = &EventTableClusteringAction{
			DropClusteringKey: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s DROP CLUSTERING KEY`, id.FullyQualifiedName())
	})

	t.Run("alter: search optimization action", func(t *testing.T) {
		opts := defaultOpts()
		opts.SearchOptimizationAction = &EventTableSearchOptimizationAction{
			Add: &SearchOptimization{
				On: []string{"EQUALITY(*)", "SUBSTRING(*)"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s ADD SEARCH OPTIMIZATION ON EQUALITY(*), SUBSTRING(*)`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.SearchOptimizationAction = &EventTableSearchOptimizationAction{
			Drop: &SearchOptimization{
				On: []string{"EQUALITY(*)", "SUBSTRING(*)"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s DROP SEARCH OPTIMIZATION ON EQUALITY(*), SUBSTRING(*)`, id.FullyQualifiedName())
	})

	t.Run("alter: set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &EventTableSet{
			DataRetentionTimeInDays: Int(1),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s SET DATA_RETENTION_TIME_IN_DAYS = 1`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.Set = &EventTableSet{
			MaxDataExtensionTimeInDays: Int(1),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s SET MAX_DATA_EXTENSION_TIME_IN_DAYS = 1`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.Set = &EventTableSet{
			ChangeTracking: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s SET CHANGE_TRACKING = true`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.Set = &EventTableSet{
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s SET COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &EventTableUnset{
			DataRetentionTimeInDays: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s UNSET DATA_RETENTION_TIME_IN_DAYS`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.Unset = &EventTableUnset{
			MaxDataExtensionTimeInDays: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s UNSET MAX_DATA_EXTENSION_TIME_IN_DAYS`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.Unset = &EventTableUnset{
			ChangeTracking: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s UNSET CHANGE_TRACKING`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.Unset = &EventTableUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("alter: set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s SET TAG "tag1" = 'value1'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})

	t.Run("alter: add row access policy", func(t *testing.T) {
		rowAccessPolicyId := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.AddRowAccessPolicy = &EventTableAddRowAccessPolicy{
			RowAccessPolicy: rowAccessPolicyId,
			On:              []string{"a", "b"},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE IF NOT EXISTS %s ADD ROW ACCESS POLICY %s ON (a, b)", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("alter: drop row access policy", func(t *testing.T) {
		rowAccessPolicyId := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.DropRowAccessPolicy = &EventTableDropRowAccessPolicy{
			RowAccessPolicy: rowAccessPolicyId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE IF NOT EXISTS %s DROP ROW ACCESS POLICY %s", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("alter: drop and add row access policy", func(t *testing.T) {
		rowAccessPolicy1Id := RandomSchemaObjectIdentifier()
		rowAccessPolicy2Id := RandomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.DropAndAddRowAccessPolicy = &EventTableDropAndAddRowAccessPolicy{
			Drop: EventTableDropRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicy1Id,
			},
			Add: EventTableAddRowAccessPolicy{
				RowAccessPolicy: rowAccessPolicy2Id,
				On:              []string{"a", "b"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TABLE IF NOT EXISTS %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON (a, b)", id.FullyQualifiedName(), rowAccessPolicy1Id.FullyQualifiedName(), rowAccessPolicy2Id.FullyQualifiedName())
	})

	t.Run("alter: drop all row access policies", func(t *testing.T) {
		opts := defaultOpts()
		opts.DropAllRowAccessPolicies = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE IF NOT EXISTS %s DROP ALL ROW ACCESS POLICIES`, id.FullyQualifiedName())
	})
}
