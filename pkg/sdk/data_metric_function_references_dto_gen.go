package sdk

//go:generate go run ./dto-builder-generator/main.go

var _ optionsProvider[GetForEntityDataMetricFunctionReferenceOptions] = new(GetForEntityDataMetricFunctionReferenceRequest)

type GetForEntityDataMetricFunctionReferenceRequest struct {
	refEntityName   ObjectIdentifier                       // required
	RefEntityDomain DataMetricFuncionRefEntityDomainOption // required
}
