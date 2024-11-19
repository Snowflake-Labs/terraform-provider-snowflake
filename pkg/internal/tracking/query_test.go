package tracking

import (
	"encoding/json"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAppendMetadataToSql(t *testing.T) {
	metadata := NewMetadata("123", resources.Account, CreateOperation)
	sql := "SELECT 1"

	bytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	expectedSql := fmt.Sprintf("%s --%s %s", sql, MetadataPrefix, string(bytes))

	newSql, err := AppendMetadataToSql(sql, metadata)
	require.NoError(t, err)
	require.Equal(t, expectedSql, newSql)
}

func TestParseMetadataFromSql(t *testing.T) {
	metadata := NewMetadata("123", resources.Account, CreateOperation)
	sql, err := AppendMetadataToSql("SELECT 1", metadata)
	require.NoError(t, err)

	parsedMetadata, err := ParseMetadataFromSql(sql)
	require.NoError(t, err)
	require.Equal(t, metadata, parsedMetadata)
}

func TestParseInvalidMetadataKeysFromSql(t *testing.T) {
	sql := fmt.Sprintf(`SELECT 1 --%s {"key": "value"}`, MetadataPrefix)

	parsedMetadata, err := ParseMetadataFromSql(sql)
	require.NoError(t, err)
	require.Equal(t, Metadata{}, parsedMetadata)
}

func TestParseInvalidMetadataJsonFromSql(t *testing.T) {
	sql := fmt.Sprintf(`SELECT 1 --%s "key": "value"`, MetadataPrefix)

	parsedMetadata, err := ParseMetadataFromSql(sql)
	require.ErrorContains(t, err, "failed to unmarshal metadata from sql")
	require.Equal(t, Metadata{}, parsedMetadata)
}

func TestParseMetadataFromInvalidSqlCommentPrefix(t *testing.T) {
	metadata := NewMetadata("123", resources.Account, CreateOperation)
	sql := "SELECT 1"

	bytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	parsedMetadata, err := ParseMetadataFromSql(fmt.Sprintf("%s --invalid_prefix %s", sql, string(bytes)))
	require.ErrorContains(t, err, "failed to parse metadata from sql")
	require.Equal(t, Metadata{}, parsedMetadata)
}
