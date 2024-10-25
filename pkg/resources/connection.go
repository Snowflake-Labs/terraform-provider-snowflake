package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var connectionSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the connection. Must start with an alphabetic character and may only contain letters, decimal digits (0-9), and underscores (_). For a primary connection, the name must be unique across connection names and account names in the organization. For a secondary connection, the name must match the name of its primary connection."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"as_replica_of": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the identifier for a primary connection from which to create a replica (i.e. a secondary connection).",
	},
	"enable_failover": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Enables failover for given connection.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"to_accounts": {
					Type:        schema.TypeList,
					Required:    true,
					Description: "Specifies a list of accounts in your organization where a secondary connection for this primary connection can be promoted to serve as the primary connection. Include your organization name for each account in the list.",
					MinItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"account_identifier": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies account identifier for which replication should be enabled. The account identifiers should be in the form of `\"<organization_name>\".\"<account_name>\"`.",
							},
						},
					},
				},
				/*
					"ignore_edition_check": {
						Type:     schema.TypeBool,
						Optional: true,
						Description: "Allows replicating data to accounts on lower editions in either of the following scenarios: " +
							"1. The primary database is in a Business Critical (or higher) account but one or more of the accounts approved for replication are on lower editions. Business Critical Edition is intended for Snowflake accounts with extremely sensitive data. " +
							"2. The primary database is in a Business Critical (or higher) account and a signed business associate agreement is in place to store PHI data in the account per HIPAA and HITRUST regulations, but no such agreement is in place for one or more of the accounts approved for replication, regardless if they are Business Critical (or higher) accounts. " +
							"Both scenarios are prohibited by default in an effort to help prevent account administrators for Business Critical (or higher) accounts from inadvertently replicating sensitive data to accounts on lower editions.",
					},
				*/
			},
		},
	},
	"is_primary": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Promote connection to serve as primary connection.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the connection.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CONNECTIONS` for the given secret.",
		Elem: &schema.Resource{
			Schema: schemas.ShowConnectionSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func Connection() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextConnection,
		ReadContext:   ReadContextConnection,
		UpdateContext: UpdateContextConnection,
		DeleteContext: DeleteContextConnection,
		Description:   "Resource used to manage connections. For more information, check [connection documentation](https://docs.snowflake.com/en/sql-reference/sql/create-connection.html).",
		Schema:        connectionSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportConnection,
		},
	}
}

func CreateContextConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateConnectionRequest(id)

	if v, ok := d.GetOk("as_replica_of"); ok {
		if externalObjectId, err := sdk.ParseExternalObjectIdentifier(v.(string)); err != nil {
			request.WithAsReplicaOf(*sdk.NewAsReplicaOfRequest(externalObjectId))
		} else {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err = client.Connections.Create(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("is_primary"); ok {
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithPrimary(v.(bool)),
		)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if v, ok := d.GetOk("enable_failover"); ok {
		enableFailoverConfig := v.([]any)[0].(map[string]any)

		if _, ok := enableFailoverConfig["to_accounts"].([]any); !ok || len(enableFailoverConfig) == 0 {
			return diag.FromErr(fmt.Errorf("The %s Connection 'to_accounts' list field is required when enable_failover is set", id.FullyQualifiedName()))
		}
		enableFailoverToAccountsConfig := enableFailoverConfig["to_accounts"].([]any)

		enableFailoverToAccountsList := make([]sdk.AccountIdentifier, 0)
		for _, enableToAccount := range enableFailoverToAccountsConfig {
			accountInConfig := enableToAccount.(map[string]any)
			accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(accountInConfig["account_identifier"].(string))

			enableFailoverToAccountsList = append(enableFailoverToAccountsList, accountIdentifier)
		}

		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithEnableConnectionFailover(*sdk.NewEnableConnectionFailoverRequest().
				WithToAccounts(enableFailoverToAccountsList)),
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextConnection(ctx, d, meta)
}

func ReadContextConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	connection, err := client.Connections.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve connection. Target object not found. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Connection name: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve connection.",
				Detail:   fmt.Sprintf("Connection name: %s, Err: %s", id.FullyQualifiedName(), err),
			},
		}
	}

	return diag.FromErr(errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("as_replica_of", connection.Primary),
		d.Set("comment", connection.Comment),
		d.Set("is_primary", connection.IsPrimary),

		// should use some func to map failover accounts to schema
		// if err := d.Set("enable_failover", connection.FailoverAllowedToAccounts); err != nil {}
	))
}
