package resources

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var serviceTypeUserNotApplicableAttributes = []string{
	"password",
	"first_name",
	"middle_name",
	"last_name",
	"must_change_password",
	"mins_to_bypass_mfa",
	"disable_mfa",
}

var legacyServiceTypeUserNotApplicableAttributes = []string{
	"first_name",
	"middle_name",
	"last_name",
	"mins_to_bypass_mfa",
	"disable_mfa",
}

var (
	serviceUserSchema       = make(map[string]*schema.Schema)
	legacyServiceUserSchema = make(map[string]*schema.Schema)
)

func init() {
	for k, v := range userSchema {
		if !slices.Contains(serviceTypeUserNotApplicableAttributes, k) {
			serviceUserSchema[k] = v
		}
		if !slices.Contains(legacyServiceTypeUserNotApplicableAttributes, k) {
			legacyServiceUserSchema[k] = v
		}
	}
}
