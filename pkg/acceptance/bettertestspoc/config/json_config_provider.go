package config

import (
	"encoding/json"
	"fmt"
)

type JsonProvider interface {
	JsonFromModel(model ResourceModel) ([]byte, error)
}

type basicJsonProvider struct{}

func NewBasicJsonProvider() JsonProvider {
	return &basicJsonProvider{}
}

func (p *basicJsonProvider) JsonFromModel(model ResourceModel) ([]byte, error) {
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
