package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/measurement"
	"os"
	"testing"
	"time"
)

var testMeasurements = measurement.NewTestMeasurements()

type T struct {
	*testing.T
}

// TODO: Proposal for handling sub-testing without the need to call measureTest in them (check accounts_test.go to see how it could look like)
// in this case we could call other function e.g. t := ourTest(tt), because we could use it later for other purposes and inside we would call measureTest anyway
func (t *T) Run(name string, test func(t *T)) {
	t.T.Run(name, func(tt *testing.T) {
		measureTest(tt)
		test(&T{tt})
	})
}

func measureTest(t *testing.T) *T {
	t.Helper()
	measurement.MeasureTestTime(testMeasurements, t)
	return &T{t}
}

func TestMain(m *testing.M) {
	code := m.Run()
	measurement.PrintMeasurementSummary(testMeasurements, time.Nanosecond)
	os.Exit(code)
}
