// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowPolicyReferenceSchema represents output of SHOW query for the single PolicyReference.
var ShowPolicyReferenceSchema = map[string]*schema.Schema{
	"policy_db": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"policy_schema": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"policy_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"policy_kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"ref_database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"ref_schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"ref_entity_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"ref_entity_domain": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"ref_column_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"ref_arg_column_names": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"tag_database": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"tag_schema": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"tag_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"policy_status": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowPolicyReferenceSchema

func PolicyReferenceToSchema(policyReference *sdk.PolicyReference) map[string]any {
	policyReferenceSchema := make(map[string]any)
	if policyReference.PolicyDb != nil {
		policyReferenceSchema["policy_db"] = policyReference.PolicyDb
	}
	if policyReference.PolicySchema != nil {
		policyReferenceSchema["policy_schema"] = policyReference.PolicySchema
	}
	policyReferenceSchema["policy_name"] = policyReference.PolicyName
	policyReferenceSchema["policy_kind"] = string(policyReference.PolicyKind)
	if policyReference.RefDatabaseName != nil {
		policyReferenceSchema["ref_database_name"] = policyReference.RefDatabaseName
	}
	if policyReference.RefSchemaName != nil {
		policyReferenceSchema["ref_schema_name"] = policyReference.RefSchemaName
	}
	policyReferenceSchema["ref_entity_name"] = policyReference.RefEntityName
	policyReferenceSchema["ref_entity_domain"] = policyReference.RefEntityDomain
	if policyReference.RefColumnName != nil {
		policyReferenceSchema["ref_column_name"] = policyReference.RefColumnName
	}
	if policyReference.RefArgColumnNames != nil {
		policyReferenceSchema["ref_arg_column_names"] = policyReference.RefArgColumnNames
	}
	if policyReference.TagDatabase != nil {
		policyReferenceSchema["tag_database"] = policyReference.TagDatabase
	}
	if policyReference.TagSchema != nil {
		policyReferenceSchema["tag_schema"] = policyReference.TagSchema
	}
	if policyReference.TagName != nil {
		policyReferenceSchema["tag_name"] = policyReference.TagName
	}
	if policyReference.PolicyStatus != nil {
		policyReferenceSchema["policy_status"] = policyReference.PolicyStatus
	}
	return policyReferenceSchema
}

var _ = PolicyReferenceToSchema
