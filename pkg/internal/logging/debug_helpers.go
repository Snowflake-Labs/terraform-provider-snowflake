package logging

import (
	"io"
	"log"
	"os"
)

func init() {
	additionalDebugLoggingEnabled = os.Getenv("SF_TF_ADDITIONAL_DEBUG_LOGGING") != ""
	DebugLogger = log.New(os.Stderr, "sf-tf-additional-debug ", log.LstdFlags|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	if !additionalDebugLoggingEnabled {
		DebugLogger.SetOutput(io.Discard)
	}
}

var additionalDebugLoggingEnabled bool
var DebugLogger *log.Logger
