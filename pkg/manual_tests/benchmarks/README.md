# Authentication methods manual tests

This directory is dedicated to hold steps for manual performance tests in the provider. These tests use simple Terraform files and are run with `terraform` CLI manually to imitate the user workflow and reduce bias with Terraform SDK testing libraries and the binary itself.
The tests are organized by the resource type, e.g. schemas, users, and warehouses.

## Run tests

- Preferably use your secondary test account to avoid potential conflicts with our "main" environments.
- Configure the needed modules in `main.tf`. If you do not want to uses resource from a module, simply set `resource_count` to 0. Note that this field refers to a one "type" of the tests, meaning that one resources can have a few variations (set up dependencies and filled optional fields).
- Run `terraform init -upgrade` to enable the modules.
- Run regular Terraform commands, like `terraform apply`.
- Do not forget to remove the resources with `terraform destroy`.
- To speed up the commands, you can use `-refresh=false` and `-parallelism=N` (default is 10).

## State size

After running the `terraform` commands, the state file should be saved at `terraform.tfstate`. This file can be analyzed in terms of file size.

Run the following command to capture state size:

```bash
ls -lh terraform.tfstate
```

To check potential size reduction with removed parameters, first remove parameters from the state (with using [jq](https://github.com/jqlang/jq)):

```bash
jq 'del(.resources[].instances[].attributes.parameters)' terraform.tfstate > terraform_without_parameters.tfstate
```

And capture the size of the new state.

```bash
ls -lh terraform_without_parameters.tfstate
```
