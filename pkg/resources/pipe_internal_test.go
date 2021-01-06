package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

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
