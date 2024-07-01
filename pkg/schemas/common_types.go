package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ParameterListSchema represents Snowflake parameter object.
var ParameterListSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: ShowParameterSchema,
	},
}

// DescribePropertyListSchema represents Snowflake property object returned by DESCRIBE query.
var DescribePropertyListSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: ShowSecurityIntegrationPropertySchema,
	},
}
