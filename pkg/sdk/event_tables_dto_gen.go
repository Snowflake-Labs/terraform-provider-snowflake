package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateEventTableOptions]   = new(CreateEventTableRequest)
	_ optionsProvider[ShowEventTableOptions]     = new(ShowEventTableRequest)
	_ optionsProvider[DescribeEventTableOptions] = new(DescribeEventTableRequest)
	_ optionsProvider[DropEventTableOptions]     = new(DropEventTableRequest)
	_ optionsProvider[AlterEventTableOptions]    = new(AlterEventTableRequest)
)

type CreateEventTableRequest struct {
	OrReplace                  *bool
	IfNotExists                *bool
	name                       SchemaObjectIdentifier // required
	ClusterBy                  []string
	DataRetentionTimeInDays    *int
	MaxDataExtensionTimeInDays *int
	ChangeTracking             *bool
	DefaultDdlCollation        *string
	CopyGrants                 *bool
	Comment                    *string
	RowAccessPolicy            *TableRowAccessPolicy
	Tag                        []TagAssociation
}

type ShowEventTableRequest struct {
	Terse      *bool
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeEventTableRequest struct {
	name SchemaObjectIdentifier // required
}

type DropEventTableRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Restrict *bool
}

type AlterEventTableRequest struct {
	IfNotExists               *bool
	name                      SchemaObjectIdentifier // required
	Set                       *EventTableSetRequest
	Unset                     *EventTableUnsetRequest
	AddRowAccessPolicy        *EventTableAddRowAccessPolicyRequest
	DropRowAccessPolicy       *EventTableDropRowAccessPolicyRequest
	DropAndAddRowAccessPolicy *EventTableDropAndAddRowAccessPolicyRequest
	DropAllRowAccessPolicies  *bool
	ClusteringAction          *EventTableClusteringActionRequest
	SearchOptimizationAction  *EventTableSearchOptimizationActionRequest
	SetTags                   []TagAssociation
	UnsetTags                 []ObjectIdentifier
	RenameTo                  *SchemaObjectIdentifier
}

type EventTableSetRequest struct {
	DataRetentionTimeInDays    *int
	MaxDataExtensionTimeInDays *int
	ChangeTracking             *bool
	Comment                    *string
}

type EventTableUnsetRequest struct {
	DataRetentionTimeInDays    *bool
	MaxDataExtensionTimeInDays *bool
	ChangeTracking             *bool
	Comment                    *bool
}

type EventTableAddRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
	On              []string               // required
}

type EventTableDropRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
}

type EventTableDropAndAddRowAccessPolicyRequest struct {
	Drop EventTableDropRowAccessPolicyRequest // required
	Add  EventTableAddRowAccessPolicyRequest  // required
}

type EventTableClusteringActionRequest struct {
	ClusterBy         *[]string
	SuspendRecluster  *bool
	ResumeRecluster   *bool
	DropClusteringKey *bool
}

type EventTableSearchOptimizationActionRequest struct {
	Add  *SearchOptimizationRequest
	Drop *SearchOptimizationRequest
}

type SearchOptimizationRequest struct {
	On []string
}
