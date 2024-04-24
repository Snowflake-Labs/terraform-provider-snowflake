package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type FileFormatClient struct {
	context *TestClientContext
}

func NewFileFormatClient(context *TestClientContext) *FileFormatClient {
	return &FileFormatClient{
		context: context,
	}
}

func (c *FileFormatClient) client() sdk.FileFormats {
	return c.context.client.FileFormats
}
