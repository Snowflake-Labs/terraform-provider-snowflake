package testint

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1827310]: use generated config for these tests
func TestInt_Client_NewClient(t *testing.T) {
	t.Run("with default config", func(t *testing.T) {
		config := sdk.DefaultConfig()
		_, err := sdk.NewClient(config)
		require.NoError(t, err)
	})

	t.Run("with missing config", func(t *testing.T) {
		dir, err := os.UserHomeDir()
		require.NoError(t, err)
		t.Setenv(snowflakeenvs.ConfigPath, dir)

		config := sdk.DefaultConfig()
		_, err = sdk.NewClient(config)
		require.ErrorContains(t, err, "260000: account is empty")
	})

	t.Run("with incorrect config", func(t *testing.T) {
		tmpServiceUser := testClientHelper().SetUpTemporaryServiceUser(t)
		tmpServiceUserConfig := testClientHelper().TempIncorrectTomlConfigForServiceUser(t, tmpServiceUser)

		t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)

		config, err := sdk.ProfileConfig(tmpServiceUserConfig.Profile)
		require.NoError(t, err)
		require.NotNil(t, config)

		_, err = sdk.NewClient(config)
		require.ErrorContains(t, err, "JWT token is invalid")
	})

	t.Run("with missing config - should not care about correct env variables", func(t *testing.T) {
		config, err := sdk.ProfileConfig(testprofiles.Default)
		require.NoError(t, err)
		require.NotNil(t, config)

		account := config.Account
		parts := strings.Split(account, "-")
		t.Setenv(snowflakeenvs.OrganizationName, parts[0])
		t.Setenv(snowflakeenvs.AccountName, parts[1])

		dir, err := os.UserHomeDir()
		require.NoError(t, err)
		t.Setenv(snowflakeenvs.ConfigPath, dir)

		config = sdk.DefaultConfig()
		_, err = sdk.NewClient(config)
		require.ErrorContains(t, err, "260000: account is empty")
	})

	t.Run("registers snowflake-instrumented driver", func(t *testing.T) {
		config := sdk.DefaultConfig()
		_, err := sdk.NewClient(config)
		require.NoError(t, err)

		assert.ElementsMatch(t, sql.Drivers(), []string{"snowflake-instrumented", "snowflake"})
	})
}

func TestInt_Client_AdditionalMetadata(t *testing.T) {
	client := testClient(t)
	metadata := tracking.Metadata{SchemaVersion: "1", Version: "v1.13.1002-rc-test", Resource: resources.Database.String(), Operation: tracking.CreateOperation}

	assertQueryMetadata := func(t *testing.T, queryId string) {
		t.Helper()
		queryText := testClientHelper().InformationSchema.GetQueryHistoryByQueryId(t, 20, queryId).QueryText
		parsedMetadata, err := tracking.ParseMetadata(queryText)
		require.NoError(t, err)
		require.Equal(t, metadata, parsedMetadata)
	}

	t.Run("query one", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		row := struct {
			One int `db:"ONE"`
		}{}
		err := client.QueryOneForTests(ctx, &row, "SELECT 1 AS ONE")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})

	t.Run("query", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		var rows []struct {
			One int `db:"ONE"`
		}
		err := client.QueryForTests(ctx, &rows, "SELECT 1 AS ONE")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})

	t.Run("exec", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		_, err := client.ExecForTests(ctx, "SELECT 1")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})
}
