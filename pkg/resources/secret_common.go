package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secretCommonSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
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
		Description: "Outputs the result of `SHOW SECRET` for the given secret.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSecretSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

type commonSecretCreate struct {
	name     string
	database string
	schema   string
	comment  *string
}

func handleSecretCreate(d *schema.ResourceData) commonSecretCreate {
	create := commonSecretCreate{
		name:     d.Get("name").(string),
		database: d.Get("database").(string),
		schema:   d.Get("schema").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		create.comment = sdk.Pointer(v.(string))
	}

	return create
}

func handleSecretRead(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, secret *sdk.Secret) error {
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return err
	}
	if err := d.Set("name", id.Name()); err != nil {
		return err
	}
	if err := d.Set("database", secret.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema", secret.SchemaName); err != nil {
		return err
	}
	if err := d.Set("comment", secret.Comment); err != nil {
		return err
	}
	return nil
}

type commonSecretSet struct {
	comment *string
}

type commonSecretUnset struct {
	comment *bool
}

func handleSecretUpdate(d *schema.ResourceData) (commonSecretSet, commonSecretUnset) {
	set, unset := commonSecretSet{}, commonSecretUnset{}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.comment = sdk.Pointer(v.(string))
		} else {
			unset.comment = sdk.Pointer(true)
		}
	}

	return set, unset
}
