package bettertestspoc

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

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

	b, err := json.MarshalIndent(model, "", "  ")
	require.NoError(t, err)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`resource "%s" "%s" {`, model.Resource(), model.ResourceName()))
	sb.WriteString(strings.Trim(string(b), "{}"))
	sb.WriteString(`}`)
	sb.WriteRune('\n')
	s := sb.String()
	t.Logf("Generated config:\n%s", s)
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
