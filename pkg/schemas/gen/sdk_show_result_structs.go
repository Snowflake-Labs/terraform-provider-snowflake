package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

var SdkShowResultStructs = []any{
	sdk.Account{},
	sdk.Alert{},
	sdk.ApiIntegration{},
	sdk.ApplicationPackage{},
	sdk.ApplicationRole{},
	sdk.Application{},
	sdk.DatabaseRole{},
	sdk.Database{},
	sdk.DynamicTable{},
	sdk.EventTable{},
	sdk.ExternalFunction{},
	sdk.ExternalTable{},
	sdk.FailoverGroup{},
	sdk.FileFormat{},
	sdk.Function{},
	sdk.Grant{},
	sdk.ManagedAccount{},
	sdk.MaskingPolicy{},
	sdk.MaterializedView{},
	sdk.NetworkPolicy{},
	sdk.NetworkRule{},
	sdk.NotificationIntegration{},
	sdk.Parameter{},
	sdk.PasswordPolicy{},
	sdk.Pipe{},
	sdk.PolicyReference{},
	sdk.Procedure{},
	sdk.ReplicationAccount{},
	sdk.ReplicationDatabase{},
	sdk.Region{},
	sdk.ResourceMonitor{},
	sdk.Role{},
	sdk.RowAccessPolicy{},
	sdk.Schema{},
	sdk.SecurityIntegration{},
	sdk.Sequence{},
	sdk.SessionPolicy{},
	sdk.Share{},
	sdk.Stage{},
	sdk.StorageIntegration{},
	sdk.Streamlit{},
	sdk.Stream{},
	sdk.Table{},
	sdk.Tag{},
	sdk.Task{},
	sdk.User{},
	sdk.View{},
	sdk.Warehouse{},
}

// TODO [SNOW-1501905]: currently all this structs have the "Show" added to the schema, while these are not show outputs
// TODO [SNOW-1501905]: temporary struct, may be refactored with addition to generation of describe results; for now used to some structs needing a schema representation
var AdditionalStructs = []any{
	sdk.SecurityIntegrationProperty{},
}
