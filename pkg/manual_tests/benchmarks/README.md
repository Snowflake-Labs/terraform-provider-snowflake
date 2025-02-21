# Manual performance benchmarks

This directory is dedicated to hold steps for manual performance tests in the provider. These tests use simple Terraform files and are run with `terraform` CLI manually to imitate the user workflow and reduce bias with Terraform SDK testing libraries and the binary itself.
The tests are organized by the resource type, e.g. schemas, tasks, and warehouses.

## Run tests

- Preferably use your secondary test account to avoid potential conflicts with our "main" environments.
- Configure the needed modules in `main.tf`. The attribute `resource_count` is a count of the resources of a given type and configuration. Note that inside the modules there are resources with different configurations (i.e. only required fields set, all fields set). This means that the total number of resources may be bigger. For example, if `resource_count` is 100 and you are testing 2 different configurations using the `resource_count`, total number of resources is 200. If you do not want to uses resource from a module, simply set `resource_count` to 0.
- If you want to test different resource configurations, adjust them in the relevant module.
- Run `terraform init -upgrade` to enable the modules.
- Run regular Terraform commands, like `terraform apply`.
- The top-level objects names contain test ID and resource index, utilizing the format like `PERFORMANCE_TESTS_BED9310F_F8CE_D2CD_D4B6_B82F56D6FD42_BASIC_0`. The test ID is regenerated for every run.
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
