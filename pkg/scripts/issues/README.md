# Generating the list of open issues
1. To use the script, generate access token here: https://github.com/settings/tokens?type=beta.
2. To get all open issues invoke the [first script](./gh/main.go) setting `SF_TF_SCRIPT_GH_ACCESS_TOKEN`:
```shell
  cd gh && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```
3. File `issues.json` should be generated in the `gh` directory. This is the input file for the second script.
4. To get process the issues invoke the [second script](./file/main.go):
```shell
  cd file && go run .
```
5. File `issues.csv` should be generated in the `file` directory. This is the CSV which summarizes all the issues we have.

# Closing old issues (regarding https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/2755)
1. To use the script, generate access token here: https://github.com/settings/tokens?type=beta.
2. First get all open issues by invoking:
```shell
  cd gh && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```
3. File `issues.json` should be generated in the `gh` directory. This is the input file for the second script. The next script is based also on `presnowflake_bucket.csv` that was created based on the GH issues filtering.
4. To filter only closeable issues invoke [this script](./filter-closeable-old-issues/main.go):
```shell
  cd filter-closeable-old-issues && go run .
```
5. Script will output files `issues_to_close.csv` and `issues_edited.csv`. There are two files documenting closing action on 30.04.2024 (`20240430 - issues_edited.csv` and `20240430 - issues_to_close.csv`). In `20240430 - notes.MD` there are notes regarding the questionable issues and the decisions taken.
6. To close the issues with the appropriate comment provide `issues_to_close.csv` in `close-with-comment` dir. Example `20240430 - issues_to_close.csv` is given. The run:
```shell
  cd close-with-comment && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```