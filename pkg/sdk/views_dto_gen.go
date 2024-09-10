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
	OrReplace         *bool
	Secure            *bool
	Temporary         *bool
	Recursive         *bool
	IfNotExists       *bool
	name              SchemaObjectIdentifier // required
	Columns           []ViewColumnRequest
	CopyGrants        *bool
	Comment           *string
	RowAccessPolicy   *ViewRowAccessPolicyRequest
	AggregationPolicy *ViewAggregationPolicyRequest
	Tag               []TagAssociation
	sql               string // required
}

func (r *CreateViewRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type ViewColumnRequest struct {
	Name             string // required
	ProjectionPolicy *ViewColumnProjectionPolicyRequest
	MaskingPolicy    *ViewColumnMaskingPolicyRequest
	Comment          *string
	Tag              []TagAssociation
}

type ViewColumnProjectionPolicyRequest struct {
	ProjectionPolicy SchemaObjectIdentifier // required
}

type ViewColumnMaskingPolicyRequest struct {
	MaskingPolicy SchemaObjectIdentifier // required
	Using         []Column
}

type ViewRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
	On              []Column               // required
}

type ViewAggregationPolicyRequest struct {
	AggregationPolicy SchemaObjectIdentifier // required
	EntityKey         []Column
}

type AlterViewRequest struct {
	IfExists                      *bool
	name                          SchemaObjectIdentifier // required
	RenameTo                      *SchemaObjectIdentifier
	SetComment                    *string
	UnsetComment                  *bool
	SetSecure                     *bool
	SetChangeTracking             *bool
	UnsetSecure                   *bool
	SetTags                       []TagAssociation
	UnsetTags                     []ObjectIdentifier
	AddDataMetricFunction         *ViewAddDataMetricFunctionRequest
	DropDataMetricFunction        *ViewDropDataMetricFunctionRequest
	ModifyDataMetricFunction      *ViewModifyDataMetricFunctionsRequest
	SetDataMetricSchedule         *ViewSetDataMetricScheduleRequest
	UnsetDataMetricSchedule       *ViewUnsetDataMetricScheduleRequest
	AddRowAccessPolicy            *ViewAddRowAccessPolicyRequest
	DropRowAccessPolicy           *ViewDropRowAccessPolicyRequest
	DropAndAddRowAccessPolicy     *ViewDropAndAddRowAccessPolicyRequest
	DropAllRowAccessPolicies      *bool
	SetAggregationPolicy          *ViewSetAggregationPolicyRequest
	UnsetAggregationPolicy        *ViewUnsetAggregationPolicyRequest
	SetMaskingPolicyOnColumn      *ViewSetColumnMaskingPolicyRequest
	UnsetMaskingPolicyOnColumn    *ViewUnsetColumnMaskingPolicyRequest
	SetProjectionPolicyOnColumn   *ViewSetProjectionPolicyRequest
	UnsetProjectionPolicyOnColumn *ViewUnsetProjectionPolicyRequest
	SetTagsOnColumn               *ViewSetColumnTagsRequest
	UnsetTagsOnColumn             *ViewUnsetColumnTagsRequest
}

type ViewAddDataMetricFunctionRequest struct {
	DataMetricFunction []ViewDataMetricFunction // required
}

type ViewDropDataMetricFunctionRequest struct {
	DataMetricFunction []ViewDataMetricFunction // required
}

type ViewModifyDataMetricFunctionsRequest struct {
	DataMetricFunction []ViewModifyDataMetricFunction // required
}

type ViewSetDataMetricScheduleRequest struct {
	DataMetricSchedule string // required
}

type ViewUnsetDataMetricScheduleRequest struct{}

type ViewAddRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
	On              []Column               // required
}

type ViewDropRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
}

type ViewDropAndAddRowAccessPolicyRequest struct {
	Drop ViewDropRowAccessPolicyRequest // required
	Add  ViewAddRowAccessPolicyRequest  // required
}

type ViewSetAggregationPolicyRequest struct {
	AggregationPolicy SchemaObjectIdentifier // required
	EntityKey         []Column
	Force             *bool
}

type ViewUnsetAggregationPolicyRequest struct{}

type ViewSetColumnMaskingPolicyRequest struct {
	Name          string                 // required
	MaskingPolicy SchemaObjectIdentifier // required
	Using         []Column
	Force         *bool
}

type ViewUnsetColumnMaskingPolicyRequest struct {
	Name string // required
}

type ViewSetProjectionPolicyRequest struct {
	Name             string                 // required
	ProjectionPolicy SchemaObjectIdentifier // required
	Force            *bool
}

type ViewUnsetProjectionPolicyRequest struct {
	Name string // required
}

type ViewSetColumnTagsRequest struct {
	Name    string           // required
	SetTags []TagAssociation // required
}

type ViewUnsetColumnTagsRequest struct {
	Name      string             // required
	UnsetTags []ObjectIdentifier // required
}

type DropViewRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowViewRequest struct {
	Terse      *bool
	Like       *Like
	In         *ExtendedIn
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeViewRequest struct {
	name SchemaObjectIdentifier // required
}
