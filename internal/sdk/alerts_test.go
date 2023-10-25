// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk/internal/random"
)

func TestAlertCreate(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("with complete options", func(t *testing.T) {
		newComment := random.String()
		warehouse := AccountObjectIdentifier{"warehouse"}
		existsCondition := "SELECT 1"
		condition := AlertCondition{[]string{existsCondition}}
		schedule := "1 minute"
		action := "INSERT INTO FOO VALUES (1)"

		opts := &CreateAlertOptions{
			name:      id,
			warehouse: warehouse,
			schedule:  schedule,
			condition: []AlertCondition{condition},
			action:    action,
			Comment:   String(newComment),
		}

		assertOptsValidAndSQLEquals(t, opts, `CREATE ALERT %s WAREHOUSE = "%s" SCHEDULE = '%s' COMMENT = '%s' IF (EXISTS (%s)) THEN %s`, id.FullyQualifiedName(), warehouse.name, schedule, newComment, existsCondition, action)
	})
}

func TestAlertAlter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("fail without alter action specified", func(t *testing.T) {
		opts := &AlterAlertOptions{
			name: id,
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("Action", "Set", "Unset", "ModifyCondition", "ModifyAction"))
	})

	t.Run("fail when 2 alter actions specified", func(t *testing.T) {
		newComment := random.String()
		opts := &AlterAlertOptions{
			name:   id,
			Action: &AlertActionResume,
			Set: &AlertSet{
				Comment: String(newComment),
			},
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("Action", "Set", "Unset", "ModifyCondition", "ModifyAction"))
	})

	t.Run("with resume", func(t *testing.T) {
		opts := &AlterAlertOptions{
			name:   id,
			Action: &AlertActionResume,
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER ALERT %s RESUME", id.FullyQualifiedName())
	})

	t.Run("with suspend", func(t *testing.T) {
		opts := &AlterAlertOptions{
			name:   id,
			Action: &AlertActionSuspend,
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER ALERT %s SUSPEND", id.FullyQualifiedName())
	})

	t.Run("with set", func(t *testing.T) {
		newComment := random.String()
		opts := &AlterAlertOptions{
			name: id,
			Set: &AlertSet{
				Comment: String(newComment),
			},
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER ALERT %s SET COMMENT = '%s'", id.FullyQualifiedName(), newComment)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &AlterAlertOptions{
			name: id,
			Unset: &AlertUnset{
				Comment: Bool(true),
			},
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER ALERT %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("with modify condition", func(t *testing.T) {
		modifyCondition := "SELECT * FROM FOO"
		opts := &AlterAlertOptions{
			name:            id,
			ModifyCondition: &[]string{modifyCondition},
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER ALERT %s MODIFY CONDITION EXISTS (%s)", id.FullyQualifiedName(), modifyCondition)
	})

	t.Run("with modify action", func(t *testing.T) {
		modifyAction := String("INSERT INTO FOO VALUES (1)")
		opts := &AlterAlertOptions{
			name:         id,
			ModifyAction: modifyAction,
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER ALERT %s MODIFY ACTION %s", id.FullyQualifiedName(), *modifyAction)
	})
}

func TestAlertDrop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &dropAlertOptions{}
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &dropAlertOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP ALERT %s", id.FullyQualifiedName())
	})
}

func TestAlertShow(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowAlertOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ALERTS")
	})

	t.Run("terse", func(t *testing.T) {
		opts := &ShowAlertOptions{Terse: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "SHOW TERSE ALERTS")
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowAlertOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ALERTS LIKE '%s'", id.Name())
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := &ShowAlertOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Account: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ALERTS LIKE '%s' IN ACCOUNT", id.Name())
	})

	t.Run("with like and in database", func(t *testing.T) {
		databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
		opts := &ShowAlertOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Database: databaseIdentifier,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ALERTS LIKE '%s' IN DATABASE %s", id.Name(), databaseIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in schema", func(t *testing.T) {
		schemaIdentifier := NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())
		opts := &ShowAlertOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Schema: schemaIdentifier,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ALERTS LIKE '%s' IN SCHEMA %s", id.Name(), schemaIdentifier.FullyQualifiedName())
	})

	t.Run("with 'starts with'", func(t *testing.T) {
		opts := &ShowAlertOptions{
			StartsWith: String("FOO"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ALERTS STARTS WITH 'FOO'")
	})

	t.Run("with limit", func(t *testing.T) {
		opts := &ShowAlertOptions{
			Limit: Int(10),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW ALERTS LIMIT 10")
	})
}

func TestAlertDescribe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &describeAlertOptions{}
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &describeAlertOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE ALERT %s", id.FullyQualifiedName())
	})
}
