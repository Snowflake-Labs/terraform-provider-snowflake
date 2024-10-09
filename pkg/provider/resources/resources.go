package resources

type resource string

const (
	Account                                                resource = "snowflake_account"
	AccountRole                                            resource = "snowflake_account_role"
	Alert                                                  resource = "snowflake_alert"
	ApiAuthenticationIntegrationWithAuthorizationCodeGrant resource = "snowflake_api_authentication_integration_with_authorization_code_grant"
	ApiAuthenticationIntegrationWithClientCredentials      resource = "snowflake_api_authentication_integration_with_client_credentials"
	ApiAuthenticationIntegrationWithJwtBearer              resource = "snowflake_api_authentication_integration_with_jwt_bearer"
	ApiIntegration                                         resource = "snowflake_api_integration"
	CortexSearchService                                    resource = "snowflake_cortex_search_service"
	DatabaseOld                                            resource = "snowflake_database_old"
	Database                                               resource = "snowflake_database"
	DatabaseRole                                           resource = "snowflake_database_role"
	DynamicTable                                           resource = "snowflake_dynamic_table"
	EmailNotificationIntegration                           resource = "snowflake_email_notification_integration"
	ExternalFunction                                       resource = "snowflake_external_function"
	ExternalTable                                          resource = "snowflake_external_table"
	ExternalOauthSecurityIntegration                       resource = "snowflake_external_oauth_security_integration"
	FailoverGroup                                          resource = "snowflake_failover_group"
	FileFormat                                             resource = "snowflake_file_format"
	Function                                               resource = "snowflake_function"
	LegacyServiceUser                                      resource = "snowflake_legacy_service_user"
	ManagedAccount                                         resource = "snowflake_managed_account"
	MaskingPolicy                                          resource = "snowflake_masking_policy"
	MaterializedView                                       resource = "snowflake_materialized_view"
	NetworkPolicy                                          resource = "snowflake_network_policy"
	NetworkRule                                            resource = "snowflake_network_rule"
	NotificationIntegration                                resource = "snowflake_notification_integration"
	OauthIntegrationForCustomClients                       resource = "snowflake_oauth_integration_for_custom_clients"
	OauthIntegrationForPartnerApplications                 resource = "snowflake_oauth_integration_for_partner_applications"
	PasswordPolicy                                         resource = "snowflake_password_policy"
	Pipe                                                   resource = "snowflake_pipe"
	Procedure                                              resource = "snowflake_procedure"
	ResourceMonitor                                        resource = "snowflake_resource_monitor"
	Role                                                   resource = "snowflake_role"
	RowAccessPolicy                                        resource = "snowflake_row_access_policy"
	Saml2SecurityIntegration                               resource = "snowflake_saml2_integration"
	Schema                                                 resource = "snowflake_schema"
	ScimSecurityIntegration                                resource = "snowflake_scim_integration"
	SecondaryDatabase                                      resource = "snowflake_secondary_database"
	Sequence                                               resource = "snowflake_sequence"
	ServiceUser                                            resource = "snowflake_service_user"
	Share                                                  resource = "snowflake_share"
	SharedDatabase                                         resource = "snowflake_shared_database"
	Stage                                                  resource = "snowflake_stage"
	StorageIntegration                                     resource = "snowflake_storage_integration"
	Stream                                                 resource = "snowflake_stream"
	StreamOnTable                                          resource = "snowflake_stream_on_table"
	Streamlit                                              resource = "snowflake_streamlit"
	Table                                                  resource = "snowflake_table"
	Tag                                                    resource = "snowflake_tag"
	Task                                                   resource = "snowflake_task"
	User                                                   resource = "snowflake_user"
	View                                                   resource = "snowflake_view"
	Warehouse                                              resource = "snowflake_warehouse"
)

type Resource interface {
	xxxProtected()
	String() string
}

func (r resource) xxxProtected() {}

func (r resource) String() string {
	return string(r)
}
