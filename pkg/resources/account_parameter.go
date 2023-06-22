package resources

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/exp/maps"
)

var accountParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Name of account parameter. Valid values are those in [account parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#account-parameters).",
		ValidateFunc: validation.StringInSlice(maps.Keys(snowflake.GetParameterDefaults(snowflake.ParameterTypeAccount)), false),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of account parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
}

func AccountParameter() *schema.Resource {
	return &schema.Resource{
		Create: CreateAccountParameter,
		Read:   ReadAccountParameter,
		Update: UpdateAccountParameter,
		Delete: DeleteAccountParameter,

		Schema: accountParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAccountParameter implements schema.CreateFunc.
func CreateAccountParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	parameter := sdk.AccountParameter(key)

	opts, err := setAccountParameter(parameter, value)
	if err != nil {
		return err
	}
	err = client.Accounts.Alter(ctx, opts)
	if err != nil {
		return err
	}
	d.SetId(key)
	return ReadAccountParameter(d, meta)
}

func setAccountParameter(parameter sdk.AccountParameter, value string) (*sdk.AlterAccountOptions, error) {
	opts := sdk.AlterAccountOptions{Set: &sdk.AccountSet{Parameters: &sdk.AccountLevelParameters{AccountParameters: &sdk.AccountParameters{}}}}
	switch parameter {
	case sdk.AccountParameterAllowClientMFACaching:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.AllowClientMFACaching = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.AllowClientMFACaching = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("ALLOW_CLIENT_MFA_CACHING session parameter is a boolean value, got: %v", value)
		}
	case sdk.AccountParameterAllowIDToken:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.AllowIDToken = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.AllowIDToken = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("ALLOW_ID_TOKEN session parameter is a boolean value, got: %v", value)
		}

	case sdk.AccountParameterClientEncryptionKeySize:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("CLIENT_ENCRYPTION_KEY_SIZE session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.AccountParameters.ClientEncryptionKeySize = sdk.Pointer(v)
	case sdk.AccountParameterEnableInternalStagesPrivatelink:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.EnableInternalStagesPrivatelink = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.EnableInternalStagesPrivatelink = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("ENABLE_INTERNAL_STAGES_PRIVATELINK session parameter is a boolean value, got: %v", value)
		}
	case sdk.AccountParameterEventTable:
		opts.Set.Parameters.AccountParameters.EventTable = &value
	case sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.ExternalOAuthAddPrivilegedRolesToBlockedList = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.ExternalOAuthAddPrivilegedRolesToBlockedList = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST session parameter is a boolean value, got: %v", value)
		}
	case sdk.AccountParameterInitialReplicationSizeLimitInTB:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("INITIAL_REPLICATION_SIZE_LIMIT_IN_TB session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.AccountParameters.InitialReplicationSizeLimitInTB = sdk.Pointer(v)

	case sdk.AccountParameterMinDataRetentionTimeInDays:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("MIN_DATA_RETENTION_TIME_IN_DAYS session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.AccountParameters.MinDataRetentionTimeInDays = sdk.Pointer(v)
	case sdk.AccountParameterNetworkPolicy:
		opts.Set.Parameters.AccountParameters.NetworkPolicy = &value
	case sdk.AccountParameterPeriodicDataRekeying:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.PeriodicDataRekeying = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.PeriodicDataRekeying = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("PERIODIC_DATA_REKEYING session parameter is a boolean value, got: %v", value)
		}
	case sdk.AccountParameterPreventUnloadToInlineURL:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.PreventUnloadToInlineURL = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.PreventUnloadToInlineURL = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("PREVENT_UNLOAD_TO_INLINE_URL session parameter is a boolean value, got: %v", value)
		}
	case sdk.AccountParameterPreventUnloadToInternalStages:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.PreventUnloadToInternalStages = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.PreventUnloadToInternalStages = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("PREVENT_UNLOAD_TO_INTERNAL_STAGES session parameter is a boolean value, got: %v", value)
		}
	case sdk.AccountParameterRequireStorageIntegrationForStageCreation:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.RequireStorageIntegrationForStageCreation = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.RequireStorageIntegrationForStageCreation = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION session parameter is a boolean value, got: %v", value)
		}
	case sdk.AccountParameterRequireStorageIntegrationForStageOperation:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.RequireStorageIntegrationForStageOperation = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.RequireStorageIntegrationForStageOperation = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION session parameter is a boolean value, got: %v", value)
		}
	case sdk.AccountParameterSSOLoginPage:
		if value == "true" {
			opts.Set.Parameters.AccountParameters.SSOLoginPage = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.AccountParameters.SSOLoginPage = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("SSO_LOGIN_PAGE session parameter is a boolean value, got: %v", value)
		}
	default:
		return nil, fmt.Errorf("Invalid account parameter: %v", string(parameter))
	}
	return &opts, nil
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	parameterName := d.Id()
	parameter, err := client.Sessions.ShowAccountParameter(ctx, sdk.AccountParameter(parameterName))
	if err != nil {
		return fmt.Errorf("error reading account parameter err = %w", err)
	}
	err = d.Set("value", parameter.Value)
	if err != nil {
		return fmt.Errorf("error setting account parameter err = %w", err)
	}
	return nil
}

// UpdateAccountParameter implements schema.UpdateFunc.
func UpdateAccountParameter(d *schema.ResourceData, meta interface{}) error {
	return CreateAccountParameter(d, meta)
}

// DeleteAccountParameter implements schema.DeleteFunc.
func DeleteAccountParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	parameter := sdk.AccountParameter(key)

	defaultParameter, err := client.Sessions.ShowAccountParameter(ctx, sdk.AccountParameter(key))
	if err != nil {
		return err
	}
	defaultValue := defaultParameter.Default
	opts, err := setAccountParameter(parameter, defaultValue)
	if err != nil {
		return err
	}
	err = client.Accounts.Alter(ctx, opts)
	if err != nil {
		return fmt.Errorf("error creating account parameter err = %w", err)
	}

	d.SetId("")
	return nil
}
