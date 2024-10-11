package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secretCommonSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the secret, must be unique in your schema."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the secret"),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the secret."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the secret.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SECRETS` for the given secret.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSecretSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SECRET` for the given secret.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeSecretSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func handleSecretImport(d *schema.ResourceData) error {
	if _, err := ImportName[sdk.SchemaObjectIdentifier](context.Background(), d, nil); err != nil {
		return err
	}
	return nil
}

func handleSecretCreate(d *schema.ResourceData) (database, schema, name string) {
	return d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string)
}

func handleSecretRead(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, secret *sdk.Secret, secretDescription *sdk.SecretDetails) error {
	err := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("comment", secret.Comment),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecretToSchema(secret)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.SecretDescriptionToSchema(*secretDescription)}),
	)
	if err != nil {
		return err
	}
	return nil
}

func handleSecretUpdate(d *schema.ResourceData, set *sdk.SecretSetRequest, unset *sdk.SecretUnsetRequest) {
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.WithComment(v.(string))
		} else {
			unset.WithComment(true)
		}
	}
}
