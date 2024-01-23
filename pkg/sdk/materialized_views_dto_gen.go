package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateMaterializedViewOptions]   = new(CreateMaterializedViewRequest)
	_ optionsProvider[AlterMaterializedViewOptions]    = new(AlterMaterializedViewRequest)
	_ optionsProvider[DropMaterializedViewOptions]     = new(DropMaterializedViewRequest)
	_ optionsProvider[ShowMaterializedViewOptions]     = new(ShowMaterializedViewRequest)
	_ optionsProvider[DescribeMaterializedViewOptions] = new(DescribeMaterializedViewRequest)
)

type CreateMaterializedViewRequest struct {
	OrReplace              *bool
	Secure                 *bool
	IfNotExists            *bool
	name                   SchemaObjectIdentifier // required
	CopyGrants             *bool
	Columns                []MaterializedViewColumnRequest
	ColumnsMaskingPolicies []MaterializedViewColumnMaskingPolicyRequest
	Comment                *string
	RowAccessPolicy        *MaterializedViewRowAccessPolicyRequest
	Tag                    []TagAssociation
	ClusterBy              *MaterializedViewClusterByRequest
	sql                    string // required
}

func (r *CreateMaterializedViewRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type MaterializedViewColumnRequest struct {
	Name    string // required
	Comment *string
}

type MaterializedViewColumnMaskingPolicyRequest struct {
	Name          string                 // required
	MaskingPolicy SchemaObjectIdentifier // required
	Using         []string
	Tag           []TagAssociation
}

type MaterializedViewRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
	On              []string               // required
}

type MaterializedViewClusterByRequest struct {
	Expressions []MaterializedViewClusterByExpressionRequest
}

type MaterializedViewClusterByExpressionRequest struct {
	Name string // required
}

type AlterMaterializedViewRequest struct {
	name              SchemaObjectIdentifier // required
	RenameTo          *SchemaObjectIdentifier
	ClusterBy         *MaterializedViewClusterByRequest
	DropClusteringKey *bool
	SuspendRecluster  *bool
	ResumeRecluster   *bool
	Suspend           *bool
	Resume            *bool
	Set               *MaterializedViewSetRequest
	Unset             *MaterializedViewUnsetRequest
}

type MaterializedViewSetRequest struct {
	Secure  *bool
	Comment *string
}

type MaterializedViewUnsetRequest struct {
	Secure  *bool
	Comment *bool
}

type DropMaterializedViewRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowMaterializedViewRequest struct {
	Like *Like
	In   *In
}

type DescribeMaterializedViewRequest struct {
	name SchemaObjectIdentifier // required
}
