package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateViewOptions]   = new(CreateViewRequest)
	_ optionsProvider[AlterViewOptions]    = new(AlterViewRequest)
	_ optionsProvider[DropViewOptions]     = new(DropViewRequest)
	_ optionsProvider[ShowViewOptions]     = new(ShowViewRequest)
	_ optionsProvider[DescribeViewOptions] = new(DescribeViewRequest)
)

type CreateViewRequest struct {
	OrReplace              *bool
	Secure                 *bool
	Temporary              *bool
	Recursive              *bool
	IfNotExists            *bool
	name                   SchemaObjectIdentifier // required
	Columns                []ViewColumnRequest
	ColumnsMaskingPolicies []ViewColumnMaskingPolicyRequest
	CopyGrants             *bool
	Comment                *string
	RowAccessPolicy        *ViewRowAccessPolicyRequest
	Tag                    []TagAssociation
	sql                    string // required
}

func (r *CreateViewRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type ViewColumnRequest struct {
	Name    string // required
	Comment *string
}

type ViewColumnMaskingPolicyRequest struct {
	Name          string                 // required
	MaskingPolicy SchemaObjectIdentifier // required
	Using         []string
	Tag           []TagAssociation
}

type ViewRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
	On              []string               // required
}

type AlterViewRequest struct {
	IfExists                   *bool
	name                       SchemaObjectIdentifier // required
	RenameTo                   *SchemaObjectIdentifier
	SetComment                 *string
	UnsetComment               *bool
	SetSecure                  *bool
	SetChangeTracking          *bool
	UnsetSecure                *bool
	SetTags                    []TagAssociation
	UnsetTags                  []ObjectIdentifier
	AddRowAccessPolicy         *ViewAddRowAccessPolicyRequest
	DropRowAccessPolicy        *ViewDropRowAccessPolicyRequest
	DropAndAddRowAccessPolicy  *ViewDropAndAddRowAccessPolicyRequest
	DropAllRowAccessPolicies   *bool
	SetMaskingPolicyOnColumn   *ViewSetColumnMaskingPolicyRequest
	UnsetMaskingPolicyOnColumn *ViewUnsetColumnMaskingPolicyRequest
	SetTagsOnColumn            *ViewSetColumnTagsRequest
	UnsetTagsOnColumn          *ViewUnsetColumnTagsRequest
}

type ViewAddRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
	On              []string               // required
}

type ViewDropRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
}

type ViewDropAndAddRowAccessPolicyRequest struct {
	Drop ViewDropRowAccessPolicyRequest // required
	Add  ViewAddRowAccessPolicyRequest  // required
}

type ViewSetColumnMaskingPolicyRequest struct {
	Name          string                 // required
	MaskingPolicy SchemaObjectIdentifier // required
	Using         []string
	Force         *bool
}

type ViewUnsetColumnMaskingPolicyRequest struct {
	Name string // required
}

type ViewSetColumnTagsRequest struct {
	Name    string // required
	SetTags []TagAssociation
}

type ViewUnsetColumnTagsRequest struct {
	Name      string // required
	UnsetTags []ObjectIdentifier
}

type DropViewRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowViewRequest struct {
	Terse      *bool
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeViewRequest struct {
	name SchemaObjectIdentifier // required
}
