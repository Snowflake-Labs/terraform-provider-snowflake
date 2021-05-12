package resources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userPublicKeyProperties = []string{
	"rsa_public_key",
	"rsa_public_key_2",
}

var userPublicKeysSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the user.",
		ForceNew:    true,
	},

	"rsa_public_key": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s RSA public key; used for key-pair authentication. Must be on 1 line without header and trailer.",
	},
	"rsa_public_key_2": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s second RSA public key; used to rotate the public and Public keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.",
	},

	// computed
	"rsa_public_key_fp": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the user’s RSA public key; used for key-pair authentication. Must be on 1 line without header and trailer.",
	},
	"rsa_public_key_2_fp": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the user’s second RSA public key; used to rotate the public and Public keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.",
	},
}

func UserPublicKeys() *schema.Resource {
	return &schema.Resource{
		Create: CreateUserPublicKeys,
		Read:   ReadUserPublicKeys,
		Update: UpdateUserPublicKeys,
		Delete: DeleteUserPublicKeys,

		Schema: userPublicKeysSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func checkUserExists(db *sql.DB, name string) (bool, error) {
	// First check if user exists
	stmt := snowflake.User(name).Show()
	row := snowflake.QueryRow(db, stmt)
	_, err := snowflake.ScanUser(row)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] user (%s) not found", name)
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func ReadUserPublicKeys(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	exists, err := checkUserExists(db, id)
	if err != nil {
		return err
	}
	// If not found, mark resource to be removed from statefile during apply or refresh
	if !exists {
		d.SetId("")
		return nil
	}

	// at this point, we know we have a user. Read keys
	var rsaKeyFP string
	var rsaKey2FP string

	stmt := fmt.Sprintf(`DESCRIBE USER "%s"`, d.Get("name").(string))
	props, err := snowflake.Query(db, stmt)
	if err != nil {
		return err
	}
	defer props.Close()

	for props.Next() {
		prop := &snowflake.DescribeUserProp{}
		err := props.StructScan(prop)
		if err != nil {
			return err
		}

		switch prop.Property {
		case "RSA_PUBLIC_KEY_FP":
			if prop.Value.Valid {
				rsaKeyFP = prop.Value.String
			}
		case "RSA_PUBLIC_KEY_2_FP":
			if prop.Value.Valid {
				rsaKey2FP = prop.Value.String
			}
		default:
			log.Printf("[DEBUG] skipping user property %s", prop.Property)
		}
	}

	d.Set("rsa_public_key_fp", rsaKeyFP)
	d.Set("rsa_public_key_2_fp", rsaKey2FP)
	return nil
}

func CreateUserPublicKeys(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	for _, prop := range userPublicKeyProperties {
		publicKey, publicKeyOK := d.GetOk(prop)
		if !publicKeyOK {
			continue
		}
		err := updateUserPublicKeys(db, name, prop, publicKey.(string))
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return ReadUserPublicKeys(d, meta)
}

func UpdateUserPublicKeys(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()

	propsToSet := map[string]string{}
	propsToUnset := []string{}

	for _, prop := range userPublicKeyProperties {
		fingerprintProp := fmt.Sprintf("%s_fp", prop)
		if !d.HasChange(fingerprintProp) {
			continue
		}
		publicKey, publicKeyOK := d.GetOk(prop)
		if publicKeyOK {
			propsToSet[prop] = publicKey.(string)
		} else {
			propsToUnset = append(propsToUnset, prop)
		}
	}

	for prop, value := range propsToSet {
		err := updateUserPublicKeys(db, name, prop, value)
		if err != nil {
			return err
		}
	}

	for _, prop := range propsToUnset {
		err := unsetUserPublicKeys(db, name, prop)
		if err != nil {
			return err
		}
	}
	return ReadUserPublicKeys(d, meta)
}

func DeleteUserPublicKeys(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Id()

	for _, prop := range userPublicKeyProperties {
		err := unsetUserPublicKeys(db, name, prop)
		if err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}

func updateUserPublicKeys(db *sql.DB, name string, prop string, value string) error {
	stmt := fmt.Sprintf(`ALTER USER "%s" SET %s = '%s'`, name, prop, value)
	return snowflake.Exec(db, stmt)
}
func unsetUserPublicKeys(db *sql.DB, name string, prop string) error {
	stmt := fmt.Sprintf(`ALTER USER "%s" UNSET %s`, name, prop)
	return snowflake.Exec(db, stmt)
}
