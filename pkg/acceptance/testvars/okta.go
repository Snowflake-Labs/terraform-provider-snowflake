package testvars

import "net/url"

var ExampleOktaUrlString = "https://example-tf.com"
var ExampleOktaUrl, _ = url.Parse(ExampleOktaUrlString)

var ExampleOktaUrlFromEnvString = "https://example-tf-env.com"
var ExampleOktaUrlFromEnv, _ = url.Parse(ExampleOktaUrlFromEnvString)
