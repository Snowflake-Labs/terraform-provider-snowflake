package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Users(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	randomPrefix := random.AlphaN(6)

	user, userCleanup := testClientHelper().User.CreateUserWithPrefix(t, randomPrefix+"_")
	t.Cleanup(userCleanup)

	user2, user2Cleanup := testClientHelper().User.CreateUserWithPrefix(t, randomPrefix)
	t.Cleanup(user2Cleanup)

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)

	t.Run("create: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToUpper(random.AlphaN(6))
		tagValue := random.String()
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := random.Password()
		loginName := random.String()

		opts := &sdk.CreateUserOptions{
			OrReplace: sdk.Bool(true),
			ObjectProperties: &sdk.UserObjectProperties{
				Password:    &password,
				LoginName:   &loginName,
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
			},
			ObjectParameters: &sdk.UserObjectParameters{
				EnableUnredactedQuerySyntaxError: sdk.Bool(true),
			},
			SessionParameters: &sdk.SessionParameters{
				Autocommit: sdk.Bool(true),
			},
			With: sdk.Bool(true),
			Tags: tags,
		}
		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)
		assert.Equal(t, defaultRole, userDetails.DefaultRole.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), user.Name)
	})

	t.Run("create: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		tagValue := random.String()
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := random.Password()
		loginName := random.String()

		opts := &sdk.CreateUserOptions{
			IfNotExists: sdk.Bool(true),
			ObjectProperties: &sdk.UserObjectProperties{
				Password:  &password,
				LoginName: &loginName,
			},
			ObjectParameters: &sdk.UserObjectParameters{
				EnableUnredactedQuerySyntaxError: sdk.Bool(true),
			},
			SessionParameters: &sdk.SessionParameters{
				Autocommit: sdk.Bool(true),
			},
			With: sdk.Bool(true),
			Tags: tags,
		}
		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), user.Name)
	})

	t.Run("create: no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(id.Name()), userDetails.LoginName.Value)
		assert.Empty(t, userDetails.Password.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), user.Name)
	})

	t.Run("create: default role with hyphen", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToUpper(random.AlphaN(4) + "-" + random.AlphaN(4))

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		createdUser, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, defaultRole, createdUser.DefaultRole)
	})

	t.Run("create: default role in lowercase", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToLower(random.AlphaN(6))

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		createdUser, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, defaultRole, createdUser.DefaultRole)
	})

	t.Run("create: other params with hyphen and mixed cases", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		randomWithHyphenAndMixedCase := strings.ToUpper(random.AlphaN(4)) + "-" + strings.ToLower(random.AlphaN(4))
		var namespaceId sdk.ObjectIdentifier = sdk.NewDatabaseObjectIdentifier(randomWithHyphenAndMixedCase, randomWithHyphenAndMixedCase)

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				LoginName:        sdk.String(randomWithHyphenAndMixedCase),
				DisplayName:      sdk.String(randomWithHyphenAndMixedCase),
				FirstName:        sdk.String(randomWithHyphenAndMixedCase),
				MiddleName:       sdk.String(randomWithHyphenAndMixedCase),
				LastName:         sdk.String(randomWithHyphenAndMixedCase),
				DefaultWarehouse: sdk.Pointer(sdk.NewAccountObjectIdentifier(randomWithHyphenAndMixedCase)),
				DefaultNamespace: sdk.Pointer(namespaceId),
				DefaultRole:      sdk.Pointer(sdk.NewAccountObjectIdentifier(randomWithHyphenAndMixedCase)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		createdUser, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		// login name is always case-insensitive
		assert.Equal(t, strings.ToUpper(randomWithHyphenAndMixedCase), createdUser.LoginName)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.DisplayName)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.FirstName)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.LastName)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.DefaultWarehouse)
		assert.Equal(t, randomWithHyphenAndMixedCase+"."+randomWithHyphenAndMixedCase, createdUser.DefaultNamespace)
		assert.Equal(t, randomWithHyphenAndMixedCase, createdUser.DefaultRole)

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.MiddleName.Value)
		// login name is always case-insensitive
		assert.Equal(t, strings.ToUpper(randomWithHyphenAndMixedCase), userDetails.LoginName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DisplayName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.FirstName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.LastName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DefaultWarehouse.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase+"."+randomWithHyphenAndMixedCase, userDetails.DefaultNamespace.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DefaultRole.Value)
	})

	// TODO: add tests for alter

	t.Run("create: with all parameters set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		opts := &sdk.CreateUserOptions{
			SessionParameters: &sdk.SessionParameters{
				AbortDetachedQuery:                       sdk.Bool(true),
				Autocommit:                               sdk.Bool(false),
				BinaryInputFormat:                        sdk.Pointer(sdk.BinaryInputFormatUTF8),
				BinaryOutputFormat:                       sdk.Pointer(sdk.BinaryOutputFormatBase64),
				ClientMemoryLimit:                        sdk.Int(1024),
				ClientMetadataRequestUseConnectionCtx:    sdk.Bool(true),
				ClientPrefetchThreads:                    sdk.Int(2),
				ClientResultChunkSize:                    sdk.Int(48),
				ClientResultColumnCaseInsensitive:        sdk.Bool(true),
				ClientSessionKeepAlive:                   sdk.Bool(true),
				ClientSessionKeepAliveHeartbeatFrequency: sdk.Int(2400),
				ClientTimestampTypeMapping:               sdk.Pointer(sdk.ClientTimestampTypeMappingNtz),
				DateInputFormat:                          sdk.String("YYYY-MM-DD"),
				DateOutputFormat:                         sdk.String("YY-MM-DD"),
				EnableUnloadPhysicalTypeOptimization:     sdk.Bool(false),
				ErrorOnNondeterministicMerge:             sdk.Bool(false),
				ErrorOnNondeterministicUpdate:            sdk.Bool(true),
				GeographyOutputFormat:                    sdk.Pointer(sdk.GeographyOutputFormatWKB),
				GeometryOutputFormat:                     sdk.Pointer(sdk.GeometryOutputFormatWKB),
				JdbcTreatDecimalAsInt:                    sdk.Bool(false),
				JdbcTreatTimestampNtzAsUtc:               sdk.Bool(true),
				JdbcUseSessionTimezone:                   sdk.Bool(false),
				JSONIndent:                               sdk.Int(4),
				LockTimeout:                              sdk.Int(21222),
				LogLevel:                                 sdk.Pointer(sdk.LogLevelError),
				MultiStatementCount:                      sdk.Int(0),
				NoorderSequenceAsDefault:                 sdk.Bool(false),
				OdbcTreatDecimalAsInt:                    sdk.Bool(true),
				QueryTag:                                 sdk.String("some_tag"),
				QuotedIdentifiersIgnoreCase:              sdk.Bool(true),
				RowsPerResultset:                         sdk.Int(2),
				S3StageVpceDnsName:                       sdk.String("vpce-some_dns-vpce.amazonaws.com"),
				SearchPath:                               sdk.String("$public, $current"),
				SimulatedDataSharingConsumer:             sdk.String("some_consumer"),
				StatementQueuedTimeoutInSeconds:          sdk.Int(10),
				StatementTimeoutInSeconds:                sdk.Int(10),
				StrictJSONOutput:                         sdk.Bool(true),
				TimestampDayIsAlways24h:                  sdk.Bool(true),
				TimestampInputFormat:                     sdk.String("YYYY-MM-DD"),
				TimestampLTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
				TimestampNTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
				TimestampOutputFormat:                    sdk.String("YYYY-MM-DD HH24:MI:SS"),
				TimestampTypeMapping:                     sdk.Pointer(sdk.TimestampTypeMappingLtz),
				TimestampTZOutputFormat:                  sdk.String("YYYY-MM-DD HH24:MI:SS"),
				Timezone:                                 sdk.String("Europe/Warsaw"),
				TimeInputFormat:                          sdk.String("HH24:MI"),
				TimeOutputFormat:                         sdk.String("HH24:MI"),
				TraceLevel:                               sdk.Pointer(sdk.TraceLevelOnEvent),
				TransactionAbortOnError:                  sdk.Bool(true),
				TransactionDefaultIsolationLevel:         sdk.Pointer(sdk.TransactionDefaultIsolationLevelReadCommitted),
				TwoDigitCenturyStart:                     sdk.Int(1980),
				UnsupportedDDLAction:                     sdk.Pointer(sdk.UnsupportedDDLActionFail),
				UseCachedResult:                          sdk.Bool(false),
				WeekOfYearPolicy:                         sdk.Int(1),
				WeekStart:                                sdk.Int(1),
			},
			ObjectParameters: &sdk.UserObjectParameters{
				EnableUnredactedQuerySyntaxError: sdk.Bool(true),
				NetworkPolicy:                    sdk.Pointer(networkPolicy.ID()),
				PreventUnloadToInternalStages:    sdk.Bool(true),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		createdUser, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), createdUser.Name)

		parameters := testClientHelper().Parameter.ShowUserParameters(t, id)

		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterAbortDetachedQuery).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterAutocommit).Value)
		assert.Equal(t, string(sdk.BinaryInputFormatUTF8), helpers.FindParameter(t, parameters, sdk.UserParameterBinaryInputFormat).Value)
		assert.Equal(t, string(sdk.BinaryOutputFormatBase64), helpers.FindParameter(t, parameters, sdk.UserParameterBinaryOutputFormat).Value)
		// TODO: check zero
		assert.Equal(t, "1024", helpers.FindParameter(t, parameters, sdk.UserParameterClientMemoryLimit).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterClientMetadataRequestUseConnectionCtx).Value)
		assert.Equal(t, "2", helpers.FindParameter(t, parameters, sdk.UserParameterClientPrefetchThreads).Value)
		assert.Equal(t, "48", helpers.FindParameter(t, parameters, sdk.UserParameterClientResultChunkSize).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterClientResultColumnCaseInsensitive).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterClientSessionKeepAlive).Value)
		assert.Equal(t, "2400", helpers.FindParameter(t, parameters, sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency).Value)
		assert.Equal(t, string(sdk.ClientTimestampTypeMappingNtz), helpers.FindParameter(t, parameters, sdk.UserParameterClientTimestampTypeMapping).Value)
		assert.Equal(t, "YYYY-MM-DD", helpers.FindParameter(t, parameters, sdk.UserParameterDateInputFormat).Value)
		assert.Equal(t, "YY-MM-DD", helpers.FindParameter(t, parameters, sdk.UserParameterDateOutputFormat).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterEnableUnloadPhysicalTypeOptimization).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterErrorOnNondeterministicMerge).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterErrorOnNondeterministicUpdate).Value)
		assert.Equal(t, string(sdk.GeographyOutputFormatWKB), helpers.FindParameter(t, parameters, sdk.UserParameterGeographyOutputFormat).Value)
		assert.Equal(t, string(sdk.GeometryOutputFormatWKB), helpers.FindParameter(t, parameters, sdk.UserParameterGeometryOutputFormat).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterJdbcTreatDecimalAsInt).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterJdbcTreatTimestampNtzAsUtc).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterJdbcUseSessionTimezone).Value)
		assert.Equal(t, "4", helpers.FindParameter(t, parameters, sdk.UserParameterJsonIndent).Value)
		assert.Equal(t, "21222", helpers.FindParameter(t, parameters, sdk.UserParameterLockTimeout).Value)
		assert.Equal(t, string(sdk.LogLevelError), helpers.FindParameter(t, parameters, sdk.UserParameterLogLevel).Value)
		assert.Equal(t, "0", helpers.FindParameter(t, parameters, sdk.UserParameterMultiStatementCount).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterNoorderSequenceAsDefault).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterOdbcTreatDecimalAsInt).Value)
		assert.Equal(t, "some_tag", helpers.FindParameter(t, parameters, sdk.UserParameterQueryTag).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterQuotedIdentifiersIgnoreCase).Value)
		assert.Equal(t, "2", helpers.FindParameter(t, parameters, sdk.UserParameterRowsPerResultset).Value)
		assert.Equal(t, "vpce-some_dns-vpce.amazonaws.com", helpers.FindParameter(t, parameters, sdk.UserParameterS3StageVpceDnsName).Value)
		assert.Equal(t, "$public, $current", helpers.FindParameter(t, parameters, sdk.UserParameterSearchPath).Value)
		assert.Equal(t, "some_consumer", helpers.FindParameter(t, parameters, sdk.UserParameterSimulatedDataSharingConsumer).Value)
		assert.Equal(t, "10", helpers.FindParameter(t, parameters, sdk.UserParameterStatementQueuedTimeoutInSeconds).Value)
		assert.Equal(t, "10", helpers.FindParameter(t, parameters, sdk.UserParameterStatementTimeoutInSeconds).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterStrictJsonOutput).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampDayIsAlways24h).Value)
		assert.Equal(t, "YYYY-MM-DD", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampInputFormat).Value)
		assert.Equal(t, "YYYY-MM-DD HH24:MI:SS", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampLtzOutputFormat).Value)
		assert.Equal(t, "YYYY-MM-DD HH24:MI:SS", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampNtzOutputFormat).Value)
		assert.Equal(t, "YYYY-MM-DD HH24:MI:SS", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampOutputFormat).Value)
		assert.Equal(t, string(sdk.TimestampTypeMappingLtz), helpers.FindParameter(t, parameters, sdk.UserParameterTimestampTypeMapping).Value)
		assert.Equal(t, "YYYY-MM-DD HH24:MI:SS", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampTzOutputFormat).Value)
		assert.Equal(t, "Europe/Warsaw", helpers.FindParameter(t, parameters, sdk.UserParameterTimezone).Value)
		assert.Equal(t, "HH24:MI", helpers.FindParameter(t, parameters, sdk.UserParameterTimeInputFormat).Value)
		assert.Equal(t, "HH24:MI", helpers.FindParameter(t, parameters, sdk.UserParameterTimeOutputFormat).Value)
		assert.Equal(t, string(sdk.TraceLevelOnEvent), helpers.FindParameter(t, parameters, sdk.UserParameterTraceLevel).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterTransactionAbortOnError).Value)
		assert.Equal(t, string(sdk.TransactionDefaultIsolationLevelReadCommitted), helpers.FindParameter(t, parameters, sdk.UserParameterTransactionDefaultIsolationLevel).Value)
		assert.Equal(t, "1980", helpers.FindParameter(t, parameters, sdk.UserParameterTwoDigitCenturyStart).Value)
		assert.Equal(t, string(sdk.UnsupportedDDLActionFail), helpers.FindParameter(t, parameters, sdk.UserParameterUnsupportedDdlAction).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterUseCachedResult).Value)
		assert.Equal(t, "1", helpers.FindParameter(t, parameters, sdk.UserParameterWeekOfYearPolicy).Value)
		assert.Equal(t, "1", helpers.FindParameter(t, parameters, sdk.UserParameterWeekStart).Value)

		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterEnableUnredactedQuerySyntaxError).Value)
		assert.Equal(t, networkPolicy.ID().Name(), helpers.FindParameter(t, parameters, sdk.UserParameterNetworkPolicy).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterPreventUnloadToInternalStages).Value)
	})

	t.Run("create: with all parameters default", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		createdUser, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), createdUser.Name)

		parameters := testClientHelper().Parameter.ShowUserParameters(t, id)

		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterAbortDetachedQuery).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterAutocommit).Value)
		assert.Equal(t, string(sdk.BinaryInputFormatHex), helpers.FindParameter(t, parameters, sdk.UserParameterBinaryInputFormat).Value)
		assert.Equal(t, string(sdk.BinaryOutputFormatHex), helpers.FindParameter(t, parameters, sdk.UserParameterBinaryOutputFormat).Value)
		assert.Equal(t, "1536", helpers.FindParameter(t, parameters, sdk.UserParameterClientMemoryLimit).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterClientMetadataRequestUseConnectionCtx).Value)
		assert.Equal(t, "4", helpers.FindParameter(t, parameters, sdk.UserParameterClientPrefetchThreads).Value)
		assert.Equal(t, "160", helpers.FindParameter(t, parameters, sdk.UserParameterClientResultChunkSize).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterClientResultColumnCaseInsensitive).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterClientSessionKeepAlive).Value)
		assert.Equal(t, "3600", helpers.FindParameter(t, parameters, sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency).Value)
		assert.Equal(t, string(sdk.ClientTimestampTypeMappingLtz), helpers.FindParameter(t, parameters, sdk.UserParameterClientTimestampTypeMapping).Value)
		assert.Equal(t, "AUTO", helpers.FindParameter(t, parameters, sdk.UserParameterDateInputFormat).Value)
		assert.Equal(t, "YYYY-MM-DD", helpers.FindParameter(t, parameters, sdk.UserParameterDateOutputFormat).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterEnableUnloadPhysicalTypeOptimization).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterErrorOnNondeterministicMerge).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterErrorOnNondeterministicUpdate).Value)
		assert.Equal(t, string(sdk.GeographyOutputFormatGeoJSON), helpers.FindParameter(t, parameters, sdk.UserParameterGeographyOutputFormat).Value)
		assert.Equal(t, string(sdk.GeometryOutputFormatGeoJSON), helpers.FindParameter(t, parameters, sdk.UserParameterGeometryOutputFormat).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterJdbcTreatDecimalAsInt).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterJdbcTreatTimestampNtzAsUtc).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterJdbcUseSessionTimezone).Value)
		assert.Equal(t, "2", helpers.FindParameter(t, parameters, sdk.UserParameterJsonIndent).Value)
		assert.Equal(t, "43200", helpers.FindParameter(t, parameters, sdk.UserParameterLockTimeout).Value)
		assert.Equal(t, string(sdk.LogLevelOff), helpers.FindParameter(t, parameters, sdk.UserParameterLogLevel).Value)
		assert.Equal(t, "1", helpers.FindParameter(t, parameters, sdk.UserParameterMultiStatementCount).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterNoorderSequenceAsDefault).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterOdbcTreatDecimalAsInt).Value)
		assert.Equal(t, "", helpers.FindParameter(t, parameters, sdk.UserParameterQueryTag).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterQuotedIdentifiersIgnoreCase).Value)
		assert.Equal(t, "0", helpers.FindParameter(t, parameters, sdk.UserParameterRowsPerResultset).Value)
		assert.Equal(t, "", helpers.FindParameter(t, parameters, sdk.UserParameterS3StageVpceDnsName).Value)
		assert.Equal(t, "$current, $public", helpers.FindParameter(t, parameters, sdk.UserParameterSearchPath).Value)
		assert.Equal(t, "", helpers.FindParameter(t, parameters, sdk.UserParameterSimulatedDataSharingConsumer).Value)
		assert.Equal(t, "0", helpers.FindParameter(t, parameters, sdk.UserParameterStatementQueuedTimeoutInSeconds).Value)
		assert.Equal(t, "172800", helpers.FindParameter(t, parameters, sdk.UserParameterStatementTimeoutInSeconds).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterStrictJsonOutput).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampDayIsAlways24h).Value)
		assert.Equal(t, "AUTO", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampInputFormat).Value)
		assert.Equal(t, "", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampLtzOutputFormat).Value)
		assert.Equal(t, "YYYY-MM-DD HH24:MI:SS.FF3", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampNtzOutputFormat).Value)
		assert.Equal(t, "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampOutputFormat).Value)
		assert.Equal(t, string(sdk.TimestampTypeMappingNtz), helpers.FindParameter(t, parameters, sdk.UserParameterTimestampTypeMapping).Value)
		assert.Equal(t, "", helpers.FindParameter(t, parameters, sdk.UserParameterTimestampTzOutputFormat).Value)
		assert.Equal(t, "America/Los_Angeles", helpers.FindParameter(t, parameters, sdk.UserParameterTimezone).Value)
		assert.Equal(t, "AUTO", helpers.FindParameter(t, parameters, sdk.UserParameterTimeInputFormat).Value)
		assert.Equal(t, "HH24:MI:SS", helpers.FindParameter(t, parameters, sdk.UserParameterTimeOutputFormat).Value)
		assert.Equal(t, string(sdk.TraceLevelOff), helpers.FindParameter(t, parameters, sdk.UserParameterTraceLevel).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterTransactionAbortOnError).Value)
		assert.Equal(t, string(sdk.TransactionDefaultIsolationLevelReadCommitted), helpers.FindParameter(t, parameters, sdk.UserParameterTransactionDefaultIsolationLevel).Value)
		assert.Equal(t, "1970", helpers.FindParameter(t, parameters, sdk.UserParameterTwoDigitCenturyStart).Value)
		// lowercase by default in Snowflake
		assert.Equal(t, strings.ToLower(string(sdk.UnsupportedDDLActionIgnore)), helpers.FindParameter(t, parameters, sdk.UserParameterUnsupportedDdlAction).Value)
		assert.Equal(t, "true", helpers.FindParameter(t, parameters, sdk.UserParameterUseCachedResult).Value)
		assert.Equal(t, "0", helpers.FindParameter(t, parameters, sdk.UserParameterWeekOfYearPolicy).Value)
		assert.Equal(t, "0", helpers.FindParameter(t, parameters, sdk.UserParameterWeekStart).Value)

		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterEnableUnredactedQuerySyntaxError).Value)
		assert.Equal(t, "", helpers.FindParameter(t, parameters, sdk.UserParameterNetworkPolicy).Value)
		assert.Equal(t, "false", helpers.FindParameter(t, parameters, sdk.UserParameterPreventUnloadToInternalStages).Value)
	})

	t.Run("describe: when user exists", func(t *testing.T) {
		userDetails, err := client.Users.Describe(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.Name, userDetails.Name.Value)
	})

	t.Run("describe: when user does not exist", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		_, err := client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when user exists", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		id := user.ID()
		err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{})
		require.NoError(t, err)
		_, err = client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when user does not exist", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{})
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show: with like options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(user.Name),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Equal(t, 1, len(users))
	})

	t.Run("show: with starts with options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			StartsWith: sdk.String(randomPrefix),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Contains(t, users, *user2)
		assert.Equal(t, 2, len(users))
	})

	t.Run("show: with starts with, limit and from options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Limit:      sdk.Int(10),
			From:       sdk.String(randomPrefix + "_"),
			StartsWith: sdk.String(randomPrefix),
		}

		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Equal(t, 1, len(users))
	})

	t.Run("show: search for a non-existent user", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(NonExistingAccountObjectIdentifier.Name()),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(users))
	})

	t.Run("show: limit the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Limit: sdk.Int(1),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})
}
