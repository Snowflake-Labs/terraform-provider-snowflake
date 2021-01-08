package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipeIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id := "database_name|schema_name|pipe"
	pipe, err := pipeIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", pipe.DatabaseName)
	r.Equal("schema_name", pipe.SchemaName)
	r.Equal("pipe", pipe.PipeName)

	// Bad ID -- not enough fields
	id = "database"
	_, err = pipeIDFromString(id)
	r.Equal(fmt.Errorf("3 fields allowed"), err)

	// Bad ID
	id = "||"
	_, err = pipeIDFromString(id)
	r.NoError(err)

	// 0 lines
	id = ""
	_, err = pipeIDFromString(id)
	r.Equal(fmt.Errorf("1 line per pipe"), err)

	// 2 lines
	id = `database_name|schema_name|pipe
	database_name|schema_name|pipe`
	_, err = pipeIDFromString(id)
	r.Equal(fmt.Errorf("1 line per pipe"), err)
}

func TestPipeStruct(t *testing.T) {
	r := require.New(t)

	// Vanilla
	pipe := &pipeID{
		DatabaseName: "database_name",
		SchemaName:   "schema_name",
		PipeName:     "pipe",
	}
	sID, err := pipe.String()
	r.NoError(err)
	r.Equal("database_name|schema_name|pipe", sID)

	// Empty grant
	pipe = &pipeID{}
	sID, err = pipe.String()
	r.NoError(err)
	r.Equal("||", sID)

	// Grant with extra delimiters
	pipe = &pipeID{
		DatabaseName: "database|name",
		PipeName:     "pipe|name",
	}
	sID, err = pipe.String()
	r.NoError(err)
	newPipe, err := pipeIDFromString(sID)
	r.NoError(err)
	r.Equal("database|name", newPipe.DatabaseName)
	r.Equal("pipe|name", newPipe.PipeName)
}

const pipeCopyStatementTemplate string = "COPY INTO MY_DATABASE.MY_SCHEMA.%[4]s (%[1]s%[2]sID%[1]s%[2]s,VALUE%[1]s) FROM (%[1]s%[2]sSELECT%[1]s%[2]s%[2]sSRC.$1%[1]s%[2]s%[2]s,SRC.$2%[1]s%[2]sFROM @MY_DATABASE.MY_SCHEMA.MY_STAGE AS SRC%[1]s)%[1]sFILE_FORMAT = (%[1]s%[2]sFORMAT_NAME = MY_DATABASE.MY_SCHEMA.JSON%[1]s)%[1]sON_ERROR = CONTINUE%[3]s"

func generatecopyStatement(lineEnding string, indent string, includeSemiColon bool, tableName string) string {
	semiColon := ""

	if includeSemiColon {
		semiColon = ";"
	}

	return fmt.Sprintf(
		pipeCopyStatementTemplate,
		lineEnding,
		indent,
		semiColon,
		tableName,
	)
}

func TestPipeCopyStatementDiffSuppress(t *testing.T) {
	type testCaseData struct {
		declared  string
		showPipes string
		expected  bool
	}

	testCases := map[string]testCaseData{
		"TestDiffSuppressSingleLine": {
			declared:  generatecopyStatement("", " ", false, "MY_TABLE"),
			showPipes: generatecopyStatement("", " ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressMultiLineUnix": {
			declared:  generatecopyStatement("\n", "  ", false, "MY_TABLE"),
			showPipes: generatecopyStatement("\n", "  ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressMultiLineWindows": {
			declared:  generatecopyStatement("\n", "  ", false, "MY_TABLE"),
			showPipes: generatecopyStatement("\r\n", "  ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressSingleLineWithEndStatement": {
			declared:  generatecopyStatement("", " ", true, "MY_TABLE"),
			showPipes: generatecopyStatement("", " ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressMultiLineUnixWithEndStatement": {
			declared:  generatecopyStatement("\n", "  ", true, "MY_TABLE"),
			showPipes: generatecopyStatement("\n", "  ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressMultiLineWindowsWithEndStatement": {
			declared:  generatecopyStatement("\n", "  ", true, "MY_TABLE"),
			showPipes: generatecopyStatement("\r\n", "  ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestNoDiffSuppressSingleLine": {
			declared:  generatecopyStatement("", " ", false, "MY_TABLE"),
			showPipes: generatecopyStatement("", " ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressMultiLineUnix": {
			declared:  generatecopyStatement("\n", "  ", false, "MY_TABLE"),
			showPipes: generatecopyStatement("\n", "  ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressMultiLineWindows": {
			declared:  generatecopyStatement("\n", "  ", false, "MY_TABLE"),
			showPipes: generatecopyStatement("\r\n", "  ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressSingleLineWithEndStatement": {
			declared:  generatecopyStatement("", " ", true, "MY_TABLE"),
			showPipes: generatecopyStatement("", " ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressMultiLineUnixWithEndStatement": {
			declared:  generatecopyStatement("\n", "  ", true, "MY_TABLE"),
			showPipes: generatecopyStatement("\n", "  ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressMultiLineWindowsWithEndStatement": {
			declared:  generatecopyStatement("\n", "  ", true, "MY_TABLE"),
			showPipes: generatecopyStatement("\r\n", "  ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(testCase.expected, pipeCopyStatementDiffSuppress("", testCase.declared, testCase.showPipes, nil))
		})
	}
}
