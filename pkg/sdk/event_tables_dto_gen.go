package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateEventTableOptions]   = new(CreateEventTableRequest)
	_ optionsProvider[ShowEventTableOptions]     = new(ShowEventTableRequest)
	_ optionsProvider[DescribeEventTableOptions] = new(DescribeEventTableRequest)
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
	RowAccessPolicy            *RowAccessPolicy
	Tag                        []TagAssociation
}

type ShowEventTableRequest struct {
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *int
	From       *string
}

type DescribeEventTableRequest struct {
	name SchemaObjectIdentifier // required
}

type AlterEventTableRequest struct {
	IfNotExists              *bool
	name                     SchemaObjectIdentifier // required
	Set                      *EventTableSetRequest
	Unset                    *EventTableUnsetRequest
	AddRowAccessPolicy       *RowAccessPolicy
	DropRowAccessPolicy      *EventTableDropRowAccessPolicyRequest
	DropAllRowAccessPolicies *bool
	ClusteringAction         *EventTableClusteringActionRequest
	SearchOptimizationAction *EventTableSearchOptimizationActionRequest
	SetTags                  []TagAssociation
	UnsetTags                []ObjectIdentifier
	RenameTo                 *SchemaObjectIdentifier
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

type EventTableDropRowAccessPolicyRequest struct {
	Name SchemaObjectIdentifier
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
