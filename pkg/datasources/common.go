package datasources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


func handleLike(d *schema.ResourceData, setField **sdk.Like) {
	if likePattern, ok := d.GetOk("like"); ok {
		*setField = &sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		}
	}
}

func handleLimitFrom(d *schema.ResourceData, setField **sdk.LimitFrom) {
	if v, ok := d.GetOk("limit"); ok {
		l := v.([]any)[0].(map[string]any)
		limit := &sdk.LimitFrom{}
		if v, ok := l["rows"]; ok {
			rows := v.(int)
			limit.Rows = sdk.Int(rows)
		}
		if v, ok := l["from"]; ok {
			from := v.(string)
			limit.From = sdk.String(from)
		}
		*setField = limit
	}
}

func handleExtendedIn(d *schema.ResourceData, setField **sdk.ExtendedIn) error {
	if v, ok := d.GetOk("in"); ok {
		in := v.([]any)[0].(map[string]any)
		if v, ok := in["account"]; ok && v.(bool) {
			*setField = &sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}}
		}
		if v, ok := in["database"]; ok {
			if database := v.(string); database != "" {
				*setField = &sdk.ExtendedIn{In: sdk.In{Database: sdk.NewAccountObjectIdentifier(database)}}
			}
		}
		if v, ok := in["schema"]; ok {
			if schema := v.(string); schema != "" {
				schemaId, err := sdk.ParseDatabaseObjectIdentifier(schema)
				if err != nil {
					return err
				}
				*setField = &sdk.ExtendedIn{In: sdk.In{Schema: schemaId}}
			}
		}
		if v, ok := in["application"]; ok {
			if application := v.(string); application != "" {
				*setField = &sdk.ExtendedIn{Application: sdk.NewAccountObjectIdentifier(application)}
			}
		}
		if v, ok := in["application_package"]; ok {
			if applicationPackage := v.(string); applicationPackage != "" {
				*setField = &sdk.ExtendedIn{ApplicationPackage: sdk.NewAccountObjectIdentifier(applicationPackage)}
			}
		}
	}
	return nil
}
