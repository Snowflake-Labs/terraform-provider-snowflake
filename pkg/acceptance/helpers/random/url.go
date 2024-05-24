package random

import (
	"fmt"
	"testing"
)

// GenerateURL generates a random valid URL
func GenerateURL(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("https://%s.com", AlphaN(6))
}
