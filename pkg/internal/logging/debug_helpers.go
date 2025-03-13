package logging

import (
	"io"
	"log"
	"os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
)

// TODO (next PRs): remove extra logging
func init() {
	additionalDebugLoggingEnabled = oswrapper.Getenv("SF_TF_ADDITIONAL_DEBUG_LOGGING") != ""
	DebugLogger = log.New(os.Stderr, "sf-tf-additional-debug ", log.LstdFlags|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	if !additionalDebugLoggingEnabled {
		DebugLogger.SetOutput(io.Discard)
	}
}

var (
	additionalDebugLoggingEnabled bool
	DebugLogger                   *log.Logger
)
