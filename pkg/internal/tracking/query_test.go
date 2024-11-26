package tracking

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/stretchr/testify/require"
)

func TestAppendMetadata(t *testing.T) {
	metadata := NewTestMetadata("123", resources.Account, CreateOperation)
	sql := "SELECT 1"

	bytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	expectedSql := fmt.Sprintf("%s --%s %s", sql, MetadataPrefix, string(bytes))

	newSql, err := AppendMetadata(sql, metadata)
	require.NoError(t, err)
	require.Equal(t, expectedSql, newSql)
}

func TestParseMetadata(t *testing.T) {
	metadata := NewTestMetadata("123", resources.Account, CreateOperation)
	bytes, err := json.Marshal(metadata)
	require.NoError(t, err)
	sql := fmt.Sprintf("SELECT 1 --%s %s", MetadataPrefix, string(bytes))

	parsedMetadata, err := ParseMetadata(sql)
	require.NoError(t, err)
	require.Equal(t, metadata, parsedMetadata)
}

func TestParseInvalidMetadataKeys(t *testing.T) {
	sql := fmt.Sprintf(`SELECT 1 --%s {"key": "value"}`, MetadataPrefix)

	parsedMetadata, err := ParseMetadata(sql)
	require.ErrorContains(t, err, "schema version for metadata should not be empty")
	require.ErrorContains(t, err, "provider version for metadata should not be empty")
	require.ErrorContains(t, err, "either resource or data source name for metadata should be specified")
	require.ErrorContains(t, err, "operation for metadata should not be empty")
	require.Equal(t, Metadata{}, parsedMetadata)
}

func TestParseInvalidMetadataJson(t *testing.T) {
	sql := fmt.Sprintf(`SELECT 1 --%s "key": "value"`, MetadataPrefix)

	parsedMetadata, err := ParseMetadata(sql)
	require.ErrorContains(t, err, "failed to unmarshal metadata from sql")
	require.Equal(t, Metadata{}, parsedMetadata)
}

func TestParseMetadataFromInvalidSqlCommentPrefix(t *testing.T) {
	metadata := NewTestMetadata("123", resources.Account, CreateOperation)
	sql := "SELECT 1"

	bytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	parsedMetadata, err := ParseMetadata(fmt.Sprintf("%s --invalid_prefix %s", sql, string(bytes)))
	require.ErrorContains(t, err, "failed to parse metadata from sql")
	require.Equal(t, Metadata{}, parsedMetadata)
}
