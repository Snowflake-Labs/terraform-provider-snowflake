package ids

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

// TODO [SNOW-1827335]: there are similar non existing ids in setup_test.go in integration tests, consider merging them (they use different parents - acc vs int)
var (
	NonExistingAccountObjectIdentifier = sdk.NewAccountObjectIdentifier("does_not_exist")
)
