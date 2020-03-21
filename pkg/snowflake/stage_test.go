package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStageCreate(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.QualifiedName(), `"test_db"."test_schema"."test_stage"`)

	a.Equal(s.Create(), `CREATE STAGE "test_db"."test_schema"."test_stage"`)

	s.WithCredentials("aws_role='arn:aws:iam::001234567890:role/mysnowflakerole'")
	a.Equal(s.Create(), `CREATE STAGE "test_db"."test_schema"."test_stage" CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole')`)

	s.WithEncryption("type='AWS_SSE_KMS' kms_key_id = 'aws/key'")
	a.Equal(s.Create(), `CREATE STAGE "test_db"."test_schema"."test_stage" CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key')`)

	s.WithURL("s3://load/encrypted_files/")
	a.Equal(s.Create(), `CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key')`)

	s.WithFileFormat("format_name=my_csv_format")
	a.Equal(s.Create(), `CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format)`)

	s.WithCopyOptions("on_error='skip_file'")
	a.Equal(s.Create(), `CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format) COPY_OPTIONS = (on_error='skip_file')`)

	s.WithComment("Yeehaw")
	a.Equal(s.Create(), `CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format) COPY_OPTIONS = (on_error='skip_file') COMMENT = 'Yeehaw'`)

	s.WithStorageIntegration("MY_INTEGRATION")
	a.Equal(s.Create(), `CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') STORAGE_INTEGRATION = MY_INTEGRATION ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format) COPY_OPTIONS = (on_error='skip_file') COMMENT = 'Yeehaw'`)
}

func TestStageRename(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.Rename("test_stage2"), `ALTER STAGE "test_db"."test_schema"."test_stage" RENAME TO "test_stage2"`)
}

func TestStageChangeComment(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.ChangeComment("worst stage ever"), `ALTER STAGE "test_db"."test_schema"."test_stage" SET COMMENT = 'worst stage ever'`)
}

func TestStageChangeURL(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.ChangeURL("s3://load/test/"), `ALTER STAGE "test_db"."test_schema"."test_stage" SET URL = 's3://load/test/'`)
}

func TestStageChangeFileFormat(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.ChangeFileFormat("format_name=my_csv_format"), `ALTER STAGE "test_db"."test_schema"."test_stage" SET FILE_FORMAT = (format_name=my_csv_format)`)
}

func TestStageChangeEncryption(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.ChangeEncryption("type='AWS_SSE_KMS' kms_key_id = 'aws/key'"), `ALTER STAGE "test_db"."test_schema"."test_stage" SET ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key')`)
}

func TestStageChangeCredentials(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.ChangeCredentials("aws_role='arn:aws:iam::001234567890:role/mysnowflakerole'"), `ALTER STAGE "test_db"."test_schema"."test_stage" SET CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole')`)
}

func TestStageChangeStorageIntegration(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.ChangeStorageIntegration("MY_INTEGRATION"), `ALTER STAGE "test_db"."test_schema"."test_stage" SET STORAGE_INTEGRATION = MY_INTEGRATION`)
}

func TestStageChangeCopyOptions(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.ChangeCopyOptions("on_error='skip_file'"), `ALTER STAGE "test_db"."test_schema"."test_stage" SET COPY_OPTIONS = (on_error='skip_file')`)
}

func TestStageDrop(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.Drop(), `DROP STAGE "test_db"."test_schema"."test_stage"`)
}

func TestStageUndrop(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.Undrop(), `UNDROP STAGE "test_db"."test_schema"."test_stage"`)
}

func TestStageDescribe(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.Describe(), `DESCRIBE STAGE "test_db"."test_schema"."test_stage"`)
}

func TestStageShow(t *testing.T) {
	a := assert.New(t)
	s := Stage("test_stage", "test_db", "test_schema")
	a.Equal(s.Show(), `SHOW STAGES LIKE 'test_stage' IN DATABASE "test_db"`)
}
