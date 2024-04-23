package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ApplicationPackageClient struct {
	context *TestClientContext
}

func NewApplicationPackageClient(context *TestClientContext) *ApplicationPackageClient {
	return &ApplicationPackageClient{
		context: context,
	}
}

func (c *ApplicationPackageClient) client() sdk.ApplicationPackages {
	return c.context.client.ApplicationPackages
}
