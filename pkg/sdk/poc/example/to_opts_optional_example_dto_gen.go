package example

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateToOptsOptionalExampleOptions] = new(CreateToOptsOptionalExampleRequest)
	_ optionsProvider[AlterToOptsOptionalExampleOptions]  = new(AlterToOptsOptionalExampleRequest)
)

type CreateToOptsOptionalExampleRequest struct {
	IfExists *bool
	name     DatabaseObjectIdentifier // required
}

type AlterToOptsOptionalExampleRequest struct {
	IfExists      *bool
	name          DatabaseObjectIdentifier // required
	OptionalField *OptionalFieldRequest
	RequiredField RequiredFieldRequest // required
}

type OptionalFieldRequest struct {
	SomeList []DatabaseObjectIdentifier
}

type RequiredFieldRequest struct {
	SomeRequiredList []DatabaseObjectIdentifier // required
}
