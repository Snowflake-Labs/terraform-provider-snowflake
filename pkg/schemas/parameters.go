package schemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// ParameterListSchema represents Snowflake parameter object.
var ParameterListSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: ShowParameterSchema,
	},
}
