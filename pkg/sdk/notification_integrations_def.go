package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var NotificationIntegrationsDef = g.NewInterface(
	"NotificationIntegrations",
	"NotificationIntegration",
	g.KindOfT[AccountObjectIdentifier](),
)
