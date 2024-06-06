package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[createCortexSearchServiceOptions] = new(CreateCortexSearchServiceRequest)
	_ optionsProvider[alterCortexSearchServiceOptions]  = new(AlterCortexSearchServiceRequest)
	_ optionsProvider[dropCortexSearchServiceOptions]   = new(DropCortexSearchServiceRequest)
	_ optionsProvider[showCortexSearchServiceOptions]   = new(ShowCortexSearchServiceRequest)
)

type CreateCortexSearchServiceRequest struct {
	orReplace   bool
	ifNotExists bool

	name      SchemaObjectIdentifier  // required
	on        string                  // required
	warehouse AccountObjectIdentifier // required
	targetLag string                  // required
	query     string                  // required

	attributes []string
	comment    *string
}

type AlterCortexSearchServiceRequest struct {
	name     SchemaObjectIdentifier // required
	IfExists *bool

	// One of
	set *CortexSearchServiceSetRequest
}

type CortexSearchServiceSetRequest struct {
	targetLag *string
	warehouse *AccountObjectIdentifier
}

type DropCortexSearchServiceRequest struct {
	name     SchemaObjectIdentifier // required
	IfExists *bool
}

type DescribeCortexSearchServiceRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowCortexSearchServiceRequest struct {
	like       *Like
	in         *In
	startsWith *string
	limit      *LimitFrom
}
