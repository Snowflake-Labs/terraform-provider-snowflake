package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

const (
	stageIDDelimiter = '|'
)

var stageSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the stage; must be unique for the database and schema in which the stage is created.",
		ForceNew:    true,
	},
	"database": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the stage.",
		ForceNew:    true,
	},
	"schema": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the stage.",
		ForceNew:    true,
	},
	"url": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the URL for the stage.",
	},
	"credentials": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the credentials for the stage.",
	},
	"storage_integration": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the name of the storage integration used to delegate authentication responsibility for external cloud storage to a Snowflake identity and access management (IAM) entity.",
	},
	"file_format": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the file format for the stage.",
	},
	"copy_options": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the copy options for the stage.",
	},
	"encryption": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the encryption settings for the stage.",
	},
	"comment": &schema.Schema{
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
		Exists: StageExists,

		Schema: stageSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateStage implements schema.CreateFunc
func CreateStage(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	database := data.Get("database").(string)
	schema := data.Get("schema").(string)

	builder := snowflake.Stage(name, database, schema)

	// Set optionals
	if v, ok := data.GetOk("url"); ok {
		builder.WithURL(v.(string))
	}

	if v, ok := data.GetOk("credentials"); ok {
		builder.WithCredentials(v.(string))
	}

	if v, ok := data.GetOk("storage_integration"); ok {
		builder.WithStorageIntegration(v.(string))
	}

	if v, ok := data.GetOk("file_format"); ok {
		builder.WithFileFormat(v.(string))
	}

	if v, ok := data.GetOk("copy_options"); ok {
		builder.WithCopyOptions(v.(string))
	}

	if v, ok := data.GetOk("encryption"); ok {
		builder.WithEncryption(v.(string))
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	q := builder.Create()

	err := DBExec(db, q)
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
	data.SetId(dataIDInput)

	return ReadStage(data, meta)
}

// ReadStage implements schema.ReadFunc
// credentials and encryption are omitted, they cannot be read via SHOW or DESCRIBE
func ReadStage(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stageID, err := stageIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := stageID.DatabaseName
	schema := stageID.SchemaName
	stage := stageID.StageName

	q := snowflake.Stage(stage, dbName, schema).Describe()
	stageDesc, err := descStage(db, q)
	if err != nil {
		return err
	}

	sq := snowflake.Stage(stage, dbName, schema).Show()
	stageShow, err := showStage(db, sq)
	if err != nil {
		return err
	}

	err = data.Set("name", stage)
	if err != nil {
		return err
	}

	err = data.Set("database", dbName)
	if err != nil {
		return err
	}

	err = data.Set("schema", schema)
	if err != nil {
		return err
	}

	err = data.Set("url", stageDesc.url)
	if err != nil {
		return err
	}

	err = data.Set("file_format", stageDesc.fileFormat)
	if err != nil {
		return err
	}

	err = data.Set("copy_options", stageDesc.copyOptions)
	if err != nil {
		return err
	}

	err = data.Set("storage_integration", stageShow.storageIntegration)
	if err != nil {
		return err
	}

	err = data.Set("comment", stageShow.comment)
	if err != nil {
		return err
	}

	err = data.Set("aws_external_id", stageDesc.awsExternalID)
	if err != nil {
		return err
	}

	err = data.Set("snowflake_iam_user", stageDesc.snowflakeIamUser)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStage implements schema.UpdateFunc
func UpdateStage(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	stageID, err := stageIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := stageID.DatabaseName
	schema := stageID.SchemaName
	stage := stageID.StageName

	builder := snowflake.Stage(stage, dbName, schema)

	db := meta.(*sql.DB)
	if data.HasChange("url") {
		_, url := data.GetChange("url")
		q := builder.ChangeURL(url.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage url on %v", data.Id())
		}

		data.SetPartial("url")
	}

	if data.HasChange("credentials") {
		_, credentials := data.GetChange("credentials")
		q := builder.ChangeCredentials(credentials.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage credentials on %v", data.Id())
		}

		data.SetPartial("credentials")
	}

	if data.HasChange("storage_integration") {
		_, si := data.GetChange("storage_integration")
		q := builder.ChangeStorageIntegration(si.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage storage integration on %v", data.Id())
		}

		data.SetPartial("storage_integration")
	}

	if data.HasChange("encryption") {
		_, encryption := data.GetChange("encryption")
		q := builder.ChangeEncryption(encryption.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage encryption on %v", data.Id())
		}

		data.SetPartial("encryption")
	}
	if data.HasChange("file_format") {
		_, fileFormat := data.GetChange("file_format")
		q := builder.ChangeFileFormat(fileFormat.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage file formaat on %v", data.Id())
		}

		data.SetPartial("file_format")
	}
	if data.HasChange("copy_options") {
		_, copyOptions := data.GetChange("copy_options")
		q := builder.ChangeCopyOptions(copyOptions.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage copy options on %v", data.Id())
		}

		data.SetPartial("copy_options")
	}
	if data.HasChange("comment") {
		_, comment := data.GetChange("comment")
		q := builder.ChangeComment(comment.(string))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating stage comment on %v", data.Id())
		}

		data.SetPartial("comment")
	}

	return ReadStage(data, meta)
}

// DeleteStage implements schema.DeleteFunc
func DeleteStage(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stageID, err := stageIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := stageID.DatabaseName
	schema := stageID.SchemaName
	stage := stageID.StageName

	q := snowflake.Stage(stage, dbName, schema).Drop()

	err = DBExec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting stage %v", data.Id())
	}

	data.SetId("")

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

type showStageResult struct {
	createdOn           *string
	name                *string
	databaseName        *string
	schemaName          *string
	url                 *string
	hasCredentials      *string
	hasEncryptionKey    *string
	owner               *string
	comment             *string
	region              *string
	stageType           *string
	cloud               *string
	notificationChannel *string
	storageIntegration  *string
}

func showStage(db *sql.DB, query string) (showStageResult, error) {
	var r showStageResult
	row := db.QueryRow(query)
	err := row.Scan(
		&r.createdOn,
		&r.name,
		&r.databaseName,
		&r.schemaName,
		&r.url,
		&r.hasCredentials,
		&r.hasEncryptionKey,
		&r.owner,
		&r.comment,
		&r.region,
		&r.stageType,
		&r.cloud,
		&r.notificationChannel,
		&r.storageIntegration,
	)
	if err != nil {
		return r, err
	}

	return r, nil
}

type descStageResult struct {
	url              string
	awsExternalID    string
	snowflakeIamUser string
	fileFormat       string
	copyOptions      string
}

func descStage(db *sql.DB, query string) (descStageResult, error) {
	var r descStageResult
	var ff []string
	var co []string
	rows, err := db.Query(query)
	if err != nil {
		return r, err
	}
	defer rows.Close()
	for rows.Next() {
		var parentProperty string
		var property string
		var propertyType string
		var propertyValue string
		var propertyDefault string
		if err := rows.Scan(&parentProperty, &property, &propertyType, &propertyValue, &propertyDefault); err != nil {
			return r, err
		}

		switch property {
		case "URL":
			r.url = strings.Trim(propertyValue, "[\"]")
		case "AWS_EXTERNAL_ID":
			r.awsExternalID = propertyValue
		case "SNOWFLAKE_IAM_USER":
			r.snowflakeIamUser = propertyValue
		}

		switch parentProperty {
		case "STAGE_FILE_FORMAT":
			if propertyValue != propertyDefault {
				ff = append(ff, fmt.Sprintf("%s = %s", property, propertyValue))
			}
		case "STAGE_COPY_OPTIONS":
			if propertyValue != propertyDefault {
				co = append(co, fmt.Sprintf("%s = %s", property, propertyValue))
			}
		}
	}

	r.fileFormat = strings.Join(ff, " ")
	r.copyOptions = strings.Join(co, " ")
	return r, nil
}
