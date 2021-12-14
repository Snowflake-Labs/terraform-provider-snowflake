package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

const (
	stageIDDelimiter = '|'
)

var stageSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the stage; must be unique for the database and schema in which the stage is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the stage.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the stage.",
		ForceNew:    true,
	},
	"url": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the URL for the stage.",
	},
	"credentials": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the credentials for the stage.",
		Sensitive:   true,
	},
	"storage_integration": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the name of the storage integration used to delegate authentication responsibility for external cloud storage to a Snowflake identity and access management (IAM) entity.",
	},
	"file_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the file format for the stage.",
	},
	"copy_options": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the copy options for the stage.",
	},
	"encryption": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the encryption settings for the stage.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the stage.",
	},
	"aws_external_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"snowflake_iam_user": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"tag": tagReferenceSchema,
}

type stageID struct {
	DatabaseName string
	SchemaName   string
	StageName    string
}

// String() takes in a stageID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|StageName
func (si *stageID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = stageIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.StageName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strStageID := strings.TrimSpace(buf.String())
	return strStageID, nil
}

// stageIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|StageName
// and returns a stageID object
func stageIDFromString(stringID string) (*stageID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = stageIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per stage")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	stageResult := &stageID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		StageName:    lines[0][2],
	}
	return stageResult, nil
}

// Stage returns a pointer to the resource representing a stage
func Stage() *schema.Resource {
	return &schema.Resource{
		Create: CreateStage,
		Read:   ReadStage,
		Update: UpdateStage,
		Delete: DeleteStage,

		Schema: stageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateStage implements schema.CreateFunc
func CreateStage(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)

	builder := snowflake.Stage(name, database, schema)

	// Set optionals
	if v, ok := d.GetOk("url"); ok {
		builder.WithURL(v.(string))
	}

	if v, ok := d.GetOk("credentials"); ok {
		builder.WithCredentials(v.(string))
	}

	if v, ok := d.GetOk("storage_integration"); ok {
		builder.WithStorageIntegration(v.(string))
	}

	if v, ok := d.GetOk("file_format"); ok {
		builder.WithFileFormat(v.(string))
	}

	if v, ok := d.GetOk("copy_options"); ok {
		builder.WithCopyOptions(v.(string))
	}

	if v, ok := d.GetOk("encryption"); ok {
		builder.WithEncryption(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	q := builder.Create()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating stage %v", name)
	}

	stageID := &stageID{
		DatabaseName: database,
		SchemaName:   schema,
		StageName:    name,
	}
	dataIDInput, err := stageID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadStage(d, meta)
}

// ReadStage implements schema.ReadFunc
// credentials and encryption are omitted, they cannot be read via SHOW or DESCRIBE
func ReadStage(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stageID, err := stageIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := stageID.DatabaseName
	schema := stageID.SchemaName
	stage := stageID.StageName

	q := snowflake.Stage(stage, dbName, schema).Describe()
	stageDesc, err := snowflake.DescStage(db, q)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] stage (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	sq := snowflake.Stage(stage, dbName, schema).Show()
	row := snowflake.QueryRow(db, sq)

	s, err := snowflake.ScanStageShow(row)
	if err != nil {
		return err
	}

	err = d.Set("name", s.Name)
	if err != nil {
		return err
	}

	err = d.Set("database", s.DatabaseName)
	if err != nil {
		return err
	}

	err = d.Set("schema", s.SchemaName)
	if err != nil {
		return err
	}

	err = d.Set("url", stageDesc.Url)
	if err != nil {
		return err
	}

	err = d.Set("file_format", stageDesc.FileFormat)
	if err != nil {
		return err
	}

	err = d.Set("copy_options", stageDesc.CopyOptions)
	if err != nil {
		return err
	}

	err = d.Set("storage_integration", s.StorageIntegration)
	if err != nil {
		return err
	}

	err = d.Set("comment", s.Comment)
	if err != nil {
		return err
	}

	err = d.Set("aws_external_id", stageDesc.AwsExternalID)
	if err != nil {
		return err
	}

	err = d.Set("snowflake_iam_user", stageDesc.SnowflakeIamUser)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStage implements schema.UpdateFunc
func UpdateStage(d *schema.ResourceData, meta interface{}) error {
	stageID, err := stageIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := stageID.DatabaseName
	schema := stageID.SchemaName
	stage := stageID.StageName

	builder := snowflake.Stage(stage, dbName, schema)

	db := meta.(*sql.DB)
	if d.HasChange("url") {
		url := d.Get("url")
		q := builder.ChangeURL(url.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage url on %v", d.Id())
		}
	}

	if d.HasChange("credentials") {
		credentials := d.Get("credentials")
		q := builder.ChangeCredentials(credentials.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage credentials on %v", d.Id())
		}
	}

	if d.HasChange("storage_integration") {
		si := d.Get("storage_integration")
		q := builder.ChangeStorageIntegration(si.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage storage integration on %v", d.Id())
		}
	}

	if d.HasChange("encryption") {
		encryption := d.Get("encryption")
		q := builder.ChangeEncryption(encryption.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage encryption on %v", d.Id())
		}
	}
	if d.HasChange("file_format") {
		fileFormat := d.Get("file_format")
		q := builder.ChangeFileFormat(fileFormat.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage file formaat on %v", d.Id())
		}
	}
	if d.HasChange("copy_options") {
		copyOptions := d.Get("copy_options")
		q := builder.ChangeCopyOptions(copyOptions.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage copy options on %v", d.Id())
		}
	}
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage comment on %v", d.Id())
		}
	}

	handleTagChanges(db, d, builder)

	return ReadStage(d, meta)
}

// DeleteStage implements schema.DeleteFunc
func DeleteStage(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stageID, err := stageIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := stageID.DatabaseName
	schema := stageID.SchemaName
	stage := stageID.StageName

	q := snowflake.Stage(stage, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting stage %v", d.Id())
	}

	d.SetId("")

	return nil
}

// StageExists implements schema.ExistsFunc
func StageExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	stageID, err := stageIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := stageID.DatabaseName
	schema := stageID.SchemaName
	stage := stageID.StageName

	q := snowflake.Stage(stage, dbName, schema).Describe()
	rows, err := db.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
