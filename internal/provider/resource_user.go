package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceUserType struct {
}

func (r resourceUserType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "User Resource for the Snowflake Provider",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "Name of the user. Note that if you do not supply login_name this will be used as login_name. [doc](https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#required-parameters)",
				Type:        types.StringType,
				Required:    true,
			},
			"login_name": {
				Description: "The name users use to log in. If not supplied, snowflake will use name instead.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					caseInsensitive{},
				},
			},
			"comment": {
				Type:     types.StringType,
				Optional: true,
				// TODO validation
			},
			"password": {
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Description: "**WARNING:** this will put the password in the terraform state file. Use carefully.",
				// TODO validation https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#optional-parameters
			},
			"disabled": {
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
			},
			"default_warehouse": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Specifies the virtual warehouse that is active by default for the user’s session upon login.",
			},
			"default_namespace": {
				Type:     types.StringType,
				Optional: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					caseInsensitive{},
				},
				Description: "Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login.",
			},
			"default_role": {
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the role that is active by default for the user’s session upon login.",
			},
			"default_secondary_roles": {
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional:    true,
				Description: "Specifies the set of secondary roles that are active for the user’s session upon login.",
			},
			"rsa_public_key": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Specifies the user’s RSA public key; used for key-pair authentication. Must be on 1 line without header and trailer.",
			},
			"rsa_public_key_2": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Specifies the user’s second RSA public key; used to rotate the public and private keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.",
			},
			"has_rsa_public_key": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "Will be true if user as an RSA key set.",
			},
			"must_change_password": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Specifies whether the user is forced to change their password on next login (including their first/initial login) into the system.",
			},
			"email": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Email address for the user.",
			},
			"display_name": {
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
				Description: "Name displayed for the user in the Snowflake web interface.",
			},
			"first_name": {
				Type:        types.StringType,
				Optional:    true,
				Description: "First name of the user.",
			},
			"last_name": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Last name of the user.",
			},
			//"tag": tagReferenceSchema,

			//    MIDDLE_NAME = <string>
			//    SNOWFLAKE_LOCK = TRUE | FALSE
			//    SNOWFLAKE_SUPPORT = TRUE | FALSE
			//    DAYS_TO_EXPIRY = <integer>
			//    MINS_TO_UNLOCK = <integer>
			//    EXT_AUTHN_DUO = TRUE | FALSE
			//    EXT_AUTHN_UID = <string>
			//    MINS_TO_BYPASS_MFA = <integer>
			//    DISABLE_MFA = TRUE | FALSE
			//    MINS_TO_BYPASS_NETWORK POLICY = <integer>
		},
	}, nil
}

func (r resourceUserType) NewResource(_ context.Context, prov tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, ok := prov.(*provider)
	if !ok {
		return nil, diag.Diagnostics{errorConvertingProvider(r)}
	}
	return resourceUser{
		p: provider,
	}, nil
}

type resourceUser struct {
	p *provider
}

func (r resourceUser) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
}

// Read resource information
func (r resourceUser) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
}

// Update resource
func (r resourceUser) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

// Delete resource
func (r resourceUser) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}


// UserResource -
type UserResource struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	LoginName             types.String `tfsdk:"login_name"`
	Comment               types.String `tfsdk:"comment"`
	Password              types.String `tfsdk:"password"`
	Disabled              types.Bool   `tfsdk:"disabled"`
	DefaultWarehouse      types.String `tfsdk:"default_warehouse"`
	DefaultNamespace      types.String `tfsdk:"default_namespace"`
	DefaultRole           types.String `tfsdk:"default_role"`
	DefaultSecondaryRoles types.Set    `tfsdk:"default_secondary_roles"`
	RSAPublicKey          types.String `tfsdk:"rsa_public_key"`
	RSAPublicKey2         types.String `tfsdk:"rsa_public_key_2"`
	HasRSAPublicKey       types.Bool   `tfsdk:"has_rsa_public_key"`
	MustChangePassword    types.Bool   `tfsdk:"must_change_password"`
	Email                 types.String `tfsdk:"email"`
	DisplayName           types.String `tfsdk:"display_name"`
	FirstName             types.String `tfsdk:"first_name"`
	LastName              types.String `tfsdk:"last_name"`
	//Tag tagReference `tfsdk:"tag"`
}
