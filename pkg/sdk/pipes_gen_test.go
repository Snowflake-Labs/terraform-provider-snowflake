package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestPipesGenCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("all optional", func(t *testing.T) {
		opts := &PipeCreateOptionsGen{
			name:              id,
			IfNotExists:       Bool(true),
			Auto_ingest:       Bool(true),
			Error_integration: String("some_error_integration"),
			Aws_sns_topic:     String("some aws sns topic"),
			Integration:       String("some integration"),
			Comment:           String("some comment"),
			As:                As{Copy_statement: "<copy_statement>"},
		}

		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`CREATE PIPE IF NOT EXISTS %s AUTO_INGEST = true ERROR_INTEGRATION = some_error_integration AWS_SNS_TOPIC = 'some aws sns topic' INTEGRATION = 'some integration' COMMENT = 'some comment' AS <copy_statement>`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}
