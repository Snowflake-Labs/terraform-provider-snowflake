package resources

type resource string

const (
	Account                      resource = "snowflake_account"
	Alert                        resource = "snowflake_alert"
	ApiIntegration               resource = "snowflake_api_integration"
	Database                     resource = "snowflake_database"
	DatabaseRole                 resource = "snowflake_database_role"
	DynamicTable                 resource = "snowflake_dynamic_table"
	EmailNotificationIntegration resource = "snowflake_email_notification_integration"
	ExternalFunction             resource = "snowflake_external_function"
	ExternalTable                resource = "snowflake_external_table"
	FailoverGroup                resource = "snowflake_failover_group"
	FileFormat                   resource = "snowflake_file_format"
	Function                     resource = "snowflake_function"
	ManagedAccount               resource = "snowflake_managed_account"
	MaskingPolicy                resource = "snowflake_masking_policy"
	MaterializedView             resource = "snowflake_materialized_view"
	NetworkPolicy                resource = "snowflake_network_policy"
	NotificationIntegration      resource = "snowflake_notification_integration"
	PasswordPolicy               resource = "snowflake_password_policy"
	Pipe                         resource = "snowflake_pipe"
	Procedure                    resource = "snowflake_procedure"
	ResourceMonitor              resource = "snowflake_resource_monitor"
	Schema                       resource = "snowflake_schema"
	Stage                        resource = "snowflake_stage"
	View                         resource = "snowflake_view"
)

type Resource interface {
	xxxProtected()
	String() string
}

func (r resource) xxxProtected() {}

func (r resource) String() string {
	return string(r)
}
