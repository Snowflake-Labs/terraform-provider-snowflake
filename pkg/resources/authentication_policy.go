package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var authenticationPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: blocklistedCharactersFieldDescription("Specifies the identifier for the authentication policy."),
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: blocklistedCharactersFieldDescription("The schema in which to create the authentication policy."),
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: blocklistedCharactersFieldDescription("The database in which to create the authentication policy."),
	},
	"authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: StringInSlice([]string{"ALL", "SAML", "PASSWORD", "OAUTH", "KEYPAIR"}, false),
		},
		Optional:    true,
		Description: "A list of authentication methods that are allowed during login. This parameter accepts one or more of the following values: ALL, SAML, PASSWORD, OAUTH, KEYPAIR.",
	},
	"mfa_authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: StringInSlice([]string{"SAML", "PASSWORD"}, false),
		},
		Optional:    true,
		Description: "A list of authentication methods that enforce multi-factor authentication (MFA) during login. Authentication methods not listed in this parameter do not prompt for multi-factor authentication. Allowed values are SAML and PASSWORD.",
	},
	"mfa_enrollment": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Determines whether a user must enroll in multi-factor authentication. Allowed values are REQUIRED and OPTIONAL. When REQUIRED is specified, Enforces users to enroll in MFA. If this value is used, then the CLIENT_TYPES parameter must include SNOWFLAKE_UI, because Snowsight is the only place users can enroll in multi-factor authentication (MFA).",
		ValidateFunc: validation.StringInSlice([]string{"REQUIRED", "OPTIONAL"}, false),
		Default:      "OPTIONAL",
	},
	"client_types": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: StringInSlice([]string{"ALL", "SNOWFLAKE_UI", "DRIVERS", "SNOWSQL"}, false),
		},
		Optional:    true,
		Description: "A list of clients that can authenticate with Snowflake. If a client tries to connect, and the client is not one of the valid CLIENT_TYPES, then the login attempt fails. Allowed values are ALL, SNOWFLAKE_UI, DRIVERS, SNOWSQL. The CLIENT_TYPES property of an authentication policy is a best effort method to block user logins based on specific clients. It should not be used as the sole control to establish a security boundary.",
	},
	"security_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		Description: "A list of security integrations the authentication policy is associated with. This parameter has no effect when SAML or OAUTH are not in the AUTHENTICATION_METHODS list. All values in the SECURITY_INTEGRATIONS list must be compatible with the values in the AUTHENTICATION_METHODS list. For example, if SECURITY_INTEGRATIONS contains a SAML security integration, and AUTHENTICATION_METHODS contains OAUTH, then you cannot create the authentication policy. To allow all security integrations use ALL as parameter.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the authentication policy.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// AuthenticationPolicy returns a pointer to the resource representing an authentication policy.
func AuthenticationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateContextAuthenticationPolicy,
		ReadContext:   ReadContextAuthenticationPolicy,
		UpdateContext: UpdateContextAuthenticationPolicy,
		DeleteContext: DeleteContextAuthenticationPolicy,

		Schema: authenticationPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	req := sdk.NewCreateAuthenticationPolicyRequest(id)

	// Set optionals
	authenticationMethodsRawList := expandStringList(d.Get("authentication_methods").(*schema.Set).List())
	authenticationMethods := make([]sdk.AuthenticationMethods, len(authenticationMethodsRawList))
	for i, v := range authenticationMethodsRawList {
		authenticationMethods[i] = sdk.AuthenticationMethods{Method: sdk.AuthenticationMethodsOption(v)}
	}
	req.WithAuthenticationMethods(authenticationMethods)

	mfaAuthenticationMethodsRawList := expandStringList(d.Get("mfa_authentication_methods").(*schema.Set).List())
	mfaAuthenticationMethods := make([]sdk.MfaAuthenticationMethods, len(mfaAuthenticationMethodsRawList))
	for i, v := range mfaAuthenticationMethodsRawList {
		mfaAuthenticationMethods[i] = sdk.MfaAuthenticationMethods{Method: sdk.MfaAuthenticationMethodsOption(v)}
	}
	req.WithMfaAuthenticationMethods(mfaAuthenticationMethods)

	if v, ok := d.GetOk("mfa_enrollment"); ok {
		req = req.WithMfaEnrollment(sdk.MfaEnrollmentOption(v.(string)))
	}

	clientTypesRawList := expandStringList(d.Get("client_types").(*schema.Set).List())
	clientTypes := make([]sdk.ClientTypes, len(clientTypesRawList))
	for i, v := range clientTypesRawList {
		clientTypes[i] = sdk.ClientTypes{ClientType: sdk.ClientTypesOption(v)}
	}
	req.WithClientTypes(clientTypes)

	securityIntegrationsRawList := expandStringList(d.Get("security_integrations").(*schema.Set).List())
	securityIntegrations := make([]sdk.SecurityIntegrationsOption, len(securityIntegrationsRawList))
	for i, v := range securityIntegrationsRawList {
		securityIntegrations[i] = sdk.SecurityIntegrationsOption{Name: sdk.NewAccountObjectIdentifier(v)}
	}
	req.WithSecurityIntegrations(securityIntegrations)

	if v, ok := d.GetOk("comment"); ok {
		req = req.WithComment(v.(string))
	}

	client := meta.(*provider.Context).Client
	if err := client.AuthenticationPolicies.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextAuthenticationPolicy(ctx, d, meta)
}

func ReadContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	authenticationPolicy, err := client.AuthenticationPolicies.ShowByID(ctx, id)
	if authenticationPolicy == nil || err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve authentication policy. Target object not found. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Id: %s", d.Id()),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve authentication policy",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	authenticationPolicyDescriptions, err := client.AuthenticationPolicies.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", authenticationPolicy.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("database", authenticationPolicy.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("schema", authenticationPolicy.SchemaName); err != nil {
		return diag.FromErr(err)
	}

	authenticationMethods := make([]string, 0)
	if authenticationMethodsProperty, err := collections.FindFirst(authenticationPolicyDescriptions, func(prop sdk.AuthenticationPolicyDescription) bool { return prop.Property == "AUTHENTICATION_METHODS" }); err == nil {
		authenticationMethods = append(authenticationMethods, sdk.ParseCommaSeparatedStringArray(authenticationMethodsProperty.Value, false)...)
	}
	if err = d.Set("authentication_methods", authenticationMethods); err != nil {
		return diag.FromErr(err)
	}

	mfaAuthenticationMethods := make([]string, 0)
	if mfaAuthenticationMethodsProperty, err := collections.FindFirst(authenticationPolicyDescriptions, func(prop sdk.AuthenticationPolicyDescription) bool {
		return prop.Property == "MFA_AUTHENTICATION_METHODS"
	}); err == nil {
		mfaAuthenticationMethods = append(mfaAuthenticationMethods, sdk.ParseCommaSeparatedStringArray(mfaAuthenticationMethodsProperty.Value, false)...)
	}
	if err = d.Set("mfa_authentication_methods", mfaAuthenticationMethods); err != nil {
		return diag.FromErr(err)
	}

	mfaEnrollment, err := collections.FindFirst(authenticationPolicyDescriptions, func(prop sdk.AuthenticationPolicyDescription) bool { return prop.Property == "MFA_ENROLLMENT" })
	if err == nil {
		if err = d.Set("mfa_enrollment", mfaEnrollment.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	clientTypes := make([]string, 0)
	if clientTypesProperty, err := collections.FindFirst(authenticationPolicyDescriptions, func(prop sdk.AuthenticationPolicyDescription) bool { return prop.Property == "CLIENT_TYPES" }); err == nil {
		clientTypes = append(clientTypes, sdk.ParseCommaSeparatedStringArray(clientTypesProperty.Value, false)...)
	}
	if err = d.Set("client_types", clientTypes); err != nil {
		return diag.FromErr(err)
	}

	securityIntegrations := make([]string, 0)
	if securityIntegrationsProperty, err := collections.FindFirst(authenticationPolicyDescriptions, func(prop sdk.AuthenticationPolicyDescription) bool { return prop.Property == "SECURITY_INTEGRATIONS" }); err == nil {
		securityIntegrations = append(securityIntegrations, sdk.ParseCommaSeparatedStringArray(securityIntegrationsProperty.Value, false)...)
	}
	if err = d.Set("security_integrations", securityIntegrations); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("comment", authenticationPolicy.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func UpdateContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	// change to authentication methods
	if d.HasChange("authentication_methods") {
		authenticationMethods := expandStringList(d.Get("authentication_methods").(*schema.Set).List())
		authenticationMethodsValues := make([]sdk.AuthenticationMethods, len(authenticationMethods))
		for i, v := range authenticationMethods {
			authenticationMethodsValues[i] = sdk.AuthenticationMethods{Method: sdk.AuthenticationMethodsOption(v)}
		}

		baseReq := sdk.NewAlterAuthenticationPolicyRequest(id)
		if len(authenticationMethodsValues) == 0 {
			unsetReq := sdk.NewAuthenticationPolicyUnsetRequest().WithAuthenticationMethods(*sdk.Bool(true))
			baseReq.WithUnset(*unsetReq)
		} else {
			setReq := sdk.NewAuthenticationPolicySetRequest().WithAuthenticationMethods(authenticationMethodsValues)
			baseReq.WithSet(*setReq)
		}

		if err := client.AuthenticationPolicies.Alter(ctx, baseReq); err != nil {
			return diag.FromErr(err)
		}
	}

	// change to mfa authentication methods
	if d.HasChange("mfa_authentication_methods") {
		mfaAuthenticationMethods := expandStringList(d.Get("mfa_authentication_methods").(*schema.Set).List())
		mfaAuthenticationMethodsValues := make([]sdk.MfaAuthenticationMethods, len(mfaAuthenticationMethods))
		for i, v := range mfaAuthenticationMethods {
			mfaAuthenticationMethodsValues[i] = sdk.MfaAuthenticationMethods{Method: sdk.MfaAuthenticationMethodsOption(v)}
		}

		baseReq := sdk.NewAlterAuthenticationPolicyRequest(id)
		if len(mfaAuthenticationMethodsValues) == 0 {
			unsetReq := sdk.NewAuthenticationPolicyUnsetRequest().WithMfaAuthenticationMethods(*sdk.Bool(true))
			baseReq.WithUnset(*unsetReq)
		} else {
			setReq := sdk.NewAuthenticationPolicySetRequest().WithMfaAuthenticationMethods(mfaAuthenticationMethodsValues)
			baseReq.WithSet(*setReq)
		}

		if err := client.AuthenticationPolicies.Alter(ctx, baseReq); err != nil {
			return diag.FromErr(err)
		}
	}

	// change to mfa enrollment
	if d.HasChange("mfa_enrollment") {
		mfaEnrollment := d.Get("mfa_enrollment").(string)

		baseReq := sdk.NewAlterAuthenticationPolicyRequest(id)
		if len(mfaEnrollment) == 0 {
			unsetReq := sdk.NewAuthenticationPolicyUnsetRequest().WithMfaEnrollment(*sdk.Bool(true))
			baseReq.WithUnset(*unsetReq)
		} else {
			setReq := sdk.NewAuthenticationPolicySetRequest().WithMfaEnrollment(sdk.MfaEnrollmentOption(mfaEnrollment))
			baseReq.WithSet(*setReq)
		}

		if err := client.AuthenticationPolicies.Alter(ctx, baseReq); err != nil {
			return diag.FromErr(err)
		}
	}

	// change to client types
	if d.HasChange("client_types") {
		clientTypes := expandStringList(d.Get("client_types").(*schema.Set).List())
		clientTypesValues := make([]sdk.ClientTypes, len(clientTypes))
		for i, v := range clientTypes {
			clientTypesValues[i] = sdk.ClientTypes{ClientType: sdk.ClientTypesOption(v)}
		}

		baseReq := sdk.NewAlterAuthenticationPolicyRequest(id)
		if len(clientTypesValues) == 0 {
			unsetReq := sdk.NewAuthenticationPolicyUnsetRequest().WithClientTypes(*sdk.Bool(true))
			baseReq.WithUnset(*unsetReq)
		} else {
			setReq := sdk.NewAuthenticationPolicySetRequest().WithClientTypes(clientTypesValues)
			baseReq.WithSet(*setReq)
		}

		if err := client.AuthenticationPolicies.Alter(ctx, baseReq); err != nil {
			return diag.FromErr(err)
		}
	}

	// change to security integrations
	if d.HasChange("security_integrations") {
		securityIntegrations := expandStringList(d.Get("security_integrations").(*schema.Set).List())
		securityIntegrationsValues := make([]sdk.SecurityIntegrationsOption, len(securityIntegrations))
		for i, v := range securityIntegrations {
			securityIntegrationsValues[i] = sdk.SecurityIntegrationsOption{Name: sdk.NewAccountObjectIdentifier(v)}
		}

		baseReq := sdk.NewAlterAuthenticationPolicyRequest(id)
		if len(securityIntegrationsValues) == 0 {
			unsetReq := sdk.NewAuthenticationPolicyUnsetRequest().WithSecurityIntegrations(*sdk.Bool(true))
			baseReq.WithUnset(*unsetReq)
		} else {
			setReq := sdk.NewAuthenticationPolicySetRequest().WithSecurityIntegrations(securityIntegrationsValues)
			baseReq.WithSet(*setReq)
		}

		if err := client.AuthenticationPolicies.Alter(ctx, baseReq); err != nil {
			return diag.FromErr(err)
		}
	}

	// change to comment
	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		baseReq := sdk.NewAlterAuthenticationPolicyRequest(id)
		if len(comment) == 0 {
			unsetReq := sdk.NewAuthenticationPolicyUnsetRequest().WithComment(*sdk.Bool(true))
			baseReq.WithUnset(*unsetReq)
		} else {
			setReq := sdk.NewAuthenticationPolicySetRequest().WithComment(comment)
			baseReq.WithSet(*setReq)
		}

		if err := client.AuthenticationPolicies.Alter(ctx, baseReq); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextAuthenticationPolicy(ctx, d, meta)
}

func DeleteContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Id()
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(name).(sdk.SchemaObjectIdentifier)

	if err := client.AuthenticationPolicies.Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(id).WithIfExists(*sdk.Bool(true))); err != nil {
		diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
