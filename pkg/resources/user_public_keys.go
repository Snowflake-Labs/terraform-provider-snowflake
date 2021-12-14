package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/go-misc/sets"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userPublicKeyProperties = []string{
	"rsa_public_key",
	"rsa_public_key_2",
}

// sanitize input to supress diffs, etc
func publicKeyStateFunc(v interface{}) string {
	value := v.(string)
	value = strings.TrimSuffix(value, "\n")
	return value
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
		StateFunc:   publicKeyStateFunc,
	},
	"rsa_public_key_2": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s second RSA public key; used to rotate the public and Public keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.",
		StateFunc:   publicKeyStateFunc,
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
	stmt := snowflake.User(name).Describe()
	_, err := snowflake.Query(db, stmt)
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
	// we can't really read the public keys back from Snowflake so assume they haven't changed
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
	propsToUnset := sets.NewStringSet()

	for _, prop := range userPublicKeyProperties {
		// if key hasn't changed, continue
		if !d.HasChange(prop) {
			continue
		}
		// if it has changed then we should do something about it
		publicKey, publicKeyOK := d.GetOk(prop)
		if publicKeyOK { // if set, then we should update the value
			propsToSet[prop] = publicKey.(string)
		} else { // if now unset, we should unset the key from the user
			propsToUnset.Add(publicKey.(string))
		}
	}

	// set the keys we decided should be set
	for prop, value := range propsToSet {
		err := updateUserPublicKeys(db, name, prop, value)
		if err != nil {
			return err
		}
	}

	// unset the keys we decided should be unset
	for _, prop := range propsToUnset.List() {
		err := unsetUserPublicKeys(db, name, prop)
		if err != nil {
			return err
		}
	}
	// re-sync
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
