package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestAPIIntegration(t *testing.T) {
	r := require.New(t)
	builder := snowflake.NewAPIIntegrationBuilder("aws_api")
	r.NotNil(builder)

	q := builder.Show()
	r.Equal("SHOW API INTEGRATIONS LIKE 'aws_api'", q)

	c := builder.Create()

	c.SetRaw(`API_PROVIDER=aws_private_api_gateway`)
	c.SetString(`api_aws_role_arn`, "arn:aws:iam::xxxx:role/snowflake-execute-externalfunc-privendpoint-role")
	c.SetStringList(`api_allowed_prefixes`, []string{"https://123456.execute-api.us-west-2.amazonaws.com/prod/", "https://123456.execute-api.us-west-2.amazonaws.com/test/"})
	c.SetBool(`enabled`, true)
	q = c.Statement()

	r.Equal(`CREATE API INTEGRATION "aws_api" API_PROVIDER=aws_private_api_gateway API_AWS_ROLE_ARN='arn:aws:iam::xxxx:role/snowflake-execute-externalfunc-privendpoint-role' API_ALLOWED_PREFIXES=('https://123456.execute-api.us-west-2.amazonaws.com/prod/', 'https://123456.execute-api.us-west-2.amazonaws.com/test/') ENABLED=true`, q)
}
