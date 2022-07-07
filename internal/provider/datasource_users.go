package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/snowflakedb/terraform-provider-snowflake/sdk"
)

type datasourceUsersType struct {
}

func (datasourceUsersType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Users Data Source for the Snowflake Provider",
		Attributes: map[string]tfsdk.Attribute{
			"pattern": {
				Description: "Users pattern for which to return metadata. Please refer to LIKE keyword from Snowflake documentation [doc](https://docs.snowflake.com/en/sql-reference/sql/show-users.html#parameters)",
				Type:        types.StringType,
				Required:    true,
			},
			"users": {
				Computed:    true,
				Description: "List of users matching the pattern",
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Computed: true,
					},
					"login_name": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"comment": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"disabled": {
						Type:     types.BoolType,
						Optional: true,
						Computed: true,
					},
					"default_warehouse": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"default_namespace": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"default_role": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"default_secondary_roles": {
						Type: types.SetType{
							ElemType: types.StringType,
						},
						Optional: true,
						Computed: true,
					},
					"has_rsa_public_key": {
						Type:     types.BoolType,
						Computed: true,
					},
					"email": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"display_name": {
						Type:     types.StringType,
						Computed: true,
						Optional: true,
					},
					"first_name": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
					"last_name": {
						Type:     types.StringType,
						Optional: true,
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

func (d datasourceUsersType) NewDataSource(_ context.Context, prov tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, ok := prov.(*provider)
	if !ok {
		return nil, diag.Diagnostics{errorConvertingProvider(d)}
	}
	return dsUsers{
		p: provider,
	}, nil
}

type dsUsers struct {
	p *provider
}

type dsUsersData struct {
	Pattern string `tfsdk:"pattern"`
	Users []dsUserData `tfsdk:"users"`
}
type dsUserData struct {
	Name string `tfsdk:"name"`
	LoginName string `tfsdk:"login_name"`
	Comment string `tfsdk:"comment"`
	Disabled bool `tfsdk:"disabled"`
	DefaultWarehouse string `tfsdk:"default_warehouse"`
	DefaultNamespace string `tfsdk:"default_namespace"`
	DefaultRole string `tfsdk:"default_role"`
	DefaultSecondaryRoles []string `tfsdk:"default_secondary_roles"`
	HasRsaPublicKey bool `tfsdk:"has_rsa_public_key"`
	Email string `tfsdk:"email"`
	DisplayName string `tfsdk:"display_name"`
	FirstName string `tfsdk:"first_name"`
	LastName string `tfsdk:"last_name"`
}

func (d dsUsers) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data dsUsersData
	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	userList, err := d.p.client.Users.List(ctx, &sdk.UserListOptions {
		Pattern: data.Pattern,
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to list users: %s", err.Error())
		return
	}
	for _, u := range userList {
		data.Users = append(data.Users, dsUserData{
			Name: u.Name,
			LoginName: u.LoginName,
			Comment: u.Comment,
			Disabled: u.Disabled,
			DefaultWarehouse: u.DefaultWarehouse,
			DefaultNamespace: u.DefaultNamespace,
			DefaultRole: u.DefaultRole,
			DefaultSecondaryRoles: u.DefaultSecondaryRoles,
			HasRsaPublicKey: u.HasRsaPublicKey,
			Email: u.Email,
			DisplayName: u.DisplayName,
			FirstName: u.FirstName,
			LastName: u.LastName,
		})
	}
	resp.Diagnostics = resp.State.Set(ctx, &data)
}
