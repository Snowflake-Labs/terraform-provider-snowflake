package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestAlertGenCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("with complete options", func(t *testing.T) {
		newComment := randomString(t)
		warehouse := AccountObjectIdentifier{"warehouse"}
		existsCondition := "SELECT 1"
		schedule := "1 minute"
		action := "INSERT INTO FOO VALUES (1)"

		opts := &AlertCreateOptionsGen{
			name:      id,
			Warehouse: warehouse.name,
			Schedule:  schedule,
			If:        If{Exists{existsCondition}},
			Then:      Then{action},
			Comment:   String(newComment),
		}

		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`CREATE ALERT %s WAREHOUSE = "%s" SCHEDULE = '%s' COMMENT = '%s' IF (EXISTS (%s)) THEN %s`, id.FullyQualifiedName(), warehouse.name, schedule, newComment, existsCondition, action)
		assert.Equal(t, expected, actual)
	})
}
