package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
			StateContext: ImportName[sdk.AccountObjectIdentifier],
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

		if _, ok := enableFailoverConfig["to_accounts"]; !ok || len(enableFailoverConfig) == 0 {
			return diag.FromErr(fmt.Errorf("The %s Connection 'to_accounts' list field is required when enable_failover is set", id.FullyQualifiedName()))
			// return ReadContextConnection(ctx, d, meta)
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

	errs := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ConnectionToSchema(connection)}),
		d.Set("as_replica_of", connection.Primary),
		d.Set("comment", connection.Comment),
		d.Set("is_primary", connection.IsPrimary),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	sessionDetails, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	currentAccountIdentifier := sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName)

	failoverAllowedToAccounts := make([]sdk.AccountIdentifier, 0)
	for _, allowedAccount := range connection.FailoverAllowedToAccounts {
		allowedAccountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(strings.TrimSpace(allowedAccount))
		if currentAccountIdentifier.FullyQualifiedName() == allowedAccountIdentifier.FullyQualifiedName() {
			continue
		}
		failoverAllowedToAccounts = append(failoverAllowedToAccounts, allowedAccountIdentifier)
	}

	enableToAccounts := make([]map[string]any, 0)
	for _, allowedAccount := range failoverAllowedToAccounts {
		enableToAccounts = append(enableToAccounts, map[string]any{
			"account_identifier": strings.ReplaceAll(allowedAccount.FullyQualifiedName(), `"`, ""),
		})
	}

	if len(enableToAccounts) == 0 {
		err := d.Set("enable_failover", []any{})
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err := d.Set("enable_failover", map[string]any{
			"to_accounts": enableToAccounts,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func UpdateContextConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	connectionSetRequest := new(sdk.SetConnectionRequest)
	connectionUnsetRequest := new(sdk.UnsetConnectionRequest)

	if d.HasChange("enable_failover") {
		before, after := d.GetChange("enable_failover")

		getFailoverToAccounts := func(failoverConfig []any) []sdk.AccountIdentifier {
			failoverEnabledToAccounts := make([]sdk.AccountIdentifier, 0)

			for _, enableFailoverConfigMap := range failoverConfig {
				enableFailoverConfig := enableFailoverConfigMap.(map[string]any)
				for _, toAccountsMap := range enableFailoverConfig["to_accounts"].([]any) {
					enableToAccounts := toAccountsMap.(map[string]any)
					accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(enableToAccounts["account_identifier"].(string))

					failoverEnabledToAccounts = append(failoverEnabledToAccounts, accountIdentifier)
				}
			}
			return failoverEnabledToAccounts
		}

		beforeFailover := getFailoverToAccounts(before.([]any))
		afterFailover := getFailoverToAccounts(after.([]any))

		addedFailovers, removedFailovers := ListDiff(beforeFailover, afterFailover)

		if len(addedFailovers) > 0 {
			err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
				WithEnableConnectionFailover(*sdk.NewEnableConnectionFailoverRequest().
					WithToAccounts(addedFailovers),
				),
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(removedFailovers) > 0 {
			err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
				WithDisableConnectionFailover(*sdk.NewDisableConnectionFailoverRequest().
					WithToAccounts(*sdk.NewToAccountsRequest().
						WithAccounts(removedFailovers),
					),
				),
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("is_primary") {
		is_primary := d.Get("is_primary").(bool)
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithPrimary(is_primary))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			connectionSetRequest.WithComment(comment)
		} else {
			connectionUnsetRequest.WithComment(true)
		}
	}

	if (*connectionSetRequest != sdk.SetConnectionRequest{}) {
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).WithSet(*connectionSetRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (*connectionUnsetRequest != sdk.UnsetConnectionRequest{}) {
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).WithUnset(*connectionUnsetRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextConnection(ctx, d, meta)
}

func DeleteContextConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Connections.Drop(ctx, sdk.NewDropConnectionRequest(id).WithIfExists(true))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
