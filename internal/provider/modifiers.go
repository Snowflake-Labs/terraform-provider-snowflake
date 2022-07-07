package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type caseInsensitive struct {
	emptyDescriptions
}

func (caseInsensitive) Modify(_ context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeState == nil {
		return
	}
	state := req.AttributeState.(types.String)
	plan := req.AttributePlan.(types.String)
	if strings.EqualFold(state.Value, plan.Value) {
		resp.AttributePlan = state
	} else {
		resp.AttributePlan = plan
	}
}
