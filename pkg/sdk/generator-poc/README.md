## Sdk code generation based on blueprint

This is not even full PoC, but something we needed as input to discussion.

For the sake of discussion only create statement was taken into consideration. For now only struct and interface are generated.

### Usage

There are two files ready for generation:
- [alerts_gen_poc.go](../alerts_gen_poc.go) with [blueprint json](blueprints/alert.json)
- [pipes_gen_poc.go](../pipes_gen_poc.go) with [blueprint json](blueprints/pipe.json)

To invoke both generations run:
```shell
make generate-all
```
This will generate two files: `alerts_gen_poc_generated.go` and `pipes_gen_poc_generated.go`.

To invoke only generation of given resource, run either:
```shell
make generate-pipes
```
or
```shell
make generate-alerts
```

### Next possible steps
- generating validations: based on example from [blueprint](blueprints/alert_proposals.json) and [discussion file](../alerts_discussion.go)
- generating validation tests: based on validations above
- generating tests of sql generation from struct

### Discussion outcome
- putting it all to json is still too much verbose; other formats won't help; maybe moving to generation from go files (e.g. from structs specially marked with custom tags could help, but other ideas are also possible)
- it was based only on two simple objects and only on create, so full json would be much longer
- we could benefit from generating more things that just a simple struct now
- it can be left as an input for future discussions on auto responding to changing schemas
