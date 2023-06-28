package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_FileFormatsCreateAndRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schema, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	t.Run("CSV", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, randomString(t))
		err := client.FileFormats.Create(ctx, id, &CreateFileFormatOptions{
			Type: FileFormatTypeCsv,
			FileFormatTypeOptions: FileFormatTypeOptions{
				CsvCompression:                &CsvCompressionBz2,
				CsvRecordDelimiter:            String("\\123"),
				CsvFieldDelimiter:             String("0x42"),
				CsvFileExtension:              String("c"),
				CsvParseHeader:                Bool(true),
				CsvSkipBlankLines:             Bool(true),
				CsvDateFormat:                 String("d"),
				CsvTimeFormat:                 String("e"),
				CsvTimestampFormat:            String("f"),
				CsvBinaryFormat:               &BinaryFormatBase64,
				CsvEscape:                     String(`\`),
				CsvEscapeUnenclosedField:      String("h"),
				CsvTrimSpace:                  Bool(true),
				CsvFieldOptionallyEnclosedBy:  String("'"),
				CsvNullIf:                     &[]NullString{{"j"}, {"k"}},
				CsvErrorOnColumnCountMismatch: Bool(true),
				CsvReplaceInvalidCharacters:   Bool(true),
				CsvEmptyFieldAsNull:           Bool(true),
				CsvSkipByteOrderMark:          Bool(true),
				CsvEncoding:                   &CsvEncodingGB18030,

				Comment: String("test comment"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, FileFormatTypeCsv, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)
		assert.Equal(t, &CsvCompressionBz2, result.Options.CsvCompression)
		assert.Equal(t, "S", *result.Options.CsvRecordDelimiter) // o123 == 83 == 'S' (ASCII)
		assert.Equal(t, "B", *result.Options.CsvFieldDelimiter)  // 0x42 == 66 == 'B' (ASCII)
		assert.Equal(t, "c", *result.Options.CsvFileExtension)
		assert.Equal(t, true, *result.Options.CsvParseHeader)
		assert.Equal(t, true, *result.Options.CsvSkipBlankLines)
		assert.Equal(t, "d", *result.Options.CsvDateFormat)
		assert.Equal(t, "e", *result.Options.CsvTimeFormat)
		assert.Equal(t, "f", *result.Options.CsvTimestampFormat)
		assert.Equal(t, &BinaryFormatBase64, result.Options.CsvBinaryFormat)
		assert.Equal(t, `\`, *result.Options.CsvEscape)
		assert.Equal(t, "h", *result.Options.CsvEscapeUnenclosedField)
		assert.Equal(t, true, *result.Options.CsvTrimSpace)
		assert.Equal(t, String("'"), result.Options.CsvFieldOptionallyEnclosedBy)
		assert.Equal(t, &[]NullString{{"j"}, {"k"}}, result.Options.CsvNullIf)
		assert.Equal(t, true, *result.Options.CsvErrorOnColumnCountMismatch)
		assert.Equal(t, true, *result.Options.CsvReplaceInvalidCharacters)
		assert.Equal(t, true, *result.Options.CsvEmptyFieldAsNull)
		assert.Equal(t, true, *result.Options.CsvSkipByteOrderMark)
		assert.Equal(t, &CsvEncodingGB18030, result.Options.CsvEncoding)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeCsv, describeResult.Type)
		assert.Equal(t, &CsvCompressionBz2, describeResult.Options.CsvCompression)
		assert.Equal(t, "S", *describeResult.Options.CsvRecordDelimiter) // o123 == 83 == 'S' (ASCII)
		assert.Equal(t, "B", *describeResult.Options.CsvFieldDelimiter)  // 0x42 == 66 == 'B' (ASCII)
		assert.Equal(t, "c", *describeResult.Options.CsvFileExtension)
		assert.Equal(t, true, *describeResult.Options.CsvParseHeader)
		assert.Equal(t, true, *describeResult.Options.CsvSkipBlankLines)
		assert.Equal(t, "d", *describeResult.Options.CsvDateFormat)
		assert.Equal(t, "e", *describeResult.Options.CsvTimeFormat)
		assert.Equal(t, "f", *describeResult.Options.CsvTimestampFormat)
		assert.Equal(t, &BinaryFormatBase64, describeResult.Options.CsvBinaryFormat)
		assert.Equal(t, `\\`, *describeResult.Options.CsvEscape) // Describe does not un-escape backslashes, but show does ....
		assert.Equal(t, "h", *describeResult.Options.CsvEscapeUnenclosedField)
		assert.Equal(t, true, *describeResult.Options.CsvTrimSpace)
		assert.Equal(t, String("'"), describeResult.Options.CsvFieldOptionallyEnclosedBy)
		assert.Equal(t, &[]NullString{{"j"}, {"k"}}, describeResult.Options.CsvNullIf)
		assert.Equal(t, true, *describeResult.Options.CsvErrorOnColumnCountMismatch)
		assert.Equal(t, true, *describeResult.Options.CsvReplaceInvalidCharacters)
		assert.Equal(t, true, *describeResult.Options.CsvEmptyFieldAsNull)
		assert.Equal(t, true, *describeResult.Options.CsvSkipByteOrderMark)
		assert.Equal(t, &CsvEncodingGB18030, describeResult.Options.CsvEncoding)
	})
	t.Run("JSON", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, randomString(t))
		err := client.FileFormats.Create(ctx, id, &CreateFileFormatOptions{
			Type: FileFormatTypeJson,
			FileFormatTypeOptions: FileFormatTypeOptions{
				JsonCompression:       &JsonCompressionBrotli,
				JsonDateFormat:        String("a"),
				JsonTimeFormat:        String("b"),
				JsonTimestampFormat:   String("c"),
				JsonBinaryFormat:      &BinaryFormatHex,
				JsonTrimSpace:         Bool(true),
				JsonNullIf:            &[]NullString{{"d"}, {"e"}},
				JsonFileExtension:     String("f"),
				JsonEnableOctal:       Bool(true),
				JsonAllowDuplicate:    Bool(true),
				JsonStripOuterArray:   Bool(true),
				JsonStripNullValues:   Bool(true),
				JsonIgnoreUtf8Errors:  Bool(true),
				JsonSkipByteOrderMark: Bool(true),

				Comment: String("test comment"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, FileFormatTypeJson, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)

		assert.Equal(t, JsonCompressionBrotli, *result.Options.JsonCompression)
		assert.Equal(t, "a", *result.Options.JsonDateFormat)
		assert.Equal(t, "b", *result.Options.JsonTimeFormat)
		assert.Equal(t, "c", *result.Options.JsonTimestampFormat)
		assert.Equal(t, BinaryFormatHex, *result.Options.JsonBinaryFormat)
		assert.Equal(t, true, *result.Options.JsonTrimSpace)
		assert.Equal(t, []NullString{{"d"}, {"e"}}, *result.Options.JsonNullIf)
		assert.Equal(t, "f", *result.Options.JsonFileExtension)
		assert.Equal(t, true, *result.Options.JsonEnableOctal)
		assert.Equal(t, true, *result.Options.JsonAllowDuplicate)
		assert.Equal(t, true, *result.Options.JsonStripOuterArray)
		assert.Equal(t, true, *result.Options.JsonStripNullValues)
		assert.Equal(t, true, *result.Options.JsonIgnoreUtf8Errors)
		assert.Equal(t, true, *result.Options.JsonSkipByteOrderMark)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeJson, describeResult.Type)
		assert.Equal(t, JsonCompressionBrotli, *describeResult.Options.JsonCompression)
		assert.Equal(t, "a", *describeResult.Options.JsonDateFormat)
		assert.Equal(t, "b", *describeResult.Options.JsonTimeFormat)
		assert.Equal(t, "c", *describeResult.Options.JsonTimestampFormat)
		assert.Equal(t, BinaryFormatHex, *describeResult.Options.JsonBinaryFormat)
		assert.Equal(t, true, *describeResult.Options.JsonTrimSpace)
		assert.Equal(t, []NullString{{"d"}, {"e"}}, *describeResult.Options.JsonNullIf)
		assert.Equal(t, "f", *describeResult.Options.JsonFileExtension)
		assert.Equal(t, true, *describeResult.Options.JsonEnableOctal)
		assert.Equal(t, true, *describeResult.Options.JsonAllowDuplicate)
		assert.Equal(t, true, *describeResult.Options.JsonStripOuterArray)
		assert.Equal(t, true, *describeResult.Options.JsonStripNullValues)
		assert.Equal(t, true, *describeResult.Options.JsonIgnoreUtf8Errors)
		assert.Equal(t, true, *describeResult.Options.JsonSkipByteOrderMark)
	})
	t.Run("AVRO", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, randomString(t))
		err := client.FileFormats.Create(ctx, id, &CreateFileFormatOptions{
			Type: FileFormatTypeAvro,
			FileFormatTypeOptions: FileFormatTypeOptions{
				AvroCompression:              &AvroCompressionGzip,
				AvroTrimSpace:                Bool(true),
				AvroReplaceInvalidCharacters: Bool(true),
				AvroNullIf:                   &[]NullString{{"a"}, {"b"}},

				Comment: String("test comment"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, FileFormatTypeAvro, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)

		assert.Equal(t, AvroCompressionGzip, *result.Options.AvroCompression)
		assert.Equal(t, true, *result.Options.AvroTrimSpace)
		assert.Equal(t, true, *result.Options.AvroReplaceInvalidCharacters)
		assert.Equal(t, []NullString{{"a"}, {"b"}}, *result.Options.AvroNullIf)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeAvro, describeResult.Type)
		assert.Equal(t, AvroCompressionGzip, *describeResult.Options.AvroCompression)
		assert.Equal(t, true, *describeResult.Options.AvroTrimSpace)
		assert.Equal(t, true, *describeResult.Options.AvroReplaceInvalidCharacters)
		assert.Equal(t, []NullString{{"a"}, {"b"}}, *describeResult.Options.AvroNullIf)
	})
	t.Run("ORC", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, randomString(t))
		err := client.FileFormats.Create(ctx, id, &CreateFileFormatOptions{
			Type: FileFormatTypeOrc,
			FileFormatTypeOptions: FileFormatTypeOptions{
				OrcTrimSpace:                Bool(true),
				OrcReplaceInvalidCharacters: Bool(true),
				OrcNullIf:                   &[]NullString{{"a"}, {"b"}},

				Comment: String("test comment"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, FileFormatTypeOrc, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)

		assert.Equal(t, true, *result.Options.OrcTrimSpace)
		assert.Equal(t, true, *result.Options.OrcReplaceInvalidCharacters)
		assert.Equal(t, []NullString{{"a"}, {"b"}}, *result.Options.OrcNullIf)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeOrc, describeResult.Type)
		assert.Equal(t, true, *describeResult.Options.OrcTrimSpace)
		assert.Equal(t, true, *describeResult.Options.OrcReplaceInvalidCharacters)
		assert.Equal(t, []NullString{{"a"}, {"b"}}, *describeResult.Options.OrcNullIf)
	})
	t.Run("PARQUET", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, randomString(t))
		err := client.FileFormats.Create(ctx, id, &CreateFileFormatOptions{
			Type: FileFormatTypeParquet,
			FileFormatTypeOptions: FileFormatTypeOptions{
				ParquetCompression:              &ParquetCompressionLzo,
				ParquetBinaryAsText:             Bool(true),
				ParquetTrimSpace:                Bool(true),
				ParquetReplaceInvalidCharacters: Bool(true),
				ParquetNullIf:                   &[]NullString{{"a"}, {"b"}},

				Comment: String("test comment"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, FileFormatTypeParquet, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)

		assert.Equal(t, ParquetCompressionLzo, *result.Options.ParquetCompression)
		assert.Equal(t, true, *result.Options.ParquetBinaryAsText)
		assert.Equal(t, true, *result.Options.ParquetTrimSpace)
		assert.Equal(t, true, *result.Options.ParquetReplaceInvalidCharacters)
		assert.Equal(t, []NullString{{"a"}, {"b"}}, *result.Options.ParquetNullIf)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeParquet, describeResult.Type)
		assert.Equal(t, ParquetCompressionLzo, *describeResult.Options.ParquetCompression)
		assert.Equal(t, true, *describeResult.Options.ParquetBinaryAsText)
		assert.Equal(t, true, *describeResult.Options.ParquetTrimSpace)
		assert.Equal(t, true, *describeResult.Options.ParquetReplaceInvalidCharacters)
		assert.Equal(t, []NullString{{"a"}, {"b"}}, *describeResult.Options.ParquetNullIf)
	})
	t.Run("XML", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, randomString(t))
		err := client.FileFormats.Create(ctx, id, &CreateFileFormatOptions{
			Type: FileFormatTypeXml,
			FileFormatTypeOptions: FileFormatTypeOptions{
				XmlCompression:          &XmlCompressionDeflate,
				XmlIgnoreUtf8Errors:     Bool(true),
				XmlPreserveSpace:        Bool(true),
				XmlStripOuterElement:    Bool(true),
				XmlDisableSnowflakeData: Bool(true),
				XmlDisableAutoConvert:   Bool(true),
				XmlSkipByteOrderMark:    Bool(true),

				Comment: String("test comment"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.FileFormats.Drop(ctx, id, nil)
			require.NoError(t, err)
		})

		result, err := client.FileFormats.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, id, result.Name)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
		assert.Equal(t, FileFormatTypeXml, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)

		assert.Equal(t, XmlCompressionDeflate, *result.Options.XmlCompression)
		assert.Equal(t, true, *result.Options.XmlIgnoreUtf8Errors)
		assert.Equal(t, true, *result.Options.XmlPreserveSpace)
		assert.Equal(t, true, *result.Options.XmlStripOuterElement)
		assert.Equal(t, true, *result.Options.XmlDisableSnowflakeData)
		assert.Equal(t, true, *result.Options.XmlDisableAutoConvert)
		assert.Equal(t, true, *result.Options.XmlSkipByteOrderMark)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeXml, describeResult.Type)
		assert.Equal(t, XmlCompressionDeflate, *describeResult.Options.XmlCompression)
		assert.Equal(t, true, *describeResult.Options.XmlIgnoreUtf8Errors)
		assert.Equal(t, true, *describeResult.Options.XmlPreserveSpace)
		assert.Equal(t, true, *describeResult.Options.XmlStripOuterElement)
		assert.Equal(t, true, *describeResult.Options.XmlDisableSnowflakeData)
		assert.Equal(t, true, *describeResult.Options.XmlDisableAutoConvert)
		assert.Equal(t, true, *describeResult.Options.XmlSkipByteOrderMark)
	})
}

func TestInt_FileFormatsAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, cleanupDatabase := createDatabase(t, client)
	t.Cleanup(cleanupDatabase)
	schemaTest, cleanupSchema := createSchema(t, client, databaseTest)
	t.Cleanup(cleanupSchema)

	t.Run("rename", func(t *testing.T) {
		fileFormat, fileFormatCleanup := createFileFormat(t, client, schemaTest.ID())
		t.Cleanup(fileFormatCleanup)
		oldId := fileFormat.ID()
		newId := NewSchemaObjectIdentifier(oldId.databaseName, oldId.schemaName, randomString(t))

		err := client.FileFormats.Alter(ctx, oldId, &AlterFileFormatOptions{
			Rename: &AlterFileFormatRenameOptions{
				NewName: newId,
			},
		})
		require.NoError(t, err)

		_, err = client.FileFormats.ShowByID(ctx, oldId)
		require.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)

		result, err := client.FileFormats.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId, result.Name)

		// Undo rename so we can clean up
		err = client.FileFormats.Alter(ctx, newId, &AlterFileFormatOptions{
			Rename: &AlterFileFormatRenameOptions{
				NewName: oldId,
			},
		})
		require.NoError(t, err)
	})

	t.Run("set", func(t *testing.T) {
		fileFormat, fileFormatCleanup := createFileFormatWithOptions(t, client, schemaTest.ID(), &CreateFileFormatOptions{
			Type: FileFormatTypeCsv,
			FileFormatTypeOptions: FileFormatTypeOptions{
				CsvCompression: &CsvCompressionAuto,
				CsvParseHeader: Bool(false),
			},
		})
		t.Cleanup(fileFormatCleanup)

		err := client.FileFormats.Alter(ctx, fileFormat.ID(), &AlterFileFormatOptions{
			Set: &FileFormatTypeOptions{
				CsvCompression: &CsvCompressionBz2,
				CsvParseHeader: Bool(true),
			},
		})
		require.NoError(t, err)

		result, err := client.FileFormats.ShowByID(ctx, fileFormat.ID())
		require.NoError(t, err)
		assert.Equal(t, CsvCompressionBz2, *result.Options.CsvCompression)
		assert.Equal(t, true, *result.Options.CsvParseHeader)
	})
}

func TestInt_FileFormatsDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, cleanupDatabase := createDatabase(t, client)
	t.Cleanup(cleanupDatabase)
	schemaTest, cleanupSchema := createSchema(t, client, databaseTest)
	t.Cleanup(cleanupSchema)
	t.Run("no options", func(t *testing.T) {
		fileFormat, _ := createFileFormat(t, client, schemaTest.ID())
		err := client.FileFormats.Drop(ctx, fileFormat.ID(), nil)
		require.NoError(t, err)

		_, err = client.FileFormats.ShowByID(ctx, fileFormat.ID())
		require.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("with IfExists", func(t *testing.T) {
		fileFormat, _ := createFileFormat(t, client, schemaTest.ID())
		err := client.FileFormats.Drop(ctx, fileFormat.ID(), &DropFileFormatOptions{
			IfExists: Bool(true),
		})
		require.NoError(t, err)

		_, err = client.FileFormats.ShowByID(ctx, fileFormat.ID())
		require.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_FileFormatsShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, cleanupDatabase := createDatabase(t, client)
	t.Cleanup(cleanupDatabase)
	schemaTest, cleanupSchema := createSchema(t, client, databaseTest)
	t.Cleanup(cleanupSchema)
	fileFormatTest, cleanupFileFormat := createFileFormat(t, client, schemaTest.ID())
	t.Cleanup(cleanupFileFormat)
	fileFormatTest2, cleanupFileFormat2 := createFileFormat(t, client, schemaTest.ID())
	t.Cleanup(cleanupFileFormat2)

	t.Run("without options", func(t *testing.T) {
		fileFormats, err := client.FileFormats.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(fileFormats))
		assert.Contains(t, fileFormats, fileFormatTest)
		assert.Contains(t, fileFormats, fileFormatTest2)
	})

	t.Run("LIKE", func(t *testing.T) {
		fileFormats, err := client.FileFormats.Show(ctx, &ShowFileFormatsOptions{
			Like: &Like{
				Pattern: String(fileFormatTest.Name.name),
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(fileFormats))
		assert.Contains(t, fileFormats, fileFormatTest)
	})

	t.Run("IN", func(t *testing.T) {
		fileFormats, err := client.FileFormats.Show(ctx, &ShowFileFormatsOptions{
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(fileFormats))
		assert.Contains(t, fileFormats, fileFormatTest)
		assert.Contains(t, fileFormats, fileFormatTest2)
	})
}
