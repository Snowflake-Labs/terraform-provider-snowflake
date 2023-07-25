package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestAlertCreate1(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("with complete options", func(t *testing.T) {
		warehouse := AccountObjectIdentifier{"warehouse"}
		schedule := "1 MINUTE"

		opts := &AlertCreateOptions1{
			name:      id,
			Warehouse: warehouse.name,
			Schedule: Schedule{
				minute: &Minute{N: 1},
				cron:   nil,
			},
		}

		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`CREATE ALERT %s WAREHOUSE = "%s" SCHEDULE = '%s'`, id.FullyQualifiedName(), warehouse.name, schedule)
		assert.Equal(t, expected, actual)
	})
}
