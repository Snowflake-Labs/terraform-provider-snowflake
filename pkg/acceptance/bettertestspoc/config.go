package bettertestspoc

import (
	"encoding/json"
	"testing"

	hclJson "github.com/hashicorp/hcl2/hcl/json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/stretchr/testify/require"
)

type ResourceModel interface {
	Resource() resources.Resource
	ResourceName() string
	SetResourceName(name string)
}

type resourceModelMeta struct {
	name     string
	resource resources.Resource
}

func (m *resourceModelMeta) Resource() resources.Resource {
	return m.resource
}

func (m *resourceModelMeta) ResourceName() string {
	return m.name
}

func (m *resourceModelMeta) SetResourceName(name string) {
	m.name = name
}

const DefaultResourceName = "test"

func defaultMeta(resource resources.Resource) *resourceModelMeta {
	return &resourceModelMeta{name: DefaultResourceName, resource: resource}
}

func ConfigurationFromModel(t *testing.T, model ResourceModel) string {
	t.Helper()

	m1 := make(map[string]map[string]ResourceModel)
	m2 := make(map[string]ResourceModel)
	m2[model.ResourceName()] = model
	m1[model.Resource().String()] = m2
	b, err := json.Marshal(ResourceWrapper{m1})
	require.NoError(t, err)
	t.Logf("Generated json:\n%s", string(b))
	// TODO: https://pkg.go.dev/github.com/hashicorp/hcl2/hcl/json#Parse
	f, diag := hclJson.Parse(b, "")
	if diag.HasErrors() {
		t.Fatal("Could not parse")
	}
	s := string(f.Bytes)
	t.Logf("Generated hcl:\n%s", s)
	return s
}

// TODO: save to tmp file and return path to it
func ConfigurationFromModelProvider(t *testing.T, model ResourceModel) func(config.TestStepConfigRequest) string {
	t.Helper()
	return func(req config.TestStepConfigRequest) string {
		t.Logf("Generating config for test %s, step %d for resource %s", req.TestName, req.StepNumber, model.Resource())
		content := ConfigurationFromModel(t, model)
		_ = content
		return ""
	}
}

type ResourceWrapper struct {
	Resource map[string]map[string]ResourceModel `json:"resource"`
}
