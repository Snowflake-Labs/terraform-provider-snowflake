package config

import (
	"encoding/json"
	"fmt"
)

var DefaultJsonProvider = NewBasicJsonProvider()

type JsonProvider interface {
	ResourceJsonFromModel(model ResourceModel) ([]byte, error)
	// Variable
	// Output
	// Locals
	// Module
	// Provider
	// Terraform
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
