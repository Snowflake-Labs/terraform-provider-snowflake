package snowflakeroles

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

var (
	Orgadmin      = sdk.NewAccountObjectIdentifier("ORGADMIN")
	Accountadmin  = sdk.NewAccountObjectIdentifier("ACCOUNTADMIN")
	SecurityAdmin = sdk.NewAccountObjectIdentifier("SECURITYADMIN")

	OktaProvisioner        = sdk.NewAccountObjectIdentifier("OKTA_PROVISIONER")
	AadProvisioner         = sdk.NewAccountObjectIdentifier("AAD_PROVISIONER")
	GenericScimProvisioner = sdk.NewAccountObjectIdentifier("GENERIC_SCIM_PROVISIONER")
)
