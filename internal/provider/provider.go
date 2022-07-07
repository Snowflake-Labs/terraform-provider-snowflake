package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/snowflakedb/terraform-provider-snowflake/sdk"

	"fmt"
)

type provider struct {
	client *sdk.Client
}

func New() tfsdk.Provider {
	return &provider{}
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "A provider for managing Snowflake",
		Attributes: map[string]tfsdk.Attribute{
			"account": {
				Description: "The name of the Snowflake account. Can also come from the `SNOWFLAKE_ACCOUNT` environment variable.",
				Type:        types.StringType,
				Required:    true,
			},
			"user": {
				Type:        types.StringType,
				Description: "Username for username+password authentication. Can come from the `SNOWFLAKE_USER` environment variable.",
				Required:    true,
			},
			"password": {
				Type:        types.StringType,
				Description: "Password for username+password auth. Cannot be used with `browser_auth` or `private_key_path`. Can be source from `SNOWFLAKE_PASSWORD` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"role": {
				Type:        types.StringType,
				Description: "Snowflake role to use for operations. If left unset, default role for user will be used. Can come from the `SNOWFLAKE_ROLE` environment variable.",
				Optional:    true,
			},
			"region": {
				Type:        types.StringType,
				Description: "[Snowflake region](https://docs.snowflake.com/en/user-guide/intro-regions.html) to use. Can be source from the `SNOWFLAKE_REGION` environment variable.",
				Optional:    true,
			},
			"host": {
				Type:        types.StringType,
				Description: "Supports passing in a custom host value to the snowflake go driver for use with privatelink.",
				Optional:    true,
			},
			"warehouse": {
				Type:        types.StringType,
				Description: "Sets the default warehouse. Optional. Can be sourced from SNOWFLAKE_WAREHOUSE enviornment variable.",
				Optional:    true,
			},
		},
	}, nil
}

type providerData struct {
	Account   types.String `tfsdk:"account"`
	User      types.String `tfsdk:"user"`
	Password  types.String `tfsdk:"password"`
	Role      types.String `tfsdk:"role"`
	Region    types.String `tfsdk:"region"`
	Host      types.String `tfsdk:"host"`
	Warehouse types.String `tfsdk:"warehouse"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	if diags := req.Config.Get(ctx, &config); diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	// interpolation not allowed in provider block
	if config.Account.Unknown {
		addCannotInterpolateInProviderBlockError(resp, "account")
		return
	}
	if config.User.Unknown {
		addCannotInterpolateInProviderBlockError(resp, "user")
		return
	}
	if config.Role.Unknown {
		addCannotInterpolateInProviderBlockError(resp, "role")
		return
	}
	if config.Region.Unknown {
		addCannotInterpolateInProviderBlockError(resp, "region")
		return
	}
	if config.Host.Unknown {
		addCannotInterpolateInProviderBlockError(resp, "host")
		return
	}
	if config.Warehouse.Unknown {
		addCannotInterpolateInProviderBlockError(resp, "warehouse")
		return
	}

	// if unset, fallback to env
	if config.Account.Null {
		config.Account.Value = os.Getenv("SNOWFLAKE_ACCOUNT")
	}
	if config.User.Null {
		config.User.Value = os.Getenv("SNOWFLAKE_USER")
	}
	if config.Role.Null {
		config.Role.Value = os.Getenv("SNOWFLAKE_ROLE")
	}
	if config.Region.Null {
		config.Region.Value = os.Getenv("SNOWFLAKE_REGION")
	}
	if config.Host.Null {
		config.Host.Value = os.Getenv("SNOWFLAKE_HOST")
	}
	if config.Warehouse.Null {
		config.Warehouse.Value = os.Getenv("SNOWFLAKE_WAREHOUSE")
	}

	// default values
	if config.Region.Value == "" {
		config.Region.Value = "us-west-2"
	}

	// required if still unset
	if config.Account.Value == "" {
		addAttributeMustBeSetError(resp, "account")
		return
	}
	if config.User.Value == "" {
		addAttributeMustBeSetError(resp, "user")
		return
	}
	if config.Role.Value == "" {
		addAttributeMustBeSetError(resp, "role")
		return
	}

	client, err := sdk.NewClient(&sdk.Config{
		Account:   config.Account.Value,
		User:      config.User.Value,
		Password:  config.Password.Value,
		Role:      config.Role.Value,
		Region:    config.Region.Value,
		Host:      config.Host.Value,
		Warehouse: config.Warehouse.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Error creating snowflake client", err.Error())
		return
	}
	p.client = client
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		//"snowflake_user": userType{},
	}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"snowflake_users": datasourceUsersType{},
	}, nil
}

func errorConvertingProvider(typ interface{}) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("Error converting provider", fmt.Sprintf("An unexpected error was encountered converting the provider. This is always a bug in the provider.\n\nType: %T", typ))
}

func addAttributeMustBeSetError(resp *tfsdk.ConfigureProviderResponse, attr string) {
	resp.Diagnostics.AddAttributeError(
		tftypes.NewAttributePath().WithAttributeName(attr),
		"Invalid provider config",
		fmt.Sprintf("%s must be set.", attr),
	)
}

func addCannotInterpolateInProviderBlockError(resp *tfsdk.ConfigureProviderResponse, attr string) {
	resp.Diagnostics.AddAttributeError(
		tftypes.NewAttributePath().WithAttributeName(attr),
		"Can't interpolate into provider block",
		"Interpolating that value into the provider block doesn't give the provider enough information to run. Try hard-coding the value, instead.",
	)
}
