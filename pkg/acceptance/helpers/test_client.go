package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type TestClient struct {
	context *TestClientContext

	Ids *IdsGenerator

	Account                   *AccountClient
	AggregationPolicy         *AggregationPolicyClient
	Alert                     *AlertClient
	ApiIntegration            *ApiIntegrationClient
	Application               *ApplicationClient
	ApplicationPackage        *ApplicationPackageClient
	Context                   *ContextClient
	CortexSearchService       *CortexSearchServiceClient
	CatalogIntegration        *CatalogIntegrationClient
	Database                  *DatabaseClient
	DatabaseRole              *DatabaseRoleClient
	DynamicTable              *DynamicTableClient
	ExternalAccessIntegration *ExternalAccessIntegrationClient
	ExternalVolume            *ExternalVolumeClient
	FailoverGroup             *FailoverGroupClient
	FileFormat                *FileFormatClient
	MaskingPolicy             *MaskingPolicyClient
	MaterializedView          *MaterializedViewClient
	NetworkPolicy             *NetworkPolicyClient
	NetworkRule               *NetworkRuleClient
	Parameter                 *ParameterClient
	PasswordPolicy            *PasswordPolicyClient
	Pipe                      *PipeClient
	ProjectionPolicy          *ProjectionPolicyClient
	ResourceMonitor           *ResourceMonitorClient
	Role                      *RoleClient
	RowAccessPolicy           *RowAccessPolicyClient
	Schema                    *SchemaClient
	SecurityIntegration       *SecurityIntegrationClient
	SessionPolicy             *SessionPolicyClient
	Share                     *ShareClient
	Stage                     *StageClient
	Streamlit                 *StreamlitClient
	Table                     *TableClient
	Tag                       *TagClient
	Task                      *TaskClient
	User                      *UserClient
	View                      *ViewClient
	Warehouse                 *WarehouseClient
}

func NewTestClient(c *sdk.Client, database string, schema string, warehouse string, testObjectSuffix string) *TestClient {
	context := &TestClientContext{
		client:           c,
		database:         database,
		schema:           schema,
		warehouse:        warehouse,
		testObjectSuffix: testObjectSuffix,
	}
	idsGenerator := NewIdsGenerator(context)
	return &TestClient{
		context: context,

		Ids: idsGenerator,

		Account:                   NewAccountClient(context),
		AggregationPolicy:         NewAggregationPolicyClient(context, idsGenerator),
		Alert:                     NewAlertClient(context, idsGenerator),
		ApiIntegration:            NewApiIntegrationClient(context, idsGenerator),
		Application:               NewApplicationClient(context, idsGenerator),
		ApplicationPackage:        NewApplicationPackageClient(context, idsGenerator),
		Context:                   NewContextClient(context),
		CortexSearchService:       NewCortexSearchServiceClient(context, idsGenerator),
		CatalogIntegration:        NewCatalogIntegrationClient(context, idsGenerator),
		Database:                  NewDatabaseClient(context, idsGenerator),
		DatabaseRole:              NewDatabaseRoleClient(context, idsGenerator),
		DynamicTable:              NewDynamicTableClient(context, idsGenerator),
		ExternalAccessIntegration: NewExternalAccessIntegrationClient(context, idsGenerator),
		ExternalVolume:            NewExternalVolumeClient(context, idsGenerator),
		FailoverGroup:             NewFailoverGroupClient(context, idsGenerator),
		FileFormat:                NewFileFormatClient(context, idsGenerator),
		MaskingPolicy:             NewMaskingPolicyClient(context, idsGenerator),
		MaterializedView:          NewMaterializedViewClient(context, idsGenerator),
		NetworkPolicy:             NewNetworkPolicyClient(context, idsGenerator),
		NetworkRule:               NewNetworkRuleClient(context, idsGenerator),
		Parameter:                 NewParameterClient(context),
		PasswordPolicy:            NewPasswordPolicyClient(context, idsGenerator),
		Pipe:                      NewPipeClient(context, idsGenerator),
		ProjectionPolicy:          NewProjectionPolicyClient(context, idsGenerator),
		ResourceMonitor:           NewResourceMonitorClient(context, idsGenerator),
		Role:                      NewRoleClient(context, idsGenerator),
		RowAccessPolicy:           NewRowAccessPolicyClient(context, idsGenerator),
		Schema:                    NewSchemaClient(context, idsGenerator),
		SecurityIntegration:       NewSecurityIntegrationClient(context, idsGenerator),
		SessionPolicy:             NewSessionPolicyClient(context, idsGenerator),
		Share:                     NewShareClient(context, idsGenerator),
		Streamlit:                 NewStreamlitClient(context, idsGenerator),
		Stage:                     NewStageClient(context, idsGenerator),
		Table:                     NewTableClient(context, idsGenerator),
		Tag:                       NewTagClient(context, idsGenerator),
		Task:                      NewTaskClient(context, idsGenerator),
		User:                      NewUserClient(context, idsGenerator),
		View:                      NewViewClient(context, idsGenerator),
		Warehouse:                 NewWarehouseClient(context, idsGenerator),
	}
}

type TestClientContext struct {
	client           *sdk.Client
	database         string
	schema           string
	warehouse        string
	testObjectSuffix string
}
