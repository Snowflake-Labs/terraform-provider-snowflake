package provider

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/gookit/color"
)

type tfOperation string

const (
	CreateOperation tfOperation = "CREATE"
	ReadOperation   tfOperation = "READ"
	UpdateOperation tfOperation = "UPDATE"
	DeleteOperation tfOperation = "DELETE"
)

func formatSQLPreview(operation tfOperation, resourceName string, id string, commands []string) string {
	var c color.Color
	switch operation {
	case CreateOperation:
		c = color.HiGreen
	case ReadOperation:
		c = color.HiBlue
	case UpdateOperation:
		c = color.HiYellow
	case DeleteOperation:
		c = color.HiRed
	}
	var sb strings.Builder
	sb.WriteString(c.Sprintf("\n[ %s %s %s ]", operation, resourceName, id))
	for _, command := range commands {
		sb.WriteString(c.Sprintf("\n  - %s", command))
	}
	sb.WriteString("\n")
	return sb.String()
}

type sensitiveAttributes struct {
	m map[string]bool
}

var (
	sa   *sensitiveAttributes
	lock = sync.Mutex{}
)

func isSensitive(s string) bool {
	if sa == nil {
		lock.Lock()
		defer lock.Unlock()
		if sa == nil {
			sa = &sensitiveAttributes{
				m: make(map[string]bool),
			}
			dir, err := os.UserHomeDir()
			if err != nil {
				return false
			}
			// sensitive path is ~/.snowflake/sensitive.
			f := filepath.Join(dir, ".snowflake", "sensitive")
			dat, err := os.ReadFile(f)
			if err != nil {
				return false
			}
			lines := strings.Split(string(dat), "\n")
			r := regexp.MustCompile("(data[.])?snowflake_(.*)[.](.+)[.](.+)")
			for _, line := range lines {
				strippedLine := strings.TrimSpace(line)
				if r.MatchString(strippedLine) {
					sa.m[strippedLine] = true
				}
			}
		}
	}
	if _, ok := sa.m[s]; ok {
		return true
	}
	return false
}
