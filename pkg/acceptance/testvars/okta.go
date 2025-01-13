package testvars

import "net/url"

var (
	ExampleOktaUrlString = "https://example-tf.com"
	ExampleOktaUrl, _    = url.Parse(ExampleOktaUrlString)
)

var (
	ExampleOktaUrlFromEnvString = "https://example-tf-env.com"
	ExampleOktaUrlFromEnv, _    = url.Parse(ExampleOktaUrlFromEnvString)
)
