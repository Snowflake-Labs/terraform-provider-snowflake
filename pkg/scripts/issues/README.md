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
