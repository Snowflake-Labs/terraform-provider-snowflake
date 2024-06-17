package schemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var DatabaseShowSchema = &schema.Schema{
	Type:        schema.TypeList,
	Computed:    true,
	Description: "Holds the output of SHOW DATABASES.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_transient": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_current": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retention_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_role_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"describe_output": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Holds the output of DESCRIBE DATABASE.",
				Elem: &schema.Resource{
					Schema: DatabaseDescribeSchema,
				},
			},
			"parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Holds the output of SHOW PARAMETERS FOR DATABASE.",
				Elem: &schema.Resource{
					Schema: ParameterSchema,
				},
			},
		},
	},
}

var DatabaseDescribeSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
}
