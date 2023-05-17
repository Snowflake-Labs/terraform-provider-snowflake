# Snowflake Go SDK

Required environment variables for running integration tests that provision resources in multiple accounts:

SNOWFLAKE_ACCOUNT_SECOND=<org_name>.<account_name> (e.g. `acme.data-dev2`)

## SQL clause types

| ddl tag            | function              | output format                                                                 |
| ------------------ | --------------------- | ----------------------------------------------------------------------------- |
| `ddl:"static"`     | `sqlStaticClause`     | `WORD`                                                                        |
| `ddl:"keyword"`    | `sqlKeywordClause`    | `"WORD"` (quotes configurable)                                                |
| `ddl:"identifier"` | `sqlIdentifierClause` | `"a.b.c"` or `OBJ_TYPE "a.b.c"`                                               |
| `ddl:"parameter"`  | `sqlParameterClause`  | `PARAM = "value"` (quotes configurable) or `PARAM = 2`                        |                                          |
| `ddl:"list"`       | `sqlListClause`       | `WORD (<subclause>, <subclause>)` (WORD, parentheses, separator configurable) |
