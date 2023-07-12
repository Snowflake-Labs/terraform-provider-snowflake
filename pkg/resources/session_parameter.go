package resources

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/exp/maps"
)

var sessionParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Name of session parameter. Valid values are those in [session parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#session-parameters).",
		ValidateFunc: validation.StringInSlice(maps.Keys(snowflake.GetParameterDefaults(snowflake.ParameterTypeSession)), false),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of session parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
	"on_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the session parameter will be set on the account level.",
	},
	"user": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The user to set the session parameter for. Required if on_account is false",
	},
}

func SessionParameter() *schema.Resource {
	return &schema.Resource{
		Create: CreateSessionParameter,
		Read:   ReadSessionParameter,
		Update: UpdateSessionParameter,
		Delete: DeleteSessionParameter,

		Schema: sessionParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSessionParameter implements schema.CreateFunc.
func CreateSessionParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	onAccount := d.Get("on_account").(bool)
	user := d.Get("user").(string)
	parameter := sdk.SessionParameter(key)
	builder := snowflake.NewSessionParameter(key, value, db)

	var err error
	if onAccount {
		opts, err := setSessionParameter(parameter, value)
		if err != nil {
			return err
		}
		err = client.Accounts.Alter(ctx, opts)
		if err != nil {
			return err
		}
	} else {
		if user == "" {
			return fmt.Errorf("user is required if on_account is false")
		}
		builder.SetUser(user)
		err = builder.SetParameter()
		if err != nil {
			return fmt.Errorf("error creating session parameter err = %w", err)
		}
	}

	d.SetId(key)

	return ReadSessionParameter(d, meta)
}

func setSessionParameter(parameter sdk.SessionParameter, value string) (*sdk.AlterAccountOptions, error) {
	opts := sdk.AlterAccountOptions{Set: &sdk.AccountSet{Parameters: &sdk.AccountLevelParameters{SessionParameters: &sdk.SessionParameters{}}}}
	switch parameter {
	case sdk.SessionParameterAbortDetachedQuery:
		if value == "true" {
			opts.Set.Parameters.SessionParameters.AbortDetachedQuery = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.SessionParameters.AbortDetachedQuery = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("ABORT_DETACHED_QUERY session parameter is a boolean value, got: %v", value)
		}
	case sdk.SessionParameterAutocommit:
		if value == "true" {
			opts.Set.Parameters.SessionParameters.Autocommit = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.SessionParameters.Autocommit = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("AUTO_COMMIT session parameter is a boolean value, got: %v", value)
		}
	case sdk.SessionParameterBinaryInputFormat:
		opts.Set.Parameters.SessionParameters.BinaryInputFormat = sdk.Pointer(sdk.BinaryInputFormat(value))
	case sdk.SessionParameterBinaryOutputFormat:
		opts.Set.Parameters.SessionParameters.BinaryOutputFormat = sdk.Pointer(sdk.BinaryOutputFormat(value))
	case sdk.SessionParameterDateInputFormat:
		opts.Set.Parameters.SessionParameters.DateInputFormat = &value
	case sdk.SessionParameterDateOutputFormat:
		opts.Set.Parameters.SessionParameters.DateOutputFormat = &value
	case sdk.SessionParameterErrorOnNondeterministicMerge:
		if value == "true" {
			opts.Set.Parameters.SessionParameters.ErrorOnNondeterministicMerge = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.SessionParameters.ErrorOnNondeterministicMerge = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("ERROR_ON_NONDETERMINISTIC_MERGE session parameter is a boolean value, got: %v", value)
		}
	case sdk.SessionParameterErrorOnNondeterministicUpdate:
		if value == "true" {
			opts.Set.Parameters.SessionParameters.ErrorOnNondeterministicUpdate = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.SessionParameters.ErrorOnNondeterministicUpdate = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("ERROR_ON_NONDETERMINISTIC_UPDATE session parameter is a boolean value, got: %v", value)
		}
	case sdk.SessionParameterGeographyOutputFormat:
		opts.Set.Parameters.SessionParameters.GeographyOutputFormat = sdk.Pointer(sdk.GeographyOutputFormat(value))
	case sdk.SessionParameterJSONIndent:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("JSON_INDENT session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.SessionParameters.JSONIndent = sdk.Pointer(v)
	case sdk.SessionParameterLockTimeout:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("LOCK_TIMEOUT session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.SessionParameters.LockTimeout = sdk.Pointer(v)
	case sdk.SessionParameterQueryTag:
		opts.Set.Parameters.SessionParameters.QueryTag = &value
	case sdk.SessionParameterRowsPerResultset:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("ROWS_PER_RESULTSET session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.SessionParameters.RowsPerResultset = sdk.Pointer(v)
	case sdk.SessionParameterSimulatedDataSharingConsumer:
		opts.Set.Parameters.SessionParameters.SimulatedDataSharingConsumer = &value
	case sdk.SessionParameterStatementTimeoutInSeconds:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("STATEMENT_TIMEOUT_IN_SECONDS session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.SessionParameters.StatementTimeoutInSeconds = sdk.Pointer(v)
	case sdk.SessionParameterStrictJSONOutput:
		if value == "true" {
			opts.Set.Parameters.SessionParameters.StrictJSONOutput = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.SessionParameters.StrictJSONOutput = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("STRICT_JSON_OUTPUT session parameter is a boolean value, got: %v", value)
		}
	case sdk.SessionParameterTimestampDayIsAlways24h:
		if value == "true" {
			opts.Set.Parameters.SessionParameters.TimestampDayIsAlways24h = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.SessionParameters.TimestampDayIsAlways24h = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("TIMESTAMP_DAY_IS_ALWAYS_24H session parameter is a boolean value, got: %v", value)
		}
	case sdk.SessionParameterTimestampInputFormat:
		opts.Set.Parameters.SessionParameters.TimestampInputFormat = &value
	case sdk.SessionParameterTimestampLTZOutputFormat:
		opts.Set.Parameters.SessionParameters.TimestampLTZOutputFormat = &value
	case sdk.SessionParameterTimestampNTZOutputFormat:
		opts.Set.Parameters.SessionParameters.TimestampNTZOutputFormat = &value
	case sdk.SessionParameterTimestampOutputFormat:
		opts.Set.Parameters.SessionParameters.TimestampOutputFormat = &value
	case sdk.SessionParameterTimestampTypeMapping:
		opts.Set.Parameters.SessionParameters.TimestampTypeMapping = &value
	case sdk.SessionParameterTimestampTZOutputFormat:
		opts.Set.Parameters.SessionParameters.TimestampTZOutputFormat = &value
	case sdk.SessionParameterTimezone:
		opts.Set.Parameters.SessionParameters.Timezone = &value
	case sdk.SessionParameterTimeInputFormat:
		opts.Set.Parameters.SessionParameters.TimeInputFormat = &value
	case sdk.SessionParameterTimeOutputFormat:
		opts.Set.Parameters.SessionParameters.TimeOutputFormat = &value
	case sdk.SessionParameterTransactionDefaultIsolationLevel:
		opts.Set.Parameters.SessionParameters.TransactionDefaultIsolationLevel = sdk.Pointer(sdk.TransactionDefaultIsolationLevel(value))
	case sdk.SessionParameterTwoDigitCenturyStart:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("TWO_DIGIT_CENTURY_START session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.SessionParameters.TwoDigitCenturyStart = sdk.Pointer(v)
	case sdk.SessionParameterUnsupportedDDLAction:
		opts.Set.Parameters.SessionParameters.UnsupportedDDLAction = sdk.Pointer(sdk.UnsupportedDDLAction(value))
	case sdk.SessionParameterUseCachedResult:
		if value == "true" {
			opts.Set.Parameters.SessionParameters.UseCachedResult = sdk.Bool(true)
		} else if value == "false" {
			opts.Set.Parameters.SessionParameters.UseCachedResult = sdk.Bool(false)
		} else {
			return nil, fmt.Errorf("USE_CACHED_RESULT session parameter is a boolean value, got: %v", value)
		}
	case sdk.SessionParameterWeekOfYearPolicy:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("WEEK_OF_YEAR_POLICY session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.SessionParameters.WeekOfYearPolicy = sdk.Pointer(v)
	case sdk.SessionParameterWeekStart:

		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("WEEK_START session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.SessionParameters.WeekStart = sdk.Pointer(v)
	default:
		return nil, fmt.Errorf("Invalid session parameter: %v", string(parameter))
	}

	return &opts, nil
}

// ReadSessionParameter implements schema.ReadFunc.
func ReadSessionParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	parameter := d.Id()

	onAccount := d.Get("on_account").(bool)
	var err error
	var p *sdk.Parameter
	if onAccount {
		p, err = client.Sessions.ShowAccountParameter(ctx, sdk.AccountParameter(parameter))
	} else {
		user := d.Get("user").(string)
		userId := sdk.NewAccountObjectIdentifier(user)
		p, err = client.Sessions.ShowUserParameter(ctx, sdk.UserParameter(parameter), userId)
	}
	if err != nil {
		return fmt.Errorf("error reading session parameter err = %w", err)
	}
	err = d.Set("value", p.Value)
	if err != nil {
		return fmt.Errorf("error setting session parameter err = %w", err)
	}
	return nil
}

// UpdateSessionParameter implements schema.UpdateFunc.
func UpdateSessionParameter(d *schema.ResourceData, meta interface{}) error {
	return CreateSessionParameter(d, meta)
}

// DeleteSessionParameter implements schema.DeleteFunc.
func DeleteSessionParameter(d *schema.ResourceData, meta interface{}) error {

	db := meta.(*sql.DB)
	key := d.Get("key").(string)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	onAccount := d.Get("on_account").(bool)
	parameter := sdk.SessionParameter(key)

	var err error
	if onAccount {
		defaultParameter, err := client.Sessions.ShowAccountParameter(ctx, sdk.AccountParameter(key))
		if err != nil {
			return err
		}
		defaultValue := defaultParameter.Default
		opts, err := setSessionParameter(parameter, defaultValue)
		if err != nil {
			return err
		}
		err = client.Accounts.Alter(ctx, opts)
		if err != nil {
			return err
		}
	} else {
		user := d.Get("user").(string)
		if user == "" {
			return fmt.Errorf("user is required if on_account is false")
		}
		parameterDefault := snowflake.GetParameterDefaults(snowflake.ParameterTypeSession)[key]
		defaultValue := parameterDefault.DefaultValue
		value := fmt.Sprintf("%v", defaultValue)
		typeString := reflect.TypeOf("")
		if reflect.TypeOf(parameterDefault.DefaultValue) == typeString {
			value = fmt.Sprintf("'%s'", value)
		}
		builder := snowflake.NewSessionParameter(key, value, db)
		builder.SetUser(user)
		err = builder.SetParameter()
		if err != nil {
			return fmt.Errorf("error creating session parameter err = %w", err)
		}
	}

	d.SetId(key)
	return nil
}
