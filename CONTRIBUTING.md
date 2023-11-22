# Contributing

## Development

1. Install Go (eg. `brew install golang` on MacOS)
2. Ensure that your `GOPATH` is set to the desired location
3. Fork this repo and clone it into `$GOPATH/src/github.com/Snowflake-Labs/terraform-provider-snowflake`
4. cd to `terraform-provider-snowflake` and install all the required packages with `go get`
5. Build and install provider with `make install`

## Testing

The following environment variables need to be set for acceptance tests to run:

* `SNOWFLAKE_ACCOUNT` - The account name
* `SNOWFLAKE_USER` - A snowflake user for running tests.
* `SNOWFLAKE_PASSWORD` - Password for that user.
* `SNOWFLAKE_ROLE` - Needs to be ACCOUNTADMIN or similar.
* `SNOWFLAKE_REGION` - Default is us-west-2, set this if your snowflake account is in a different region.
* `TF_ACC` - to enable acceptance tests.

For example:

```sh
export SNOWFLAKE_ACCOUNT=TESTACCOUNT
export SNOWFLAKE_USER=TEST_USER
export SNOWFLAKE_PASSWORD=hunter2
export SNOWFLAKE_ROLE=ACCOUNTADMIN
export SNOWFLAKE_REGION=us-west-2
export TF_ACC=true
```

You can also read the config from a `~/.snowflake/config` file, although you will still need to set `TF_ACC` to true.


~/.snowflake/config
```sh
[default]
account='TESTACCOUNT'
user='TEST_USER'
password='hunter2'
role='ACCOUNTADMIN'
```

**Note: PRs for new resources will not be accepted without passing acceptance tests.**

For the Terraform resources, there are 3 levels of testing - internal, unit and acceptance tests.

The 'internal' tests are run in the `github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources` package so that they can test functions that are not exported. These tests are intended to be limited to unit tests for simple functions.

The 'unit' tests are run in  `github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources_test`, so they only have access to the exported methods of `resources`. These tests exercise the CRUD methods that on the terraform resources. Note that all tests here make use of database mocking and are run locally. This means the tests are fast, but are liable to be wrong in subtle ways (since the mocks are unlikely to be perfect).

You can run these first two sets of tests with `make test`.

The 'acceptance' tests run the full stack, creating, modifying and destroying resources in a live snowflake account. To run them you need a snowflake account and the proper authentication set up. These tests are slower but have higher fidelity. You can create a new Snowflake Enterprise trial account and setup the environment variables for running acceptance tests.

To run all tests, including the acceptance tests, run `make test-acceptance`.

### Running tests in VSCode

If you're using VSCode, this project comes pre-configured to source the `test.env` file before each test so you can run acceptance tests directly for the editor.
We've included an example env file `test.env.example` with the environment variables described above so you can set up your acceptance tests by:
- running `cp test.env.example test.env`
- editing `test.env` with your Snowflake account values
- installing the `golang.go` extension
- running tests directly from VSCode!

## Advanced Debugging

If you want to build and test the provider locally you should edit you `~.terraformrc` file to include the following:

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/Snowflake-Labs/snowflake" = "<path_to_binary>"
  }

  direct {}
}
```

To debug the provider with a debugger:

1. Launch the provider with the `-debug` command line argument in your debugger session. Once the provider starts, it will print instructions on setting the `TF_REATTACH_PROVIDERS` environment variable.

   ```sh
   Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

   Command Prompt:	set "TF_REATTACH_PROVIDERS={"registry.terraform.io/Snowflake-Labs/snowflake":{"Protocol":"grpc","ProtocolVersion":5,"Pid":35140,"Test":true,"Addr": {"Network":"tcp","String":"127.0.0.1:54706"}}}"
   PowerShell:	$env:TF_REATTACH_PROVIDERS='{"registry.terraform.io/Snowflake-Labs/snowflake":{"Protocol":"grpc","ProtocolVersion":5,"Pid":35140,"Test":true,"Addr":{"Network":"tcp","String":"127.0.0.1:54706"}}}'
   ```

2. Open a terminal where you will execute Terraform and set the `TF_REATTACH_PROVIDERS` environment variable using the command from the first step.
3. Run Terraform as usual from this terminal. Any breakpoints you set will halt execution and you can troubleshoot the provider from your debugger.

**Note**: The `TF_REATTACH_PROVIDERS` environment variable needs to be set every time you restart your debugger session as some values like the `Pid` or the TCP port will change with every execution.

For further instructions, please check the official [Terraform Plugin Development guide](https://www.terraform.io/plugin/debugging#starting-a-provider-in-debug-mode).

## Contributing

We use [Conventional Commits](https://www.conventionalcommits.org/) for commit message formatting and PR titles. Please try to adhere to the standard.
Refer to the [regular expression](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/.github/workflows/titleLint.yml#L17) for PR title validation.

## Releasing

Releases will be performed as needed, typically once every 1-2 weeks. If your change is more urgent and you need to use it sooner, use the commit hash.

Releases are done by [goreleaser](https://goreleaser.com/) and run by our make files. There two goreleaser configs, `.goreleaser.yml` for regular releases and `.goreleaser.prerelease.yml` for doing prereleases (for testing).

Releases are [published to the terraform registry](https://registry.terraform.io/providers/chanzuckerberg/snowflake/latest), which requires that releases by signed.
