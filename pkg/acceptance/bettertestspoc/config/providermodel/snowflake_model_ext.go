package providermodel

import "encoding/json"

// Based on https://medium.com/picus-security-engineering/custom-json-marshaller-in-go-and-common-pitfalls-c43fa774db05.
func (m *SnowflakeModel) MarshalJSON() ([]byte, error) {
	type AliasModelType SnowflakeModel
	return json.Marshal(&struct {
		*AliasModelType
		Alias string `json:"alias,omitempty"`
	}{
		AliasModelType: (*AliasModelType)(m),
		Alias:          m.Alias(),
	})
}
