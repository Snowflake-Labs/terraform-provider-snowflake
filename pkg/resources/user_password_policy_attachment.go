package resources

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userPasswordPolicyAttachmentSchema = map[string]*schema.Schema{
	"user_name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "User name of the user you want to attach the password policy to",
	},
	"password_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Fully qualified name of the password policy",
	},
}

func UserPasswordPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the password policy to use for a certain user.",
		Create:      CreateUserPasswordPolicyAttachment,
		Read:        ReadUserPasswordPolicyAttachment,
		Delete:      DeleteUserPasswordPolicyAttachment,
		Schema:      userPasswordPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateUserPasswordPolicyAttachment(d *schema.ResourceData, meta any) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))
	passwordPolicy := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Get("password_policy_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			PasswordPolicy: &passwordPolicy,
		},
	})
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(userName.FullyQualifiedName(), passwordPolicy.FullyQualifiedName()))

	return ReadUserPasswordPolicyAttachment(d, meta)
}

func ReadUserPasswordPolicyAttachment(d *schema.ResourceData, meta any) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	if len(parts) != 2 {
		return fmt.Errorf("required id format 'user_name|password_policy_name', but got: '%s'", d.Id())
	}

	// Note: there is no alphanumeric id for an attachment, so we retrieve the password policies attached to a certain user.
	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[0])
	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(userName, sdk.PolicyEntityDomainUser))
	if err != nil {
		return err
	}

	// Note: this should never happen, but just in case: so far, Snowflake only allows one Password Policy per user.
	if len(policyReferences) > 1 {
		return fmt.Errorf("internal error: multiple policy references attached to a user. This should never happen")
	}

	// Note: this means the resource has been deleted outside of Terraform.
	if len(policyReferences) == 0 {
		d.SetId("")
		return nil
	}

	if err := d.Set("user_name", userName.Name()); err != nil {
		return err
	}
	if err := d.Set(
		"password_policy_name",
		sdk.NewSchemaObjectIdentifier(
			policyReferences[0].PolicyDb,
			policyReferences[0].PolicySchema,
			policyReferences[0].PolicyName,
		).FullyQualifiedName()); err != nil {
		return err
	}

	return err
}

func DeleteUserPasswordPolicyAttachment(d *schema.ResourceData, meta any) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Unset: &sdk.UserUnset{
			PasswordPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
