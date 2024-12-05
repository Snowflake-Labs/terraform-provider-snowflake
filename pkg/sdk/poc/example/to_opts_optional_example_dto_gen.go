package example

//go:generate go run ./dto-builder-generator/main.go

var _ optionsProvider[AlterFeaturesExamplesOptions] = new(AlterFeaturesExamplesRequest)

type AlterFeaturesExamplesRequest struct {
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
