package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const showOutputAttributeName = "show_output"

// handleExternalChangesToObject assumes that show output is kept in showOutputAttributeName attribute
// TODO [after discussion/next PR]: fix/make safer (casting)
// TODO [after discussion/next PR]: replace func with generic struct to build this internally?
func handleExternalChangesToObject(d *schema.ResourceData, handler func(map[string]any) error) error {
	if showOutput, ok := d.GetOk(showOutputAttributeName); ok {
		showOutputList := showOutput.([]any)
		if len(showOutputList) == 1 {
			result := showOutputList[0].(map[string]any)
			return handler(result)
		}
	}
	return nil
}
