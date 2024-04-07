package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[createDynamicTableOptions] = new(CreateDynamicTableRequest)
	_ optionsProvider[alterDynamicTableOptions]  = new(AlterDynamicTableRequest)
	_ optionsProvider[dropDynamicTableOptions]   = new(DropDynamicTableRequest)
	_ optionsProvider[showDynamicTableOptions]   = new(ShowDynamicTableRequest)
)

type CreateDynamicTableRequest struct {
	name      SchemaObjectIdentifier  // required
	warehouse AccountObjectIdentifier // required
	targetLag TargetLag               // required
	query     string                  // required

	comment     *string
	refreshMode *DynamicTableRefreshMode
	initialize  *DynamicTableInitialize
}

type AlterDynamicTableRequest struct {
	name SchemaObjectIdentifier // required

	// One of
	suspend *bool
	resume  *bool
	refresh *bool
	set     *DynamicTableSetRequest
}

type DynamicTableSetRequest struct {
	targetLag  *TargetLag
	warehourse *AccountObjectIdentifier
}

type DropDynamicTableRequest struct {
	name SchemaObjectIdentifier // required
}

type DescribeDynamicTableRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowDynamicTableRequest struct {
	like       *Like
	in         *In
	startsWith *string
	limit      *LimitFrom
}
