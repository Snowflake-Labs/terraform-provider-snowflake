package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var authenticationPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToAuthenticationMethodsOption),
			DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToAuthenticationMethodsOption)),
		},
		Optional:    true,
		Description: fmt.Sprintf("A list of authentication methods that are allowed during login. This parameter accepts one or more of the following values: %s", possibleValuesListed(sdk.AllAuthenticationMethods)),
	},
	"mfa_authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToMfaAuthenticationMethodsOption),
			DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToMfaAuthenticationMethodsOption)),
		},
		Optional:    true,
		Description: fmt.Sprintf("A list of authentication methods that enforce multi-factor authentication (MFA) during login. Authentication methods not listed in this parameter do not prompt for multi-factor authentication. Allowed values are %s.", possibleValuesListed(sdk.AllMfaAuthenticationMethods)),
	},
	"mfa_enrollment": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Determines whether a user must enroll in multi-factor authentication. Allowed values are REQUIRED and OPTIONAL. When REQUIRED is specified, Enforces users to enroll in MFA. If this value is used, then the CLIENT_TYPES parameter must include SNOWFLAKE_UI, because Snowsight is the only place users can enroll in multi-factor authentication (MFA).",
		ValidateDiagFunc: sdkValidation(sdk.ToMfaEnrollmentOption),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToMfaEnrollmentOption)),
		Default:          "OPTIONAL",
	},
	"client_types": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToClientTypesOption),
			DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToClientTypesOption)),
		},
		Optional:    true,
		Description: fmt.Sprintf("A list of clients that can authenticate with Snowflake. If a client tries to connect, and the client is not one of the valid CLIENT_TYPES, then the login attempt fails. Allowed values are %s. The CLIENT_TYPES property of an authentication policy is a best effort method to block user logins based on specific clients. It should not be used as the sole control to establish a security boundary.", possibleValuesListed(sdk.AllClientTypes)),
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
		Description:   "Resource used to manage authentication policy objects. For more information, check [authentication policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy).",

		Schema: authenticationPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportAuthenticationPolicy,
		},
	}
}

func ImportAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting authentication policy import")
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	authenticationPolicy, err := client.AuthenticationPolicies.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("name", authenticationPolicy.Name); err != nil {
		return nil, err
	}
	if err = d.Set("database", authenticationPolicy.DatabaseName); err != nil {
		return nil, err
	}
	if err = d.Set("schema", authenticationPolicy.SchemaName); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	req := sdk.NewCreateAuthenticationPolicyRequest(id)

	// Set optionals
	if v, ok := d.GetOk("authentication_methods"); ok {
		authenticationMethodsRawList := expandStringList(v.(*schema.Set).List())
		authenticationMethods := make([]sdk.AuthenticationMethods, len(authenticationMethodsRawList))
		for i, v := range authenticationMethodsRawList {
			authenticationMethods[i] = sdk.AuthenticationMethods{Method: sdk.AuthenticationMethodsOption(v)}
		}
		req.WithAuthenticationMethods(authenticationMethods)
	}

	if v, ok := d.GetOk("mfa_authentication_methods"); ok {
		mfaAuthenticationMethodsRawList := expandStringList(v.(*schema.Set).List())
		mfaAuthenticationMethods := make([]sdk.MfaAuthenticationMethods, len(mfaAuthenticationMethodsRawList))
		for i, v := range mfaAuthenticationMethodsRawList {
			mfaAuthenticationMethods[i] = sdk.MfaAuthenticationMethods{Method: sdk.MfaAuthenticationMethodsOption(v)}
		}
		req.WithMfaAuthenticationMethods(mfaAuthenticationMethods)
	}

	if v, ok := d.GetOk("mfa_enrollment"); ok {
		req = req.WithMfaEnrollment(sdk.MfaEnrollmentOption(v.(string)))
	}

	if v, ok := d.GetOk("client_types"); ok {
		clientTypesRawList := expandStringList(v.(*schema.Set).List())
		clientTypes := make([]sdk.ClientTypes, len(clientTypesRawList))
		for i, v := range clientTypesRawList {
			clientTypes[i] = sdk.ClientTypes{ClientType: sdk.ClientTypesOption(v)}
		}
		req.WithClientTypes(clientTypes)
	}

	if v, ok := d.GetOk("security_integrations"); ok {
		securityIntegrationsRawList := expandStringList(v.(*schema.Set).List())
		securityIntegrations := make([]sdk.SecurityIntegrationsOption, len(securityIntegrationsRawList))
		for i, v := range securityIntegrationsRawList {
			securityIntegrations[i] = sdk.SecurityIntegrationsOption{Name: sdk.NewAccountObjectIdentifier(v)}
		}
		req.WithSecurityIntegrations(securityIntegrations)
	}

	if v, ok := d.GetOk("comment"); ok {
		req = req.WithComment(v.(string))
	}

	client := meta.(*provider.Context).Client
	if err := client.AuthenticationPolicies.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

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
	if err != nil {
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
		return diag.FromErr(err)
	}

	authenticationPolicyDescriptions, err := client.AuthenticationPolicies.Describe(ctx, id)
	if err != nil {
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
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewAuthenticationPolicySetRequest(), sdk.NewAuthenticationPolicyUnsetRequest()

	// change to name
	if d.HasChange("name") {
		newId, err := sdk.ParseSchemaObjectIdentifier(d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// change to authentication methods
	if d.HasChange("authentication_methods") {
		if v, ok := d.GetOk("authentication_methods"); ok {
			authenticationMethods := expandStringList(v.(*schema.Set).List())
			authenticationMethodsValues := make([]sdk.AuthenticationMethods, len(authenticationMethods))
			for i, v := range authenticationMethods {
				authenticationMethodsValues[i] = sdk.AuthenticationMethods{Method: sdk.AuthenticationMethodsOption(v)}
			}
			if len(authenticationMethodsValues) == 0 {
				unset.WithAuthenticationMethods(true)
			} else {
				set.AuthenticationMethods = authenticationMethodsValues
			}
		}
	}

	// change to mfa authentication methods
	if d.HasChange("mfa_authentication_methods") {
		if v, ok := d.GetOk("mfa_authentication_methods"); ok {
			mfaAuthenticationMethods := expandStringList(v.(*schema.Set).List())
			mfaAuthenticationMethodsValues := make([]sdk.MfaAuthenticationMethods, len(mfaAuthenticationMethods))
			for i, v := range mfaAuthenticationMethods {
				mfaAuthenticationMethodsValues[i] = sdk.MfaAuthenticationMethods{Method: sdk.MfaAuthenticationMethodsOption(v)}
			}

			if len(mfaAuthenticationMethodsValues) == 0 {
				unset.WithMfaAuthenticationMethods(true)
			} else {
				set.MfaAuthenticationMethods = mfaAuthenticationMethodsValues
			}
		}
	}

	// change to mfa enrollment
	if d.HasChange("mfa_enrollment") {
		if v, ok := d.GetOk("mfa_enrollment"); ok {
			set.MfaEnrollment = sdk.Pointer(sdk.MfaEnrollmentOption(v.(string)))
		} else {
			unset.WithMfaEnrollment(true)
		}
	}

	// change to client types
	if d.HasChange("client_types") {
		if v, ok := d.GetOk("client_types"); ok {
			clientTypes := expandStringList(v.(*schema.Set).List())
			clientTypesValues := make([]sdk.ClientTypes, len(clientTypes))
			for i, v := range clientTypes {
				clientTypesValues[i] = sdk.ClientTypes{ClientType: sdk.ClientTypesOption(v)}
			}

			if len(clientTypesValues) == 0 {
				unset.WithClientTypes(true)
			} else {
				set.ClientTypes = clientTypesValues
			}
		}
	}

	// change to security integrations
	if d.HasChange("security_integrations") {
		if v, ok := d.GetOk("security_integrations"); ok {
			securityIntegrations := expandStringList(v.(*schema.Set).List())
			securityIntegrationsValues := make([]sdk.SecurityIntegrationsOption, len(securityIntegrations))
			for i, v := range securityIntegrations {
				securityIntegrationsValues[i] = sdk.SecurityIntegrationsOption{Name: sdk.NewAccountObjectIdentifier(v)}
			}

			if len(securityIntegrationsValues) == 0 {
				unset.WithSecurityIntegrations(true)
			} else {
				set.SecurityIntegrations = securityIntegrationsValues
			}
		}
	}

	// change to comment
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.Comment = sdk.String(v.(string))
		} else {
			unset.WithComment(true)
		}
	}

	if !reflect.DeepEqual(*set, *sdk.NewAuthenticationPolicySetRequest()) {
		req := sdk.NewAlterAuthenticationPolicyRequest(id).WithSet(*set)
		if err := client.AuthenticationPolicies.Alter(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewAuthenticationPolicyUnsetRequest()) {
		req := sdk.NewAlterAuthenticationPolicyRequest(id).WithUnset(*unset)
		if err := client.AuthenticationPolicies.Alter(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextAuthenticationPolicy(ctx, d, meta)
}

func DeleteContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting authentication policy",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	if err := client.AuthenticationPolicies.Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(id).WithIfExists(true)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
