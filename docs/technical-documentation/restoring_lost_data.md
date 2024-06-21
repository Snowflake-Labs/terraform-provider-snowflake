# Restoring Lost Data with Time-Travel

If you've ever accidentally deleted important data either by using Terraform or by hand, there's still hope to recover the data.
By using the Snowflake's Time-Travel feature, you can restore lost data and undo those accidental deletions.

> Note: Currently, the recovery process is predominantly manual, relying on SQL commands and the Terraform CLI. 
We made a strategic decision not to integrate it as a provider feature at this time, as demand for this functionality was not significant.
Following the release of V1, we intend to reassess the topic of data recovery and UNDROP functionality to explore potential integration into the provider, evaluating its necessity and feasibility.

You should be prepared beforehand by specifying how much of the historical data Snowflake should keep by setting the [DATA_RETENTION_TIME_IN_DAYS](https://docs.snowflake.com/en/sql-reference/parameters#data-retention-time-in-days) parameter.
When using our provider, you can set this by using one of our parameter-setting resources (like [snowflake_account_parameter](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/account_parameter) or [snowflake_object_parameter](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/object_parameter))
or set it on the resource level (e.g. `data_retention_time_in_days` in [snowflake_database](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/database)).

> Note: If some of the resources support `data_retention_time_in_days` parameter in Snowflake, but it's not available in the provider, we'll add it during [the resource preparation for V1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1).

Now, with [DATA_RETENTION_TIME_IN_DAYS](https://docs.snowflake.com/en/sql-reference/parameters#data-retention-time-in-days) set up, 
let's imagine we accidentally dropped a database that was managed by Terraform and contained a lot of important data we would like to recover.
But before we start, let's clearly understand the initial state of the database we'll be recovering.

The configuration will contain only a database as we want to mainly focus on the data recovery from the Terraform point of view:
```terraform
resource "snowflake_database" "test" {
  name = "TEST_DATABASE"
}
```

The rest of the objects will be created in SQL by hand (e.g. in the worksheet):
```sql
CREATE SCHEMA TEST_DATABASE.TEST_SCHEMA;
CREATE TABLE TEST_DATABASE.TEST_SCHEMA.TEST_TABLE(id INT, name STRING);
INSERT INTO TEST_DATABASE.TEST_SCHEMATEST_TABLE(id, name) VALUES (0, 'john'), (1, 'doe');
```
As you can see, we created a schema and a table with "important" data.

As the next step we will remove the database from the configuration and run `terraform apply` that will trigger the `DROP DATABASE`.
Now, we're in a state where the data is lost and can only be recovered by Time-Travel. By checking up the [DATA_RETENTION_TIME_IN_DAYS](https://docs.snowflake.com/en/sql-reference/parameters#data-retention-time-in-days)
parameter, we can calculate if the dropped database can still be recovered or not (that's why it's important to prepare for such situations beforehand,
to avoid situations where the parameter was set to low value and the data is lost).

To recover the database (and the data inside it), we have to call `UNDROP DATABASE TEST_DATABASE` manually.
To bring the database back to the Terraform configuration, we have to specify the same configuration as previously, but now, 
instead of running `terraform apply` we have to import it by calling `terraform import 'TEST_DATABASE'`.
After successful import the `terraform plan` shouldn't produce any plan for the database. 
To ensure all the important data we inserted before is there, we can call `SELECT * FROM TEST_DATABASE.TEST_SCHEMA.TEST_TABLE;`.
 
To learn more about how to use Time-Travel, check out the links below:
1. [Understanding & using Time-Travel](https://docs.snowflake.com/en/user-guide/data-time-travel)
2. [DATA_RETENTION_TIME_IN_DAYS parameter](https://docs.snowflake.com/en/sql-reference/parameters#data-retention-time-in-days)
3. [UNDROP command with available objects to restore](https://docs.snowflake.com/en/sql-reference/sql/undrop)
