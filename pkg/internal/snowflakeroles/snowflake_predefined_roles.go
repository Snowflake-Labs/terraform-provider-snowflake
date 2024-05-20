package snowflakeroles

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

var (
	Orgadmin               = sdk.NewAccountObjectIdentifier("ORGADMIN")
	Accountadmin           = sdk.NewAccountObjectIdentifier("ACCOUNTADMIN")
	GenericScimProvisioner = sdk.NewAccountObjectIdentifier("GENERIC_SCIM_PROVISIONER")
)
