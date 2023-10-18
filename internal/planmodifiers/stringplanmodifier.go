package stringplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// useStateForUnknownModifier implements the plan modifier.
type suppressDiffIfModifier struct {
	f func(old, new string) bool
}

// Description returns a human-readable description of the plan modifier.
func (m suppressDiffIfModifier) Description(_ context.Context) string {
	return "Suppresses diff if values based on function."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m suppressDiffIfModifier) MarkdownDescription(_ context.Context) string {
	return "Suppresses diff if values based on function."
}

// PlanModifyBool implements the plan modification logic.
func (m suppressDiffIfModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if m.f(req.StateValue.ValueString(), req.PlanValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

func SuppressDiffIf(f func(old, new string) bool) planmodifier.String {
	return suppressDiffIfModifier{
		f: f,
	}
}
