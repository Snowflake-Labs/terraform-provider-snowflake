package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestAlertPolicyCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("with complete options", func(t *testing.T) {
		newComment := randomString(t)
		warehouse := AccountObjectIdentifier{"warehouse"}
		existsCondition := "SELECT 1"
		condition := AlertCondition{existsCondition}
		schedule := "1 minute"
		action := "INSERT INTO FOO VALUES (1)"

		opts := &CreateAlertOptions{
			name:      id,
			warehouse: warehouse,
			schedule:  schedule,
			condition: condition,
			action:    action,
			Comment:   String(newComment),
		}

		err := opts.validate()
		assert.NoError(t, err)
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`CREATE ALERT %s WAREHOUSE = "%s" SCHEDULE = '%s' COMMENT = '%s' IF(EXISTS(%s)) THEN %s`, id.FullyQualifiedName(), warehouse.name, schedule, newComment, existsCondition, action)
		assert.Equal(t, expected, actual)
	})
}

func TestAlertAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("fail without alter action specified", func(t *testing.T) {
		opts := &AlterAlertOptions{
			name: id,
		}
		err := opts.validate()
		assert.Error(t, err)
	})

	t.Run("fail when 2 alter actions specified", func(t *testing.T) {
		newComment := randomString(t)
		opts := &AlterAlertOptions{
			name:      id,
			Operation: &Resume,
			Set: &AlertSet{
				Comment: String(newComment),
			},
		}
		err := opts.validate()
		assert.Error(t, err)
	})

	t.Run("with resume", func(t *testing.T) {
		opts := &AlterAlertOptions{
			name:      id,
			Operation: &Resume,
		}

		err := opts.validate()
		assert.NoError(t, err)
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER ALERT %s RESUME", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with suspend", func(t *testing.T) {
		opts := &AlterAlertOptions{
			name:      id,
			Operation: &Suspend,
		}

		err := opts.validate()
		assert.NoError(t, err)
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER ALERT %s SUSPEND", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with set", func(t *testing.T) {
		newComment := randomString(t)
		opts := &AlterAlertOptions{
			name: id,
			Set: &AlertSet{
				Comment: String(newComment),
			},
		}
		err := opts.validate()
		assert.NoError(t, err)
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER ALERT %s SET COMMENT = '%s'", id.FullyQualifiedName(), newComment)
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &AlterAlertOptions{
			name: id,
			Unset: &AlertUnset{
				Comment: Bool(true),
			},
		}
		err := opts.validate()
		assert.NoError(t, err)
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER ALERT %s UNSET COMMENT", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with modify condition", func(t *testing.T) {
		modifyCondition := String("SELECT * FROM FOO")
		opts := &AlterAlertOptions{
			name:            id,
			ModifyCondition: modifyCondition,
		}
		err := opts.validate()
		assert.NoError(t, err)
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER ALERT %s MODIFY CONDITION EXISTS(%s)", id.FullyQualifiedName(), *modifyCondition)
		assert.Equal(t, expected, actual)
	})
	t.Run("with modify action", func(t *testing.T) {
		modifyAction := String("INSERT INTO FOO VALUES (1)")
		opts := &AlterAlertOptions{
			name:         id,
			ModifyAction: modifyAction,
		}
		err := opts.validate()
		assert.NoError(t, err)
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER ALERT %s MODIFY ACTION %s", id.FullyQualifiedName(), *modifyAction)
		assert.Equal(t, expected, actual)
	})
}

func TestAlertDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &DropAlertOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "DROP ALERT"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &DropAlertOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("DROP ALERT %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestAlertShow(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowAlertOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW ALERTS"
		assert.Equal(t, expected, actual)
	})

	t.Run("terse", func(t *testing.T) {
		opts := &ShowAlertOptions{Terse: Bool(true)}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW TERSE ALERTS"
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowAlertOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW ALERTS LIKE '%s'", id.Name())
		assert.Equal(t, expected, actual)
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
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW ALERTS LIKE '%s' IN ACCOUNT", id.Name())
		assert.Equal(t, expected, actual)
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
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW ALERTS LIKE '%s' IN DATABASE %s", id.Name(), databaseIdentifier.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and in schema", func(t *testing.T) {
		schemaIdentifier := NewSchemaIdentifier(id.DatabaseName(), id.SchemaName())
		opts := &ShowAlertOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Schema: schemaIdentifier,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW ALERTS LIKE '%s' IN SCHEMA %s", id.Name(), schemaIdentifier.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with 'starts with'", func(t *testing.T) {
		opts := &ShowAlertOptions{
			StartsWith: String("FOO"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW ALERTS STARTS WITH 'FOO'"
		assert.Equal(t, expected, actual)
	})

	t.Run("with limit", func(t *testing.T) {
		opts := &ShowAlertOptions{
			Limit: Int(10),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW ALERTS LIMIT 10"
		assert.Equal(t, expected, actual)
	})
}

func TestAlertDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &describeAlertOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "DESCRIBE ALERT"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &describeAlertOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("DESCRIBE ALERT %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}
