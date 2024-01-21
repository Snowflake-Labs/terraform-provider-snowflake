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
	// "password_policy": {
	// 	Type:        schema.TypeString,
	// 	Required:    true,
	// 	ForceNew:    true,
	// 	Description: "Qualified name (`\"db\".\"schema\".\"policy_name\"`) of the password policy to apply to the current account.",
	// },
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
		Description: "Specifies the password policy to use for the current account. To set the password policy of a different account, use a provider alias.",

		Create: CreateUserPasswordPolicyAttachment,
		Read:   ReadUserPasswordPolicyAttachment,
		Delete: DeleteUserPasswordPolicyAttachment,

		Schema: userPasswordPolicyAttachmentSchema,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), helpers.IDDelimiter)
				if len(parts) != 4 {
					return nil, fmt.Errorf("id should be in the format 'database|schema|password_policy|user_name|roles', but I got '%s'", d.Id())
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
				// TODO: change for a fully qualified name
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
	fmt.Printf("CREATE FUNCTION: '%s'\n", d.Get("user_name").(string))

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))
	passwordPolicyDatabase := sdk.NewAccountIdentifierFromFullyQualifiedName(d.Get("password_policy_database").(string))
	passwordPolicySchema := sdk.NewAccountIdentifierFromFullyQualifiedName(d.Get("password_policy_schema").(string))
	passwordPolicyName := sdk.NewAccountIdentifierFromFullyQualifiedName(d.Get("password_policy_name").(string))
	passwordPolicy := sdk.NewSchemaObjectIdentifier(passwordPolicyDatabase.Name(), passwordPolicySchema.Name(), passwordPolicyName.Name())

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			PasswordPolicy: passwordPolicy,
		},
	})
	if err != nil {
		return err
	}
	if err := d.Set("password_policy_database", passwordPolicyDatabase.Name()); err != nil {
		return err
	}
	if err := d.Set("password_policy_schema", passwordPolicySchema.Name()); err != nil {
		return err
	}
	if err := d.Set("password_policy_name", passwordPolicyName.Name()); err != nil {
		return err
	}
	if err := d.Set("user_name", helpers.EncodeSnowflakeID(userName)); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf(`%s|%s`, helpers.EncodeSnowflakeID(passwordPolicy), helpers.EncodeSnowflakeID(userName)))

	return ReadUserPasswordPolicyAttachment(d, meta)
}

func ReadUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	fmt.Printf("READ FUNCTION: '%s'\n", d.Get("user_name").(string))
	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, &sdk.GetForEntityPolicyReferenceRequest{
		RefEntityName:   d.Get("user_name").(string),
		RefEntityDomain: "user",
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
		// fmt.Printf("THE RESOURCE HAS BEEN DELETED '%s'\n", d.Get("user_name").(string))
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
	// TODO: not sure if this is needed
	// if err := d.Set("user_name", d.Get("user_name").(string)); err != nil {
	// 	return err
	// }
	// fmt.Printf(policyReference.FullyQualifiedName())
	fmt.Printf("END FUNCTION: '%s'\n", d.Get("user_name").(string))
	return err
}

// DeleteAccountPasswordPolicyAttachment implements schema.DeleteFunc.
func DeleteUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	fmt.Printf("DELETE FUNCTION: '%s'\n", d.Get("user_name").(string))
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
