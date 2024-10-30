package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
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
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "Specifies the identifier for a primary connection from which to create a replica (i.e. a secondary connection).",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"is_primary": {
		Type:         schema.TypeBool,
		Optional:     true,
		RequiredWith: []string{"as_replica_of"},
	},
	"enable_failover_to_accounts": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Enables failover for given connection to provided accounts. Specifies a list of accounts in your organization where a secondary connection for this primary connection can be promoted to serve as the primary connection. Include your organization name for each account in the list.",
		MinItems:    1,
		Elem:        &schema.Schema{Type: schema.TypeString},
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
		ReadContext:   ReadContextConnection(true),
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
		externalObjectId, err := sdk.ParseExternalObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithAsReplicaOf(externalObjectId)
	}

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err = client.Connections.Create(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	if v, ok := d.GetOk("is_primary"); ok && v.(bool) {
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithPrimary(v.(bool)),
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("enable_failover_to_accounts"); ok {
		enableFailoverConfig := v.([]any)

		enableFailoverToAccountsList := make([]sdk.AccountIdentifier, 0)
		for _, enableToAccount := range enableFailoverConfig {
			accountInConfig := enableToAccount.(string)
			accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(accountInConfig)

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

	return ReadContextConnection(false)(ctx, d, meta)
}

func ReadContextConnection(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

		if withExternalChangesMarking {
			if err := handleExternalChangesToObjectInShow(d,
				outputMapping{"is_primary", "is_primary", connection.IsPrimary, connection.IsPrimary, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			setStateToValuesFromConfig(d, connectionSchema, []string{"is_primary"}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.ConnectionToSchema(connection)}),
			d.Set("comment", connection.Comment),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		sessionDetails, err := client.ContextFunctions.CurrentSessionDetails(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		currentAccountIdentifier := sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName)

		enableFailoverToAccounts := make([]string, 0)
		for _, allowedAccount := range connection.FailoverAllowedToAccounts {
			allowedAccountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(allowedAccount)
			if currentAccountIdentifier.FullyQualifiedName() == allowedAccountIdentifier.FullyQualifiedName() {
				continue
			}
			enableFailoverToAccounts = append(enableFailoverToAccounts, allowedAccountIdentifier.Name())
		}

		if len(enableFailoverToAccounts) == 0 {
			err := d.Set("enable_failover_to_accounts", []any{})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := d.Set("enable_failover_to_accounts", enableFailoverToAccounts)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	}
}

func UpdateContextConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	connectionSetRequest := new(sdk.ConnectionSetRequest)
	connectionUnsetRequest := new(sdk.ConnectionUnsetRequest)

	if d.HasChange("enable_failover_to_accounts") {
		before, after := d.GetChange("enable_failover_to_accounts")

		getFailoverToAccounts := func(failoverConfig []any) []sdk.AccountIdentifier {
			failoverEnabledToAccounts := make([]sdk.AccountIdentifier, 0)

			for _, allowedAccount := range failoverConfig {
				accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(allowedAccount.(string))
				failoverEnabledToAccounts = append(failoverEnabledToAccounts, accountIdentifier)
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
		if is_primary := d.Get("is_primary").(bool); is_primary {
			err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).WithPrimary(is_primary))
			if err != nil {
				return diag.FromErr(err)
			}
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

	if (*connectionSetRequest != sdk.ConnectionSetRequest{}) {
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).WithSet(*connectionSetRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (*connectionUnsetRequest != sdk.ConnectionUnsetRequest{}) {
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).WithUnset(*connectionUnsetRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextConnection(false)(ctx, d, meta)
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
