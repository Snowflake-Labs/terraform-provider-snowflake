package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// StageBuilder abstracts the creation of SQL queries for a Snowflake stage
type StageBuilder struct {
	name               string
	db                 string
	schema             string
	url                string
	credentials        string
	storageIntegration string
	encryption         string
	fileFormat         string
	copyOptions        string
	comment            string
}

// QualifiedName prepends the db and schema and escapes everything nicely
func (sb *StageBuilder) QualifiedName() string {
	var n strings.Builder

	n.WriteString(fmt.Sprintf(`"%v"."%v"."%v"`, sb.db, sb.schema, sb.name))

	return n.String()
}

// WithURL adds a URL to the StageBuilder
func (sb *StageBuilder) WithURL(u string) *StageBuilder {
	sb.url = u
	return sb
}

// WithCredentials adds credentials to the StageBuilder
func (sb *StageBuilder) WithCredentials(c string) *StageBuilder {
	sb.credentials = c
	return sb
}

// WithStorageIntegration adds a storage integration to the StageBuilder
func (sb *StageBuilder) WithStorageIntegration(s string) *StageBuilder {
	sb.storageIntegration = s
	return sb
}

// WithEncryption adds encryption to the StageBuilder
func (sb *StageBuilder) WithEncryption(e string) *StageBuilder {
	sb.encryption = e
	return sb
}

// WithFileFormat adds a file format to the StageBuilder
func (sb *StageBuilder) WithFileFormat(f string) *StageBuilder {
	sb.fileFormat = f
	return sb
}

// WithCopyOptions adds copy options to the StageBuilder
func (sb *StageBuilder) WithCopyOptions(c string) *StageBuilder {
	sb.copyOptions = c
	return sb
}

// WithComment adds a comment to the StageBuilder
func (sb *StageBuilder) WithComment(c string) *StageBuilder {
	sb.comment = c
	return sb
}

// Stage returns a pointer to a Builder that abstracts the DDL operations for a stage.
//
// Supported DDL operations are:
//   - CREATE STAGE
//   - ALTER STAGE
//   - DROP STAGE
//   - UNDROP STAGE
//   - DESCRIBE STAGE
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-stage.html#stage-management)
func Stage(name, db, schema string) *StageBuilder {
	return &StageBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL query that will create a new stage.
func (sb *StageBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` STAGE %v`, sb.QualifiedName()))

	if sb.url != "" {
		q.WriteString(fmt.Sprintf(` URL = '%v'`, sb.url))
	}

	if sb.credentials != "" {
		q.WriteString(fmt.Sprintf(` CREDENTIALS = (%v)`, sb.credentials))
	}

	if sb.storageIntegration != "" {
		q.WriteString(fmt.Sprintf(` STORAGE_INTEGRATION = %v`, sb.storageIntegration))
	}

	if sb.encryption != "" {
		q.WriteString(fmt.Sprintf(` ENCRYPTION = (%v)`, sb.encryption))
	}

	if sb.fileFormat != "" {
		q.WriteString(fmt.Sprintf(` FILE_FORMAT = (%v)`, sb.fileFormat))
	}

	if sb.copyOptions != "" {
		q.WriteString(fmt.Sprintf(` COPY_OPTIONS = (%v)`, sb.copyOptions))
	}

	if sb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, sb.comment))
	}

	return q.String()
}

// Rename returns the SQL query that will rename the stage.
func (sb *StageBuilder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER STAGE %v RENAME TO "%v"`, sb.QualifiedName(), newName)
}

// ChangeComment returns the SQL query that will update the comment on the stage.
func (sb *StageBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET COMMENT = '%v'`, sb.QualifiedName(), c)
}

// RemoveComment returns the SQL query that will remove the comment on the stage.
func (sb *StageBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER STAGE %v UNSET COMMENT`, sb.QualifiedName())
}

// ChangeURL returns the SQL query that will update the url on the stage.
func (sb *StageBuilder) ChangeURL(u string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET URL = '%v'`, sb.QualifiedName(), u)
}

// ChangeCredentials returns the SQL query that will update the credentials on the stage.
func (sb *StageBuilder) ChangeCredentials(c string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET CREDENTIALS = (%v)`, sb.QualifiedName(), c)
}

// ChangeStorageIntegration returns the SQL query that will update the storage integration on the stage.
func (sb *StageBuilder) ChangeStorageIntegration(s string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET STORAGE_INTEGRATION = %v`, sb.QualifiedName(), s)
}

// ChangeEncryption returns the SQL query that will update the encryption on the stage.
func (sb *StageBuilder) ChangeEncryption(e string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET ENCRYPTION = (%v)`, sb.QualifiedName(), e)
}

// ChangeFileFormat returns the SQL query that will update the file format on the stage.
func (sb *StageBuilder) ChangeFileFormat(f string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET FILE_FORMAT = (%v)`, sb.QualifiedName(), f)
}

// ChangeCopyOptions returns the SQL query that will update the copy options on the stage.
func (sb *StageBuilder) ChangeCopyOptions(c string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET COPY_OPTIONS = (%v)`, sb.QualifiedName(), c)
}

// Drop returns the SQL query that will drop a stage.
func (sb *StageBuilder) Drop() string {
	return fmt.Sprintf(`DROP STAGE %v`, sb.QualifiedName())
}

// Undrop returns the SQL query that will undrop a stage.
func (sb *StageBuilder) Undrop() string {
	return fmt.Sprintf(`UNDROP STAGE %v`, sb.QualifiedName())
}

// Describe returns the SQL query that will describe a stage.
func (sb *StageBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE STAGE %v`, sb.QualifiedName())
}

// Show returns the SQL query that will show a stage.
func (sb *StageBuilder) Show() string {
	return fmt.Sprintf(`SHOW STAGES LIKE '%v' IN DATABASE "%v"`, sb.name, sb.db)
}

type stage struct {
	Name               *string `db:"name"`
	DatabaseName       *string `db:"database_name"`
	SchemaName         *string `db:"schema_name"`
	Comment            *string `db:"comment"`
	StorageIntegration *string `db:"storage_integration"`
}

func ScanStageShow(row *sqlx.Row) (*stage, error) {
	r := &stage{}
	err := row.StructScan(r)
	return r, err
}

type descStageResult struct {
	Url              string
	AwsExternalID    string
	SnowflakeIamUser string
	FileFormat       string
	CopyOptions      string
}

type descStageRow struct {
	ParentProperty  string `db:"parent_property"`
	Property        string `db:"property"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

func DescStage(db *sql.DB, query string) (*descStageResult, error) {
	r := &descStageResult{}
	var ff []string
	var co []string
	rows, err := Query(db, query)
	if err != nil {
		return r, err
	}
	defer rows.Close()

	for rows.Next() {

		row := &descStageRow{}
		if err := rows.StructScan(row); err != nil {
			return r, err
		}

		switch row.Property {
		case "URL":
			r.Url = strings.Trim(row.PropertyValue, "[\"]")
		case "AWS_EXTERNAL_ID":
			r.AwsExternalID = row.PropertyValue
		case "SNOWFLAKE_IAM_USER":
			r.SnowflakeIamUser = row.PropertyValue
		}

		switch row.ParentProperty {
		case "STAGE_FILE_FORMAT":
			if row.PropertyValue != row.PropertyDefault {
				ff = append(ff, fmt.Sprintf("%s = %s", row.Property, row.PropertyValue))
			}
		case "STAGE_COPY_OPTIONS":
			if row.PropertyValue != row.PropertyDefault {
				co = append(co, fmt.Sprintf("%s = %s", row.Property, row.PropertyValue))
			}
		}
	}

	r.FileFormat = strings.Join(ff, " ")
	r.CopyOptions = strings.Join(co, " ")
	return r, nil
}
