package resources

type ResourceName string

const (
	Account                                                ResourceName = "snowflake_account"
	AccountRole                                            ResourceName = "snowflake_account_role"
	Alert                                                  ResourceName = "snowflake_alert"
	ApiAuthenticationIntegrationWithAuthorizationCodeGrant ResourceName = "snowflake_api_authentication_integration_with_authorization_code_grant"
	ApiAuthenticationIntegrationWithClientCredentials      ResourceName = "snowflake_api_authentication_integration_with_client_credentials"
	ApiAuthenticationIntegrationWithJwtBearer              ResourceName = "snowflake_api_authentication_integration_with_jwt_bearer"
	ApiIntegration                                         ResourceName = "snowflake_api_integration"
	AuthenticationPolicy                                   ResourceName = "snowflake_authentication_policy"
	CortexSearchService                                    ResourceName = "snowflake_cortex_search_service"
	DatabaseOld                                            ResourceName = "snowflake_database_old"
	Database                                               ResourceName = "snowflake_database"
	DatabaseRole                                           ResourceName = "snowflake_database_role"
	DynamicTable                                           ResourceName = "snowflake_dynamic_table"
	EmailNotificationIntegration                           ResourceName = "snowflake_email_notification_integration"
	ExternalFunction                                       ResourceName = "snowflake_external_function"
	ExternalTable                                          ResourceName = "snowflake_external_table"
	ExternalOauthSecurityIntegration                       ResourceName = "snowflake_external_oauth_security_integration"
	ExternalVolume                                         ResourceName = "snowflake_external_volume"
	FailoverGroup                                          ResourceName = "snowflake_failover_group"
	FileFormat                                             ResourceName = "snowflake_file_format"
	Function                                               ResourceName = "snowflake_function"
	LegacyServiceUser                                      ResourceName = "snowflake_legacy_service_user"
	ManagedAccount                                         ResourceName = "snowflake_managed_account"
	MaskingPolicy                                          ResourceName = "snowflake_masking_policy"
	MaterializedView                                       ResourceName = "snowflake_materialized_view"
	NetworkPolicy                                          ResourceName = "snowflake_network_policy"
	NetworkRule                                            ResourceName = "snowflake_network_rule"
	NotificationIntegration                                ResourceName = "snowflake_notification_integration"
	OauthIntegrationForCustomClients                       ResourceName = "snowflake_oauth_integration_for_custom_clients"
	OauthIntegrationForPartnerApplications                 ResourceName = "snowflake_oauth_integration_for_partner_applications"
	PasswordPolicy                                         ResourceName = "snowflake_password_policy"
	Pipe                                                   ResourceName = "snowflake_pipe"
	PrimaryConnection                                      ResourceName = "snowflake_primary_connection"
	Procedure                                              ResourceName = "snowflake_procedure"
	ResourceMonitor                                        ResourceName = "snowflake_resource_monitor"
	Role                                                   ResourceName = "snowflake_role"
	RowAccessPolicy                                        ResourceName = "snowflake_row_access_policy"
	Saml2SecurityIntegration                               ResourceName = "snowflake_saml2_integration"
	Schema                                                 ResourceName = "snowflake_schema"
	ScimSecurityIntegration                                ResourceName = "snowflake_scim_integration"
	SecondaryConnection                                    ResourceName = "snowflake_secondary_connection"
	SecondaryDatabase                                      ResourceName = "snowflake_secondary_database"
	SecretWithAuthorizationCodeGrant                       ResourceName = "snowflake_secret_with_authorization_code_grant"
	SecretWithBasicAuthentication                          ResourceName = "snowflake_secret_with_basic_authentication"
	SecretWithClientCredentials                            ResourceName = "snowflake_secret_with_client_credentials"
	SecretWithGenericString                                ResourceName = "snowflake_secret_with_generic_string"
	Sequence                                               ResourceName = "snowflake_sequence"
	ServiceUser                                            ResourceName = "snowflake_service_user"
	Share                                                  ResourceName = "snowflake_share"
	SharedDatabase                                         ResourceName = "snowflake_shared_database"
	Stage                                                  ResourceName = "snowflake_stage"
	StorageIntegration                                     ResourceName = "snowflake_storage_integration"
	Stream                                                 ResourceName = "snowflake_stream"
	StreamOnDirectoryTable                                 ResourceName = "snowflake_stream_on_directory_table"
	StreamOnExternalTable                                  ResourceName = "snowflake_stream_on_external_table"
	StreamOnTable                                          ResourceName = "snowflake_stream_on_table"
	StreamOnView                                           ResourceName = "snowflake_stream_on_view"
	Streamlit                                              ResourceName = "snowflake_streamlit"
	Table                                                  ResourceName = "snowflake_table"
	Tag                                                    ResourceName = "snowflake_tag"
	Task                                                   ResourceName = "snowflake_task"
	User                                                   ResourceName = "snowflake_user"
	View                                                   ResourceName = "snowflake_view"
	Warehouse                                              ResourceName = "snowflake_warehouse"
)

type Resource interface {
	xxxProtected()
	String() string
}

func (r ResourceName) xxxProtected() {}

func (r ResourceName) String() string {
	return string(r)
}
