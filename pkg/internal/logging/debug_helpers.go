package logging

import (
	"io"
	"log"
	pkgos "os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/os"
)

// TODO (next PRs): remove extra logging
func init() {
	additionalDebugLoggingEnabled = os.Getenv("SF_TF_ADDITIONAL_DEBUG_LOGGING") != ""
	DebugLogger = log.New(pkgos.Stderr, "sf-tf-additional-debug ", log.LstdFlags|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	if !additionalDebugLoggingEnabled {
		DebugLogger.SetOutput(io.Discard)
	}
}

var (
	additionalDebugLoggingEnabled bool
	DebugLogger                   *log.Logger
)
