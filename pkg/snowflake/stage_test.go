package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStageCreate(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`"test_db"."test_schema"."test_stage"`, s.QualifiedName())

	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage"`, s.Create())

	s.WithCredentials("aws_role='arn:aws:iam::001234567890:role/mysnowflakerole'")
	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage" CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole')`, s.Create())

	s.WithEncryption("type='AWS_SSE_KMS' kms_key_id = 'aws/key'")
	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage" CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key')`, s.Create())

	s.WithURL("s3://load/encrypted_files/")
	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key')`, s.Create())

	s.WithFileFormat("format_name=my_csv_format")
	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format)`, s.Create())

	s.WithCopyOptions("on_error='skip_file'")
	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format) COPY_OPTIONS = (on_error='skip_file')`, s.Create())

	s.WithDirectory("ENABLE=TRUE")
	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format) COPY_OPTIONS = (on_error='skip_file') DIRECTORY = (ENABLE=TRUE)`, s.Create())

	s.WithComment("Yee'haw")
	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format) COPY_OPTIONS = (on_error='skip_file') DIRECTORY = (ENABLE=TRUE) COMMENT = 'Yee\'haw'`, s.Create())

	s.WithStorageIntegration("MY_INTEGRATION")
	r.Equal(`CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://load/encrypted_files/' CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole') STORAGE_INTEGRATION = "MY_INTEGRATION" ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key') FILE_FORMAT = (format_name=my_csv_format) COPY_OPTIONS = (on_error='skip_file') DIRECTORY = (ENABLE=TRUE) COMMENT = 'Yee\'haw'`, s.Create())
}

func TestStageRename(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" RENAME TO "test_stage2"`, s.Rename("test_stage2"))
}

func TestStageChangeComment(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET COMMENT = 'worst stage ever'`, s.ChangeComment("worst stage ever"))
}

func TestStageChangeURL(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET URL = 's3://load/test/'`, s.ChangeURL("s3://load/test/"))
}

func TestStageChangeFileFormat(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET FILE_FORMAT = (format_name=my_csv_format)`, s.ChangeFileFormat("format_name=my_csv_format"))
}

func TestStageChangeFileFormatToEmptyList(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET FILE_FORMAT = (TYPE = parquet NULL_IF = () COMPRESSION = none)`, s.ChangeFileFormat("TYPE = parquet NULL_IF = [] COMPRESSION = none"))
}

func TestStageChangeEncryption(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET ENCRYPTION = (type='AWS_SSE_KMS' kms_key_id = 'aws/key')`, s.ChangeEncryption("type='AWS_SSE_KMS' kms_key_id = 'aws/key'"))
}

func TestStageChangeCredentials(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET CREDENTIALS = (aws_role='arn:aws:iam::001234567890:role/mysnowflakerole')`, s.ChangeCredentials("aws_role='arn:aws:iam::001234567890:role/mysnowflakerole'"))
}

func TestStageChangeStorageIntegration(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET STORAGE_INTEGRATION = "MY_INTEGRATION"`, s.ChangeStorageIntegration("MY_INTEGRATION"))
}

func TestStageChangeStorageIntegrationAndUrl(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")

	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET STORAGE_INTEGRATION = "MY_INTEGRATION" URL = 's3://load/test'`, s.ChangeStorageIntegrationAndUrl("MY_INTEGRATION", "s3://load/test"))
}

func TestStageChangeCopyOptions(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`ALTER STAGE "test_db"."test_schema"."test_stage" SET COPY_OPTIONS = (on_error='skip_file')`, s.ChangeCopyOptions("on_error='skip_file'"))
}

func TestStageDrop(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`DROP STAGE "test_db"."test_schema"."test_stage"`, s.Drop())
}

func TestStageUndrop(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`UNDROP STAGE "test_db"."test_schema"."test_stage"`, s.Undrop())
}

func TestStageDescribe(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`DESCRIBE STAGE "test_db"."test_schema"."test_stage"`, s.Describe())
}

func TestStageShow(t *testing.T) {
	r := require.New(t)
	s := NewStageBuilder("test_stage", "test_db", "test_schema")
	r.Equal(`SHOW STAGES LIKE 'test_stage' IN SCHEMA "test_db"."test_schema"`, s.Show())
}
