package builder

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func sqlToFieldName(sql string, shouldExport bool) string {
	sqlWords := strings.Split(sql, " ")
	for i, s := range sqlWords {
		if !shouldExport && i == 0 {
			sqlWords[i] = cases.Lower(language.English).String(s)
			continue
		}
		sqlWords[i] = cases.Title(language.English).String(s)
	}
	return strings.Join(sqlWords, "")
}
