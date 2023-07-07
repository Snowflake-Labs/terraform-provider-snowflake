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
			Type: FileFormatTypeCSV,
			FileFormatTypeOptions: FileFormatTypeOptions{
				CSVCompression:                &CSVCompressionBz2,
				CSVRecordDelimiter:            String("\\123"),
				CSVFieldDelimiter:             String("0x42"),
				CSVFileExtension:              String("c"),
				CSVParseHeader:                Bool(true),
				CSVSkipBlankLines:             Bool(true),
				CSVDateFormat:                 String("d"),
				CSVTimeFormat:                 String("e"),
				CSVTimestampFormat:            String("f"),
				CSVBinaryFormat:               &BinaryFormatBase64,
				CSVEscape:                     String(`\`),
				CSVEscapeUnenclosedField:      String("h"),
				CSVTrimSpace:                  Bool(true),
				CSVFieldOptionallyEnclosedBy:  String("'"),
				CSVNullIf:                     &[]NullString{{"j"}, {"k"}},
				CSVErrorOnColumnCountMismatch: Bool(true),
				CSVReplaceInvalidCharacters:   Bool(true),
				CSVEmptyFieldAsNull:           Bool(true),
				CSVSkipByteOrderMark:          Bool(true),
				CSVEncoding:                   &CSVEncodingGB18030,

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
		assert.Equal(t, FileFormatTypeCSV, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)
		assert.Equal(t, &CSVCompressionBz2, result.Options.CSVCompression)
		assert.Equal(t, "S", *result.Options.CSVRecordDelimiter) // o123 == 83 == 'S' (ASCII)
		assert.Equal(t, "B", *result.Options.CSVFieldDelimiter)  // 0x42 == 66 == 'B' (ASCII)
		assert.Equal(t, "c", *result.Options.CSVFileExtension)
		assert.Equal(t, true, *result.Options.CSVParseHeader)
		assert.Equal(t, true, *result.Options.CSVSkipBlankLines)
		assert.Equal(t, "d", *result.Options.CSVDateFormat)
		assert.Equal(t, "e", *result.Options.CSVTimeFormat)
		assert.Equal(t, "f", *result.Options.CSVTimestampFormat)
		assert.Equal(t, &BinaryFormatBase64, result.Options.CSVBinaryFormat)
		assert.Equal(t, `\`, *result.Options.CSVEscape)
		assert.Equal(t, "h", *result.Options.CSVEscapeUnenclosedField)
		assert.Equal(t, true, *result.Options.CSVTrimSpace)
		assert.Equal(t, String("'"), result.Options.CSVFieldOptionallyEnclosedBy)
		assert.Equal(t, &[]NullString{{"j"}, {"k"}}, result.Options.CSVNullIf)
		assert.Equal(t, true, *result.Options.CSVErrorOnColumnCountMismatch)
		assert.Equal(t, true, *result.Options.CSVReplaceInvalidCharacters)
		assert.Equal(t, true, *result.Options.CSVEmptyFieldAsNull)
		assert.Equal(t, true, *result.Options.CSVSkipByteOrderMark)
		assert.Equal(t, &CSVEncodingGB18030, result.Options.CSVEncoding)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeCSV, describeResult.Type)
		assert.Equal(t, &CSVCompressionBz2, describeResult.Options.CSVCompression)
		assert.Equal(t, "S", *describeResult.Options.CSVRecordDelimiter) // o123 == 83 == 'S' (ASCII)
		assert.Equal(t, "B", *describeResult.Options.CSVFieldDelimiter)  // 0x42 == 66 == 'B' (ASCII)
		assert.Equal(t, "c", *describeResult.Options.CSVFileExtension)
		assert.Equal(t, true, *describeResult.Options.CSVParseHeader)
		assert.Equal(t, true, *describeResult.Options.CSVSkipBlankLines)
		assert.Equal(t, "d", *describeResult.Options.CSVDateFormat)
		assert.Equal(t, "e", *describeResult.Options.CSVTimeFormat)
		assert.Equal(t, "f", *describeResult.Options.CSVTimestampFormat)
		assert.Equal(t, &BinaryFormatBase64, describeResult.Options.CSVBinaryFormat)
		assert.Equal(t, `\\`, *describeResult.Options.CSVEscape) // Describe does not un-escape backslashes, but show does ....
		assert.Equal(t, "h", *describeResult.Options.CSVEscapeUnenclosedField)
		assert.Equal(t, true, *describeResult.Options.CSVTrimSpace)
		assert.Equal(t, String("'"), describeResult.Options.CSVFieldOptionallyEnclosedBy)
		assert.Equal(t, &[]NullString{{"j"}, {"k"}}, describeResult.Options.CSVNullIf)
		assert.Equal(t, true, *describeResult.Options.CSVErrorOnColumnCountMismatch)
		assert.Equal(t, true, *describeResult.Options.CSVReplaceInvalidCharacters)
		assert.Equal(t, true, *describeResult.Options.CSVEmptyFieldAsNull)
		assert.Equal(t, true, *describeResult.Options.CSVSkipByteOrderMark)
		assert.Equal(t, &CSVEncodingGB18030, describeResult.Options.CSVEncoding)
	})
	t.Run("JSON", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, randomString(t))
		err := client.FileFormats.Create(ctx, id, &CreateFileFormatOptions{
			Type: FileFormatTypeJSON,
			FileFormatTypeOptions: FileFormatTypeOptions{
				JSONCompression:       &JSONCompressionBrotli,
				JSONDateFormat:        String("a"),
				JSONTimeFormat:        String("b"),
				JSONTimestampFormat:   String("c"),
				JSONBinaryFormat:      &BinaryFormatHex,
				JSONTrimSpace:         Bool(true),
				JSONNullIf:            &[]NullString{{"d"}, {"e"}},
				JSONFileExtension:     String("f"),
				JSONEnableOctal:       Bool(true),
				JSONAllowDuplicate:    Bool(true),
				JSONStripOuterArray:   Bool(true),
				JSONStripNullValues:   Bool(true),
				JSONIgnoreUTF8Errors:  Bool(true),
				JSONSkipByteOrderMark: Bool(true),

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
		assert.Equal(t, FileFormatTypeJSON, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)

		assert.Equal(t, JSONCompressionBrotli, *result.Options.JSONCompression)
		assert.Equal(t, "a", *result.Options.JSONDateFormat)
		assert.Equal(t, "b", *result.Options.JSONTimeFormat)
		assert.Equal(t, "c", *result.Options.JSONTimestampFormat)
		assert.Equal(t, BinaryFormatHex, *result.Options.JSONBinaryFormat)
		assert.Equal(t, true, *result.Options.JSONTrimSpace)
		assert.Equal(t, []NullString{{"d"}, {"e"}}, *result.Options.JSONNullIf)
		assert.Equal(t, "f", *result.Options.JSONFileExtension)
		assert.Equal(t, true, *result.Options.JSONEnableOctal)
		assert.Equal(t, true, *result.Options.JSONAllowDuplicate)
		assert.Equal(t, true, *result.Options.JSONStripOuterArray)
		assert.Equal(t, true, *result.Options.JSONStripNullValues)
		assert.Equal(t, true, *result.Options.JSONIgnoreUTF8Errors)
		assert.Equal(t, true, *result.Options.JSONSkipByteOrderMark)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeJSON, describeResult.Type)
		assert.Equal(t, JSONCompressionBrotli, *describeResult.Options.JSONCompression)
		assert.Equal(t, "a", *describeResult.Options.JSONDateFormat)
		assert.Equal(t, "b", *describeResult.Options.JSONTimeFormat)
		assert.Equal(t, "c", *describeResult.Options.JSONTimestampFormat)
		assert.Equal(t, BinaryFormatHex, *describeResult.Options.JSONBinaryFormat)
		assert.Equal(t, true, *describeResult.Options.JSONTrimSpace)
		assert.Equal(t, []NullString{{"d"}, {"e"}}, *describeResult.Options.JSONNullIf)
		assert.Equal(t, "f", *describeResult.Options.JSONFileExtension)
		assert.Equal(t, true, *describeResult.Options.JSONEnableOctal)
		assert.Equal(t, true, *describeResult.Options.JSONAllowDuplicate)
		assert.Equal(t, true, *describeResult.Options.JSONStripOuterArray)
		assert.Equal(t, true, *describeResult.Options.JSONStripNullValues)
		assert.Equal(t, true, *describeResult.Options.JSONIgnoreUTF8Errors)
		assert.Equal(t, true, *describeResult.Options.JSONSkipByteOrderMark)
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
			Type: FileFormatTypeORC,
			FileFormatTypeOptions: FileFormatTypeOptions{
				ORCTrimSpace:                Bool(true),
				ORCReplaceInvalidCharacters: Bool(true),
				ORCNullIf:                   &[]NullString{{"a"}, {"b"}},

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
		assert.Equal(t, FileFormatTypeORC, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)

		assert.Equal(t, true, *result.Options.ORCTrimSpace)
		assert.Equal(t, true, *result.Options.ORCReplaceInvalidCharacters)
		assert.Equal(t, []NullString{{"a"}, {"b"}}, *result.Options.ORCNullIf)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeORC, describeResult.Type)
		assert.Equal(t, true, *describeResult.Options.ORCTrimSpace)
		assert.Equal(t, true, *describeResult.Options.ORCReplaceInvalidCharacters)
		assert.Equal(t, []NullString{{"a"}, {"b"}}, *describeResult.Options.ORCNullIf)
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
			Type: FileFormatTypeXML,
			FileFormatTypeOptions: FileFormatTypeOptions{
				XMLCompression:          &XMLCompressionDeflate,
				XMLIgnoreUTF8Errors:     Bool(true),
				XMLPreserveSpace:        Bool(true),
				XMLStripOuterElement:    Bool(true),
				XMLDisableSnowflakeData: Bool(true),
				XMLDisableAutoConvert:   Bool(true),
				XMLSkipByteOrderMark:    Bool(true),

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
		assert.Equal(t, FileFormatTypeXML, result.Type)
		assert.Equal(t, client.config.Role, result.Owner)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, "", result.OwnerRoleType)

		assert.Equal(t, XMLCompressionDeflate, *result.Options.XMLCompression)
		assert.Equal(t, true, *result.Options.XMLIgnoreUTF8Errors)
		assert.Equal(t, true, *result.Options.XMLPreserveSpace)
		assert.Equal(t, true, *result.Options.XMLStripOuterElement)
		assert.Equal(t, true, *result.Options.XMLDisableSnowflakeData)
		assert.Equal(t, true, *result.Options.XMLDisableAutoConvert)
		assert.Equal(t, true, *result.Options.XMLSkipByteOrderMark)

		describeResult, err := client.FileFormats.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, FileFormatTypeXML, describeResult.Type)
		assert.Equal(t, XMLCompressionDeflate, *describeResult.Options.XMLCompression)
		assert.Equal(t, true, *describeResult.Options.XMLIgnoreUTF8Errors)
		assert.Equal(t, true, *describeResult.Options.XMLPreserveSpace)
		assert.Equal(t, true, *describeResult.Options.XMLStripOuterElement)
		assert.Equal(t, true, *describeResult.Options.XMLDisableSnowflakeData)
		assert.Equal(t, true, *describeResult.Options.XMLDisableAutoConvert)
		assert.Equal(t, true, *describeResult.Options.XMLSkipByteOrderMark)
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
			Type: FileFormatTypeCSV,
			FileFormatTypeOptions: FileFormatTypeOptions{
				CSVCompression: &CSVCompressionAuto,
				CSVParseHeader: Bool(false),
			},
		})
		t.Cleanup(fileFormatCleanup)

		err := client.FileFormats.Alter(ctx, fileFormat.ID(), &AlterFileFormatOptions{
			Set: &FileFormatTypeOptions{
				CSVCompression: &CSVCompressionBz2,
				CSVParseHeader: Bool(true),
			},
		})
		require.NoError(t, err)

		result, err := client.FileFormats.ShowByID(ctx, fileFormat.ID())
		require.NoError(t, err)
		assert.Equal(t, CSVCompressionBz2, *result.Options.CSVCompression)
		assert.Equal(t, true, *result.Options.CSVParseHeader)
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

func TestInt_FileFormatsShowById(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, cleanupDatabase := createDatabase(t, client)
	t.Cleanup(cleanupDatabase)
	schemaTest, cleanupSchema := createSchema(t, client, databaseTest)
	t.Cleanup(cleanupSchema)
	fileFormatTest, cleanupFileFormat := createFileFormat(t, client, schemaTest.ID())
	t.Cleanup(cleanupFileFormat)

	databaseTest2, cleanupDatabase2 := createDatabase(t, client)
	t.Cleanup(cleanupDatabase2)
	schemaTest2, cleanupSchema2 := createSchema(t, client, databaseTest2)
	t.Cleanup(cleanupSchema2)

	t.Run("show format in different schema", func(t *testing.T) {
		err := client.Sessions.UseDatabase(ctx, databaseTest2.ID())
		require.NoError(t, err)
		err = client.Sessions.UseSchema(ctx, schemaTest2.ID())
		require.NoError(t, err)

		fileFormat, err := client.FileFormats.ShowByID(ctx, fileFormatTest.ID())
		require.NoError(t, err)
		assert.Equal(t, databaseTest.Name, fileFormat.Name.databaseName)
		assert.Equal(t, schemaTest.Name, fileFormat.Name.schemaName)
		assert.Equal(t, fileFormatTest.Name.name, fileFormat.Name.name)
	})
}
