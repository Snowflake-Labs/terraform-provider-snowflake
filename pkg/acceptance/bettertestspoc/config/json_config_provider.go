package config

import (
	"encoding/json"
	"fmt"
)

var DefaultJsonProvider = NewBasicJsonProvider()

// TODO [SNOW-1501905]: add config builders for other block types (Variable, Output, Localsl, Module, Terraform)
type JsonProvider interface {
	ResourceJsonFromModel(model ResourceModel) ([]byte, error)
	DatasourceJsonFromModel(model DatasourceModel) ([]byte, error)
	ProviderJsonFromModel(model ProviderModel) ([]byte, error)
}

type basicJsonProvider struct{}

func NewBasicJsonProvider() JsonProvider {
	return &basicJsonProvider{}
}

func (p *basicJsonProvider) ResourceJsonFromModel(model ResourceModel) ([]byte, error) {
	modelJson := resourceJson{
		Resource: map[string]map[string]ResourceModel{
			fmt.Sprintf("%s", model.Resource()): {
				fmt.Sprintf("%s", model.ResourceName()): model,
			},
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}

type resourceJson struct {
	Resource map[string]map[string]ResourceModel `json:"resource"`
}

func (p *basicJsonProvider) DatasourceJsonFromModel(model DatasourceModel) ([]byte, error) {
	modelJson := datasourceJson{
		Datasource: map[string]map[string]DatasourceModel{
			fmt.Sprintf("%s", model.Datasource()): {
				fmt.Sprintf("%s", model.DatasourceName()): model,
			},
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}

type datasourceJson struct {
	Datasource map[string]map[string]DatasourceModel `json:"data"`
}

func (p *basicJsonProvider) ProviderJsonFromModel(model ProviderModel) ([]byte, error) {
	modelJson := providerJson{
		Provider: map[string]ProviderModel{
			fmt.Sprintf("%s", model.ProviderName()): model,
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}

type providerJson struct {
	Provider map[string]ProviderModel `json:"provider"`
}
