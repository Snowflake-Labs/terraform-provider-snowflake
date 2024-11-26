package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sessionParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Name of session parameter. Valid values are those in [session parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#session-parameters).",
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of session parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
	"on_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the session parameter will be set on the account level.",
	},
	"user": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The user to set the session parameter for. Required if on_account is false",
	},
}

func SessionParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.SessionParameter, CreateSessionParameter),
		ReadContext:   TrackingReadWrapper(resources.SessionParameter, ReadSessionParameter),
		UpdateContext: TrackingUpdateWrapper(resources.SessionParameter, UpdateSessionParameter),
		DeleteContext: TrackingDeleteWrapper(resources.SessionParameter, DeleteSessionParameter),

		Schema: sessionParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSessionParameter implements schema.CreateFunc.
func CreateSessionParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	onAccount := d.Get("on_account").(bool)
	user := d.Get("user").(string)
	parameter := sdk.SessionParameter(key)

	var err error
	if onAccount {
		err := client.Parameters.SetSessionParameterOnAccount(ctx, parameter, value)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		if user == "" {
			return diag.FromErr(fmt.Errorf("user is required if on_account is false"))
		}
		userId := sdk.NewAccountObjectIdentifier(user)
		err = client.Parameters.SetSessionParameterOnUser(ctx, userId, parameter, value)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating session parameter err = %w", err))
		}
	}

	d.SetId(key)

	return ReadSessionParameter(ctx, d, meta)
}

// ReadSessionParameter implements schema.ReadFunc.
func ReadSessionParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parameter := d.Id()

	onAccount := d.Get("on_account").(bool)
	var err error
	var p *sdk.Parameter
	if onAccount {
		p, err = client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(parameter))
	} else {
		user := d.Get("user").(string)
		userId := sdk.NewAccountObjectIdentifier(user)
		p, err = client.Parameters.ShowUserParameter(ctx, sdk.UserParameter(parameter), userId)
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading session parameter err = %w", err))
	}
	err = d.Set("value", p.Value)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting session parameter err = %w", err))
	}
	return nil
}

// UpdateSessionParameter implements schema.UpdateFunc.
func UpdateSessionParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return CreateSessionParameter(ctx, d, meta)
}

// DeleteSessionParameter implements schema.DeleteFunc.
func DeleteSessionParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)

	onAccount := d.Get("on_account").(bool)
	parameter := sdk.SessionParameter(key)

	if onAccount {
		defaultParameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
		if err != nil {
			return diag.FromErr(err)
		}
		defaultValue := defaultParameter.Default
		err = client.Parameters.SetSessionParameterOnAccount(ctx, parameter, defaultValue)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating session parameter err = %w", err))
		}
	} else {
		user := d.Get("user").(string)
		if user == "" {
			return diag.FromErr(fmt.Errorf("user is required if on_account is false"))
		}
		userId := sdk.NewAccountObjectIdentifier(user)
		defaultParameter, err := client.Parameters.ShowSessionParameter(ctx, sdk.SessionParameter(key))
		if err != nil {
			return diag.FromErr(err)
		}
		defaultValue := defaultParameter.Default
		err = client.Parameters.SetSessionParameterOnUser(ctx, userId, parameter, defaultValue)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error deleting session parameter err = %w", err))
		}
	}

	d.SetId(key)
	return nil
}
