package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Provider_toProtocol(t *testing.T) {
	type test struct {
		input string
		want  protocol
	}

	valid := []test{
		// Case insensitive.
		{input: "http", want: protocolHttp},

		// Supported Values.
		{input: "HTTP", want: protocolHttp},
		{input: "HTTPS", want: protocolHttps},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := toProtocol(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := toProtocol(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Provider_toDriverLogLevel(t *testing.T) {
	type test struct {
		input string
		want  driverLogLevel
	}

	valid := []test{
		// Case insensitive.
		{input: "WARNING", want: logLevelWarning},

		// Supported Values.
		{input: "trace", want: logLevelTrace},
		{input: "debug", want: logLevelDebug},
		{input: "info", want: logLevelInfo},
		{input: "print", want: logLevelPrint},
		{input: "warning", want: logLevelWarning},
		{input: "error", want: logLevelError},
		{input: "fatal", want: logLevelFatal},
		{input: "panic", want: logLevelPanic},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := toDriverLogLevel(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := toDriverLogLevel(tc.input)
			require.Error(t, err)
		})
	}
}
