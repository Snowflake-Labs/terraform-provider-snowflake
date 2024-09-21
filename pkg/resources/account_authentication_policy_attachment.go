package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountAuthenticationPolicyAttachmentSchema = map[string]*schema.Schema{
	"authentication_policy": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Qualified name (`\"db\".\"schema\".\"policy_name\"`) of the authentication policy to apply to the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

// AccountAuthenticationPolicyAttachment returns a pointer to the resource representing an account authentication policy attachment.
func AccountAuthenticationPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the authentication policy to use for the current account. To set the authentication policy of a different account, use a provider alias.",

		Create: CreateAccountAuthenticationPolicyAttachment,
		Read:   ReadAccountAuthenticationPolicyAttachment,
		Delete: DeleteAccountAuthenticationPolicyAttachment,

		Schema: accountAuthenticationPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAccountAuthenticationPolicyAttachment implements schema.CreateFunc.
func CreateAccountAuthenticationPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	authenticationPolicy, ok := sdk.NewObjectIdentifierFromFullyQualifiedName(d.Get("authentication_policy").(string)).(sdk.SchemaObjectIdentifier)
	if !ok {
		return fmt.Errorf("authentication_policy %s is not a valid authentication policy qualified name, expected format: `\"db\".\"schema\".\"policy\"`", d.Get("authentication_policy"))
	}

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Set: &sdk.AccountSet{
			AuthenticationPolicy: authenticationPolicy,
		},
	})
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(authenticationPolicy))

	return ReadAccountAuthenticationPolicyAttachment(d, meta)
}

func ReadAccountAuthenticationPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	authenticationPolicy := helpers.DecodeSnowflakeID(d.Id())
	if err := d.Set("authentication_policy", authenticationPolicy.FullyQualifiedName()); err != nil {
		return err
	}

	return nil
}

// DeleteAccountAuthenticationPolicyAttachment implements schema.DeleteFunc.
func DeleteAccountAuthenticationPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Unset: &sdk.AccountUnset{
			AuthenticationPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
