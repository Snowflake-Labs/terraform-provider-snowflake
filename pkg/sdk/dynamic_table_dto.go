package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[createDynamicTableOptions] = new(CreateDynamicTableRequest)
	_ optionsProvider[alterDynamicTableOptions]  = new(AlterDynamicTableRequest)
	_ optionsProvider[dropDynamicTableOptions]   = new(DropDynamicTableRequest)
)

type CreateDynamicTableRequest struct {
	orReplace bool

	name      AccountObjectIdentifier // required
	warehouse AccountObjectIdentifier // required
	targetLag TargetLag               // required
	query     string                  // required

	comment *string
}

type AlterDynamicTableRequest struct {
	name AccountObjectIdentifier // required

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
	name AccountObjectIdentifier // required
}

type DescribeDynamicTableRequest struct {
	name AccountObjectIdentifier // required
}
