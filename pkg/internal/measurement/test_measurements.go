package measurement

import (
	"fmt"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
	"testing"
	"time"
)

type TestMeasurement struct {
	TestName        string
	Duration        time.Duration
	SubMeasurements map[string]*TestMeasurement
}

func NewTestMeasurement(t *testing.T) *TestMeasurement {
	m := new(TestMeasurement)
	m.TestName = t.Name()
	m.SubMeasurements = make(map[string]*TestMeasurement)
	return m
}

type TestMeasurements struct {
	Measurements map[string]*TestMeasurement
}

func NewTestMeasurements() *TestMeasurements {
	return &TestMeasurements{
		Measurements: make(map[string]*TestMeasurement),
	}
}

func PrintMeasurementSummary(measurements *TestMeasurements, truncateDuration time.Duration) {
	fmt.Println("Summary of test measurements (sorted by execution time)")
	values := maps.Values(measurements.Measurements)
	slices.SortFunc(values, func(a, b *TestMeasurement) int {
		return int(b.Duration.Nanoseconds() - a.Duration.Nanoseconds())
	})
	printMeasurements(0, values, truncateDuration)
}

func printMeasurements(level int, measurements []*TestMeasurement, truncateDuration time.Duration) {
	for _, measurement := range measurements {
		fmt.Printf("%-10s%s%s\n", measurement.Duration.Truncate(truncateDuration), strings.Repeat(" → ", level), measurement.TestName)
		printMeasurements(level+1, maps.Values(measurement.SubMeasurements), truncateDuration)
	}
}

func MeasureTestTime(testMeasurements *TestMeasurements, t *testing.T) {
	t.Helper()
	start := time.Now()
	m := NewTestMeasurement(t)
	testNameParts := strings.Split(m.TestName, "/")
	if len(testNameParts) == 1 {
		if _, ok := testMeasurements.Measurements[m.TestName]; ok {
			panic("this shouldn't happen")
		}

		testMeasurements.Measurements[m.TestName] = m
	} else {
		currentMeasurement := testMeasurements.Measurements[testNameParts[0]]
		for index, part := range testNameParts[1:] {
			if value, ok := currentMeasurement.SubMeasurements[part]; ok {
				currentMeasurement = value
			} else {
				m.TestName = strings.Join(testNameParts[index+1:], "/")
				currentMeasurement.SubMeasurements[part] = m
			}
		}
	}
	t.Cleanup(func() { m.Duration = time.Now().Sub(start) })
}
