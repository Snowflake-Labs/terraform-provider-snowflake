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
	"password_policy_database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Database name where the password policy is stored",
	},
	"password_policy_schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Schema name where the password policy is stored",
	},
	"password_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Non-qualified name of the password policy",
	},
}

func UserPasswordPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the password policy to use for a certain user.",

		Create: CreateUserPasswordPolicyAttachment,
		Read:   ReadUserPasswordPolicyAttachment,
		Delete: DeleteUserPasswordPolicyAttachment,

		Schema: userPasswordPolicyAttachmentSchema,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), helpers.IDDelimiter)
				if len(parts) != 4 {
					return nil, fmt.Errorf("id should be in the format 'database|schema|password_policy|user_name', but I got '%s'", d.Id())
				}
				passwordPolicyDatabase := sdk.NewAccountIdentifierFromFullyQualifiedName(parts[0])
				passwordPolicySchema := sdk.NewAccountIdentifierFromFullyQualifiedName(parts[1])
				passwordPolicyName := sdk.NewAccountIdentifierFromFullyQualifiedName(parts[2])
				userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[3])
				if err := d.Set("password_policy_database", passwordPolicyDatabase.Name()); err != nil {
					return nil, err
				}
				if err := d.Set("password_policy_schema", passwordPolicySchema.Name()); err != nil {
					return nil, err
				}
				if err := d.Set("password_policy_name", passwordPolicyName.Name()); err != nil {
					return nil, err
				}
				if err := d.Set("user_name", userName.Name()); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func CreateUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))
	passwordPolicy := sdk.NewSchemaObjectIdentifier(
		d.Get("password_policy_database").(string),
		d.Get("password_policy_schema").(string),
		d.Get("password_policy_name").(string),
	)

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			PasswordPolicy: passwordPolicy,
		},
	})
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf(`%s|%s`, helpers.EncodeSnowflakeID(passwordPolicy), helpers.EncodeSnowflakeID(userName)))

	return ReadUserPasswordPolicyAttachment(d, meta)
}

func ReadUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))
	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, &sdk.GetForEntityPolicyReferenceRequest{
		// Note: I cannot insert both single and double quotes in the SDK, so for now I need to do this
		RefEntityName:   sdk.String(userName.FullyQualifiedName()),
		RefEntityDomain: sdk.String("user"),
	})
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
	if err := d.Set("password_policy_database", sdk.NewAccountIdentifierFromFullyQualifiedName(policyReferences[0].PolicyDb).Name()); err != nil {
		return err
	}
	if err := d.Set("password_policy_schema", sdk.NewAccountIdentifierFromFullyQualifiedName(policyReferences[0].PolicySchema).Name()); err != nil {
		return err
	}
	if err := d.Set("password_policy_name", sdk.NewAccountIdentifierFromFullyQualifiedName(policyReferences[0].PolicyName).Name()); err != nil {
		return err
	}
	if err := d.Set("user_name", helpers.EncodeSnowflakeID(userName)); err != nil {
		return err
	}
	return err
}

// DeleteAccountPasswordPolicyAttachment implements schema.DeleteFunc.
func DeleteUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
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

	return nil
}
