package snowflake

import (
	"fmt"
	"strings"
)

func fixFileFormat(inputFileFormat string) string {
	return strings.Replace(inputFileFormat, "NULL_IF = []", "NULL_IF = ()", 1)
}

// StageBuilder abstracts the creation of SQL queries for a Snowflake stage.
type StageBuilder struct {
	name               string
	db                 string
	schema             string
	url                string
	credentials        string
	directory          string
	storageIntegration string
	encryption         string
	fileFormat         string
	copyOptions        string
	comment            string
	tags               []TagValue
}

// QualifiedName prepends the db and schema and escapes everything nicely.
func (sb *StageBuilder) QualifiedName() string {
	var n strings.Builder

	n.WriteString(fmt.Sprintf(`"%v"."%v"."%v"`, sb.db, sb.schema, sb.name))

	return n.String()
}

// WithURL adds a URL to the StageBuilder.
func (sb *StageBuilder) WithURL(u string) *StageBuilder {
	sb.url = u
	return sb
}

// WithCredentials adds credentials to the StageBuilder.
func (sb *StageBuilder) WithCredentials(c string) *StageBuilder {
	sb.credentials = c
	return sb
}

// WithStorageIntegration adds a storage integration to the StageBuilder.
func (sb *StageBuilder) WithStorageIntegration(s string) *StageBuilder {
	sb.storageIntegration = s
	return sb
}

// WithEncryption adds encryption to the StageBuilder.
func (sb *StageBuilder) WithEncryption(e string) *StageBuilder {
	sb.encryption = e
	return sb
}

// WithFileFormat adds a file format to the StageBuilder.
func (sb *StageBuilder) WithFileFormat(f string) *StageBuilder {
	sb.fileFormat = f
	return sb
}

// WithCopyOptions adds copy options to the StageBuilder.
func (sb *StageBuilder) WithCopyOptions(c string) *StageBuilder {
	sb.copyOptions = c
	return sb
}

// WithDirectory adds directory option to the StageBuilder.
func (sb *StageBuilder) WithDirectory(d string) *StageBuilder {
	sb.directory = d
	return sb
}

// WithComment adds a comment to the StageBuilder.
func (sb *StageBuilder) WithComment(c string) *StageBuilder {
	sb.comment = c
	return sb
}

// WithTags sets the tags on the ExternalTableBuilder.
func (sb *StageBuilder) WithTags(tags []TagValue) *StageBuilder {
	sb.tags = tags
	return sb
}

// AddTag returns the SQL query that will add a new tag to the view.
func (sb *StageBuilder) AddTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER STAGE %s SET TAG "%v"."%v"."%v" = "%v"`, sb.QualifiedName(), tag.Database, tag.Schema, tag.Name, tag.Value)
}

// ChangeTag returns the SQL query that will alter a tag on the view.
func (sb *StageBuilder) ChangeTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER STAGE %s SET TAG "%v"."%v"."%v" = "%v"`, sb.QualifiedName(), tag.Database, tag.Schema, tag.Name, tag.Value)
}

// UnsetTag returns the SQL query that will unset a tag on the view.
func (sb *StageBuilder) UnsetTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER STAGE %s UNSET TAG "%v"."%v"."%v"`, sb.QualifiedName(), tag.Database, tag.Schema, tag.Name)
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
func NewStageBuilder(name, db, schema string) *StageBuilder {
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
		q.WriteString(fmt.Sprintf(` STORAGE_INTEGRATION = "%v"`, sb.storageIntegration))
	}

	if sb.encryption != "" {
		q.WriteString(fmt.Sprintf(` ENCRYPTION = (%v)`, sb.encryption))
	}

	if sb.fileFormat != "" {
		q.WriteString(fmt.Sprintf(` FILE_FORMAT = (%v)`, fixFileFormat(sb.fileFormat)))
	}

	if sb.copyOptions != "" {
		q.WriteString(fmt.Sprintf(` COPY_OPTIONS = (%v)`, sb.copyOptions))
	}

	if sb.directory != "" {
		q.WriteString(fmt.Sprintf(` DIRECTORY = (%v)`, sb.directory))
	}

	if sb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(sb.comment)))
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
	return fmt.Sprintf(`ALTER STAGE %v SET STORAGE_INTEGRATION = "%v"`, sb.QualifiedName(), s)
}

func (sb *StageBuilder) ChangeStorageIntegrationAndUrl(s string, url string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET STORAGE_INTEGRATION = "%v" URL = '%v'`, sb.QualifiedName(), s, url)
}

// ChangeEncryption returns the SQL query that will update the encryption on the stage.
func (sb *StageBuilder) ChangeEncryption(e string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET ENCRYPTION = (%v)`, sb.QualifiedName(), e)
}

// ChangeFileFormat returns the SQL query that will update the file format on the stage.
func (sb *StageBuilder) ChangeFileFormat(f string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET FILE_FORMAT = (%v)`, sb.QualifiedName(), fixFileFormat(f))
}

// ChangeCopyOptions returns the SQL query that will update the copy options on the stage.
func (sb *StageBuilder) ChangeCopyOptions(c string) string {
	return fmt.Sprintf(`ALTER STAGE %v SET COPY_OPTIONS = (%v)`, sb.QualifiedName(), c)
}
