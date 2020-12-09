package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func mergeSchemas(schemaCollections ...map[string]*schema.Resource) map[string]*schema.Resource {
	out := map[string]*schema.Resource{}
	for _, schemaCollection := range schemaCollections {
		for name, s := range schemaCollection {
			out[name] = s
		}
	}
	return out
}
