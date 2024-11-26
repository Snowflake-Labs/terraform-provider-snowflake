package previewfeatures

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StringToFeature(t *testing.T) {
	type test struct {
		input string
		want  feature
	}

	valid := []test{
		// Case insensitive.
		{input: "SNOWFLAKE_CURRENT_ACCOUNT_DATASOURCE", want: CurrentAccountDatasource},

		// Supported Values.
		{input: "snowflake_current_account_datasource", want: CurrentAccountDatasource},
		{input: "snowflake_account_password_policy_attachment_resource", want: AccountPasswordPolicyAttachmentResource},
		{input: "snowflake_alert_resource", want: AlertResource},
		{input: "snowflake_alerts_datasource", want: AlertsDatasource},
		{input: "snowflake_api_integration_resource", want: ApiIntegrationResource},
		{input: "snowflake_cortex_search_service_resource", want: CortexSearchServiceResource},
		{input: "snowflake_cortex_search_services_datasource", want: CortexSearchServicesDatasource},
		{input: "snowflake_database_datasource", want: DatabaseDatasource},
		{input: "snowflake_database_role_datasource", want: DatabaseRoleDatasource},
		{input: "snowflake_dynamic_table_resource", want: DynamicTableResource},
		{input: "snowflake_dynamic_tables_datasource", want: DynamicTablesDatasource},
		{input: "snowflake_external_function_resource", want: ExternalFunctionResource},
		{input: "snowflake_external_functions_datasource", want: ExternalFunctionsDatasource},
		{input: "snowflake_external_table_resource", want: ExternalTableResource},
		{input: "snowflake_external_tables_datasource", want: ExternalTablesDatasource},
		{input: "snowflake_external_volume_resource", want: ExternalVolumeResource},
		{input: "snowflake_failover_group_resource", want: FailoverGroupResource},
		{input: "snowflake_failover_groups_datasource", want: FailoverGroupsDatasource},
		{input: "snowflake_file_format_resource", want: FileFormatResource},
		{input: "snowflake_file_formats_datasource", want: FileFormatsDatasource},
		{input: "snowflake_managed_account_resource", want: ManagedAccountResource},
		{input: "snowflake_materialized_view_resource", want: MaterializedViewResource},
		{input: "snowflake_materialized_views_datasource", want: MaterializedViewsDatasource},
		{input: "snowflake_network_policy_attachment_resource", want: NetworkPolicyAttachmentResource},
		{input: "snowflake_network_rule_resource", want: NetworkRuleResource},
		{input: "snowflake_email_notification_integration_resource", want: EmailNotificationIntegrationResource},
		{input: "snowflake_notification_integration_resource", want: NotificationIntegrationResource},
		{input: "snowflake_object_parameter_resource", want: ObjectParameterResource},
		{input: "snowflake_password_policy_resource", want: PasswordPolicyResource},
		{input: "snowflake_pipe_resource", want: PipeResource},
		{input: "snowflake_pipes_datasource", want: PipesDatasource},
		{input: "snowflake_current_role_datasource", want: CurrentRoleDatasource},
		{input: "snowflake_sequence_resource", want: SequenceResource},
		{input: "snowflake_sequences_datasource", want: SequencesDatasource},
		{input: "snowflake_share_resource", want: ShareResource},
		{input: "snowflake_shares_datasource", want: SharesDatasource},
		{input: "snowflake_parameters_datasource", want: ParametersDatasource},
		{input: "snowflake_stage_resource", want: StageResource},
		{input: "snowflake_stages_datasource", want: StagesDatasource},
		{input: "snowflake_storage_integration_resource", want: StorageIntegrationResource},
		{input: "snowflake_storage_integrations_datasource", want: StorageIntegrationsDatasource},
		{input: "snowflake_system_generate_scim_access_token_datasource", want: SystemGenerateSCIMAccessTokenDatasource},
		{input: "snowflake_system_get_aws_sns_iam_policy_datasource", want: SystemGetAWSSNSIAMPolicyDatasource},
		{input: "snowflake_system_get_privatelink_config_datasource", want: SystemGetPrivateLinkConfigDatasource},
		{input: "snowflake_system_get_snowflake_platform_info_datasource", want: SystemGetSnowflakePlatformInfoDatasource},
		{input: "snowflake_table_column_masking_policy_application_resource", want: TableColumnMaskingPolicyApplicationResource},
		{input: "snowflake_table_constraint_resource", want: TableConstraintResource},
		{input: "snowflake_user_authentication_policy_attachment_resource", want: UserAuthenticationPolicyAttachmentResource},
		{input: "snowflake_user_public_keys_resource", want: UserPublicKeysResource},
		{input: "snowflake_user_password_policy_attachment_resource", want: UserPasswordPolicyAttachmentResource},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := StringToFeature(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := StringToFeature(tc.input)
			require.ErrorContains(t, err, "invalid feature")
		})
	}
}
