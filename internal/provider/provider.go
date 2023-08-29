package provider

import (
	"context"
	"os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/snowflakedb/gosnowflake"
)

// Ensure SnowflakeProvider satisfies various provider interfaces.
var _ provider.Provider = &SnowflakeProvider{}

// SnowflakeProvider defines the provider implementation.
type SnowflakeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SnowflakeProviderModel describes the provider data model.
type SnowflakeProviderModel struct {
	Account  types.String `tfsdk:"account"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Role     types.String `tfsdk:"role"`
}

func (p *SnowflakeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "snowflake"
	resp.Version = p.version
}

func (p *SnowflakeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account": schema.StringAttribute{
				Description: "",
				//Description:         "An account identifier which uniquely identifies a Snowflake account within your organization, as well as throughout the global network of Snowflake-supported cloud platforms and cloud regions. Can also be set with the SNOWFLAKE_ACCOUNT environment variable. Required unless using profile.",
				//MarkdownDescription: "An [account identifier](https://docs.snowflake.com/en/user-guide/admin-account-identifier) which uniquely identifies a Snowflake account within your organization, as well as throughout the global network of Snowflake-supported cloud platforms and cloud regions. Can also be set with the `SNOWFLAKE_ACCOUNT` environment variable. Required unless using `profile`.",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				Description: "",
				//Description:         "The username to use to authenticate with Snowflake. Can also be set with the SNOWFLAKE_USER environment variable. Required unless using profile.",
				//MarkdownDescription: "The username to use to authenticate with Snowflake. Can also be set with the `SNOWFLAKE_USER` environment variable. Required unless using `profile`.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				Description: "",
				//Description:         "The password to use to authenticate with Snowflake. Can also be set with the SNOWFLAKE_PASSWORD environment variable. Required unless using profile.",
				//MarkdownDescription: "The password to use to authenticate with Snowflake. Can also be set with the `SNOWFLAKE_PASSWORD` environment variable. Required unless using `profile`.",
				Optional:            true,
				Sensitive:           true,
			},
			"role": schema.StringAttribute{
				Description: "",
				//Description:         "The role to use to authenticate with Snowflake. Can also be set with the SNOWFLAKE_ROLE environment variable. Required unless using profile.",
				//MarkdownDescription: "The role to use to authenticate with Snowflake. Can also be set with the `SNOWFLAKE_ROLE` environment variable. Required unless using `profile`.",
				Optional:            true,
			},
		},
	}
}

func (p *SnowflakeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Check environment variables
	account := os.Getenv("SNOWFLAKE_ACCOUNT")
	user := os.Getenv("SNOWFLAKE_USER")
	password := os.Getenv("SNOWFLAKE_PASSWORD")
	role := os.Getenv("SNOWFLAKE_ROLE")

	var data SnowflakeProviderModel

	// Read configuration data into model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if data.Account.ValueString() != "" {
		account = data.Account.ValueString()
	}

	if data.User.ValueString() != "" {
		user = data.User.ValueString()
	}

	if data.Password.ValueString() != "" {
		password = data.Password.ValueString()
	}

	if data.Role.ValueString() != "" {
		role = data.Role.ValueString()
	}

	if account == "" {
		resp.Diagnostics.AddError("Missing Snowflake Account Identifier", "While configuring the provider, the Snowflake Account identifier was not found in the SNOWFLAKE_ACCOUNT environment variable or provider configuration block `account` attribute.")
	}

	config := &gosnowflake.Config{
		Account:  account,
		User:     user,
		Password: password,
		Role:     role,
	}

	// Example client configuration for data sources and resources
	client, err := sdk.NewClient(config)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Snowflake client", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	providerData := &ProviderData{
		client: client,
	}
	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

type ProviderData struct {
	client *sdk.Client
}

func (p *SnowflakeProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{
		/*providervalidator.Conflicting(
			path.MatchRoot("account"),
			path.MatchRoot("profile"),
		),*/
	}
}

func (p *SnowflakeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceMonitorResource,
	}
}

func (p *SnowflakeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// NewDatabaseDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SnowflakeProvider{
			version: version,
		}
	}
}
