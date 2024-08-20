# Manual tests

This directory is dedicated to hold steps for manual tests that are not possible to re-recreate in automated acceptance tests.
The main limitations come from using [terraform-plugin-testing](https://github.com/hashicorp/terraform-plugin-testing) which is
not supporting every action you are able to perform with Terraform CLI. 

Here's the list of cases we currently cannot reproduce and write acceptance tests for:
- When upgrading from version to version we need to remove the state of the deprecated object (terraform state rm) and import a new type representing the same object (terraform import)
  - Specifically, `terraform state rm` is not possible. As an example we can have `snowflake_database` in version 0.92.0 which in version 0.93.0 could be represented as e.g. `snowflake_shared_database`. To fully test such upgrade path, it has to be done manually.
  - Currently tests under this category:
    - `upgrade_cloned_database`
    - `upgrade_secondary_database`
    - `upgrade_shared_database`

## How to use manual tests
- Choose the test you want to run and go into the test folder.
- Take the first step from that test copy it into a separate folder (outside the project) to initialize terraform and start the first step.
  - The tests contain provider configuration, but you have to make sure you have compliant configuration in your `~/snowflake/config` file.
  - Please mind which commands should be run to perform a given test step correctly (instructions are at the top of the file; run after project is initialized).
  - Also, please mind `TODO: Replace` comments that indicate lines where configuration should be changed before running any test command.
- To proceed with the test, take the content of the next test step file and replace the previous one (with the same rules as above).