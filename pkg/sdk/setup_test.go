package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/measurement"
	"os"
	"testing"
	"time"
)

var testMeasurements = measurement.NewTestMeasurements()

func measureTest(t *testing.T) {
	t.Helper()
	measurement.MeasureTestTime(testMeasurements, t)
}

func TestMain(m *testing.M) {
	code := m.Run()
	measurement.PrintMeasurementSummary(testMeasurements, time.Nanosecond)
	os.Exit(code)
}
