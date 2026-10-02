package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeIssueTitle(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal title",
			input:    "snowflake_grant fails on 0.92",
			expected: "snowflake_grant fails on 0.92",
		},
		{
			name:     "leading equals",
			input:    "=cmd|'/c powershell'!A1",
			expected: "cmd|'/c powershell'!A1",
		},
		{
			name:     "leading plus",
			input:    "+2+5",
			expected: "2+5",
		},
		{
			name:     "leading minus",
			input:    "-2+5",
			expected: "2+5",
		},
		{
			name:     "leading at",
			input:    "@SUM(1+1)",
			expected: "SUM(1+1)",
		},
		{
			name:     "leading tab",
			input:    "\t=HYPERLINK(\"http://evil\")",
			expected: "HYPERLINK(\"http://evil\")",
		},
		{
			name:     "leading carriage return",
			input:    "\r=cmd",
			expected: "cmd",
		},
		{
			name:     "stacked prefix",
			input:    "\t=@+-title",
			expected: "title",
		},
		{
			name:     "crlf then formula",
			input:    "\r\n=cmd",
			expected: "cmd",
		},
		{
			name:     "equals later in the title",
			input:    "snowflake_grant = fails",
			expected: "snowflake_grant = fails",
		},
		{
			name:     "only formula characters",
			input:    "=+-@\t\r\n",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, sanitizeIssueTitle(tc.input))
		})
	}
}
