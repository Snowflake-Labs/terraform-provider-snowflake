package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// ParameterType is the type of parameter.
type ParameterType string

const (
	ParameterTypeAccount ParameterType = "ACCOUNT"
	ParameterTypeSession ParameterType = "SESSION"
	ParameterTypeObject  ParameterType = "OBJECT"
)

// ParameterDefault is a parameter that can be set on an account, session, or object.
type ParameterDefault struct {
	TypeSet            []ParameterType
	DefaultValue       interface{}
	ValueType          reflect.Type
	Validate           func(string) error
	AllowedObjectTypes []ObjectType
}

// ParameterDefaults returns a map of default values for all parameters.
func ParameterDefaults() map[string]ParameterDefault {
	validateBoolFunc := func(value string) (err error) {
		_, err = strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("%v is an invalid value. Boolean value (\"true\"/\"false\") expected", value)
		}
		return nil
	}

	return map[string]ParameterDefault{
		"ALLOW_CLIENT_MFA_CACHING": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"ALLOW_ID_TOKEN": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"CLIENT_ENCRYPTION_KEY_SIZE": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: 128,
			Validate: func(value string) (err error) {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("%v cannot be cast to an integer", value)
				}
				if v != 128 && v != 256 {
					return fmt.Errorf("%v is not a valid value for CLIENT_ENCRYPTION_KEY_SIZE", value)
				}
				return nil
			},
		},
		"ENABLE_INTERNAL_STAGES_PRIVATELINK": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"ENFORCE_SESSION_POLICY": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: true,
			Validate:     validateBoolFunc,
		},
		"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: 10.0,
			Validate: func(value string) (err error) {
				v, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return fmt.Errorf("%v cannot be cast to a float", value)
				}
				if v < 0.0 || (v < 0.0 && v < 1.0) {
					return fmt.Errorf("%v must be 0.0 and above with a scale of at least 1 (e.g. 20.5, 32.25, 33.333, etc.)", v)
				}
				return nil
			},
		},
		"MIN_DATA_RETENTION_TIME_IN_DAYS": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: 0,
			Validate: func(value string) (err error) {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("%v cannot be cast to an integer", value)
				}
				if v < 0 || v > 90 {
					return fmt.Errorf("%v must be 0 or 1 for Standard Edition, or between 0 and 90 for Enterprise Edition or higher", v)
				}
				return nil
			},
		},
		"NETWORK_POLICY": {
			TypeSet:      []ParameterType{ParameterTypeAccount, ParameterTypeObject},
			DefaultValue: "none",
			Validate: func(value string) (err error) {
				if len(value) == 0 {
					return fmt.Errorf("NETWORK_POLICY cannot be empty")
				}
				_, errs := ValidateIdentifier(value, []string{})
				if len(errs) > 0 {
					return fmt.Errorf("NETWORK_POLICY %v is not a valid identifier", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeUser,
			},
		},
		"PERIODIC_DATA_REKEYING": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"PREVENT_UNLOAD_TO_INLINE_URL": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"SSO_LOGIN_PAGE": {
			TypeSet:      []ParameterType{ParameterTypeAccount},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"ABORT_DETACHED_QUERY": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"AUTOCOMMIT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: true,
			Validate:     validateBoolFunc,
		},
		"BINARY_INPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "HEX",
			Validate: func(value string) (err error) {
				validFormats := []string{"HEX", "BASE64", "UTF8", "UTF-8"}
				if !slices.Contains(validFormats, value) {
					return fmt.Errorf("%v is not a valid value for BINARY_INPUT_FORMAT", value)
				}
				return nil
			},
		},
		"BINARY_OUTPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "HEX",
			Validate: func(value string) (err error) {
				validFormats := []string{"HEX", "BASE64"}
				if !slices.Contains(validFormats, value) {
					return fmt.Errorf("%v is not a valid value for BINARY_OUTPUT_FORMAT", value)
				}
				return nil
			},
		},
		"DATE_INPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "auto",
			Validate: func(value string) (err error) {
				validFormats := getValidDateFormats(DateFormatAny, true)
				if !slices.Contains(validFormats, value) {
					return fmt.Errorf("%v is not a valid value for DATE_INPUT_FORMAT", value)
				}
				return nil
			},
		},
		"DATE_OUTPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "YYYY-MM-DD",
			Validate: func(value string) (err error) {
				validFormats := getValidDateFormats(DateFormatAny, false)
				if !slices.Contains(validFormats, value) {
					return fmt.Errorf("%v is not a valid value for DATE_INPUT_FORMAT", value)
				}
				return nil
			},
		},
		"ERROR_ON_NONDETERMINISTIC_MERGE": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: true,
			Validate:     validateBoolFunc,
		},
		"ERROR_ON_NONDETERMINISTIC_UPDATE": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"JSON_INDENT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: 2,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for JSON_INDENT, must be an integer between 0 and 16", value)
				}
				if v < 0 || v > 16 {
					return fmt.Errorf("%v is not a valid value for JSON_INDENT, must be an integer between 0 and 16", value)
				}
				return nil
			},
		},
		"LOCK_TIMEOUT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: 43200,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for LOCK_TIMEOUT, must be an integer", value)
				}
				if v < 0 {
					return fmt.Errorf("%v is not a valid value for LOCK_TIMEOUT, must be an integer greater than 0", value)
				}
				return nil
			},
		},
		"QUERY_TAG": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "",
			Validate: func(value string) (err error) {
				if len(value) > 2000 {
					return fmt.Errorf("%v is not a valid value for QUERY_TAG, must be 2000 characters or less", value)
				}
				return nil
			},
		},
		"ROWS_PER_RESULTSET": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: 0,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for LOCK_TIMEOUT, must be an integer", value)
				}
				if v < 0 {
					return fmt.Errorf("%v is not a valid value for LOCK_TIMEOUT, must be an integer greater than 0", value)
				}
				return nil
			},
		},
		"SIMULATED_DATA_SHARING_CONSUMER": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "",
			Validate:     nil,
		},
		"STATEMENT_TIMEOUT_IN_SECONDS": {
			TypeSet:      []ParameterType{ParameterTypeSession, ParameterTypeObject},
			DefaultValue: 172800,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for STATEMENT_TIMEOUT_IN_SECONDS, must be an integer", value)
				}
				if v < 0 || v > 604800 {
					return fmt.Errorf("%v is not a valid value for STATEMENT_TIMEOUT_IN_SECONDS, must be an integer between 0 and 604800", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeWarehouse,
			},
		},
		"STRICT_JSON_OUTPUT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"TIMESTAMP_DAY_IS_ALWAYS_24H": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: false,
			Validate:     validateBoolFunc,
		},
		"TIMESTAMP_INPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "auto",
			Validate: func(value string) (err error) {
				formats := getValidTimeStampFormats(TimeStampFormatAny, true)
				if !slices.Contains(formats, value) {
					return fmt.Errorf("%v is not a valid value for TIMESTAMP_INPUT_FORMAT, must be one of %v", value, formats)
				}
				return nil
			},
		},
		"TIMESTAMP_LTZ_OUTPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM",
			Validate: func(value string) (err error) {
				formats := getValidTimeStampFormats(TimeStampFormatAny, false)
				if !slices.Contains(formats, value) {
					return fmt.Errorf("%v is not a valid value for TIMESTAMP_LTZ_OUTPUT_FORMAT, must be one of %v", value, formats)
				}
				return nil
			},
		},
		"TIMESTAMP_NTZ_OUTPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3",
			Validate: func(value string) (err error) {
				formats := getValidTimeStampFormats(TimeStampFormatAny, false)
				if !slices.Contains(formats, value) {
					return fmt.Errorf("%v is not a valid value for TIMESTAMP_NTZ_OUTPUT_FORMAT, must be one of %v", value, formats)
				}
				return nil
			},
		},
		"TIMESTAMP_OUTPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM",
			Validate: func(value string) (err error) {
				formats := getValidTimeStampFormats(TimeStampFormatAny, false)
				if !slices.Contains(formats, value) {
					return fmt.Errorf("%v is not a valid value for TIMESTAMP_OUTPUT_FORMAT, must be one of %v", value, formats)
				}
				return nil
			},
		},
		"TIMESTAMP_TYPE_MAPPING": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "TIMESTAMP_NTZ",
			Validate: func(value string) (err error) {
				if !slices.Contains([]string{"TIMESTAMP_NTZ", "TIMESTAMP_LTZ", "TIMESTAMP_TZ"}, value) {
					return fmt.Errorf("%v is not a valid value for TIMESTAMP_TYPE_MAPPING, must be one of TIMESTAMP_NTZ, TIMESTAMP_LTZ, TIMESTAMP_TZ", value)
				}
				return nil
			},
		},
		"TIMESTAMP_TZ_OUTPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "",
			Validate: func(value string) (err error) {
				formats := getValidTimeStampFormats(TimeStampFormatAny, false)
				if !slices.Contains(formats, value) {
					return fmt.Errorf("%v is not a valid value for TIMESTAMP_TZ_OUTPUT_FORMAT, must be one of %v", value, formats)
				}
				return nil
			},
		},
		"TIMEZONE": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "America/Los_Angeles",
			Validate: func(value string) (err error) {
				_, err = time.LoadLocation(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for TIMEZONE, must be a valid timezone", value)
				}
				return nil
			},
		},
		"TIME_INPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "auto",
			Validate: func(value string) (err error) {
				formats := getValidTimeFormats(TimeFormatAny, true)
				if !slices.Contains(formats, value) {
					return fmt.Errorf("%v is not a valid value for TIME_INPUT_FORMAT, must be one of %v", value, formats)
				}
				return nil
			},
		},
		"TIME_OUTPUT_FORMAT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "HH24:MI:SS",
			Validate: func(value string) (err error) {
				formats := getValidTimeFormats(TimeFormatAny, false)
				if !slices.Contains(formats, value) {
					return fmt.Errorf("%v is not a valid value for TIME_OUTPUT_FORMAT, must be one of %v", value, formats)
				}
				return nil
			},
		},
		"TRANSACTION_DEFAULT_ISOLATION_LEVEL": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "READ_COMMITTED",
			Validate: func(value string) (err error) {
				if !slices.Contains([]string{"READ_COMMITTED"}, value) {
					return fmt.Errorf("%v is not a valid value for TRANSACTION_DEFAULT_ISOLATION_LEVEL, must be one of READ_UNCOMMITTED, READ_COMMITTED, REPEATABLE_READ, SERIALIZABLE", value)
				}
				return nil
			},
		},
		"TWO_DIGIT_CENTURY_START": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: 1970,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for TWO_DIGIT_CENTURY_START, must be an integer", value)
				}
				if v < 1900 || v > 2100 {
					return fmt.Errorf("%v is not a valid value for TWO_DIGIT_CENTURY_START, must be between 1900 and 2100", value)
				}
				return nil
			},
		},
		"UNSUPPORTED_DDL_ACTION": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: "IGNORE",
			Validate: func(value string) (err error) {
				if !slices.Contains([]string{"IGNORE", "FAIL"}, value) {
					return fmt.Errorf("%v is not a valid value for UNSUPPORTED_DDL_ACTION, must be one of IGNORE, FAIL", value)
				}
				return nil
			},
		},
		"USE_CACHED_RESULT": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: true,
			Validate:     validateBoolFunc,
		},
		"WEEK_OF_YEAR_POLICY": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: 0,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for WEEK_OF_YEAR_POLICY, must be an integer", value)
				}
				if v < 0 || v > 1 {
					return fmt.Errorf("%v is not a valid value for WEEK_OF_YEAR_POLICY, must be 0 or 1", value)
				}
				return nil
			},
		},
		"WEEK_START": {
			TypeSet:      []ParameterType{ParameterTypeSession},
			DefaultValue: 0,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for WEEK_START, must be an integer", value)
				}
				if v < 0 || v > 7 {
					return fmt.Errorf("%v is not a valid value for WEEK_START, must be between 0 and 7", value)
				}
				return nil
			},
		},
		"DATA_RETENTION_TIME_IN_DAYS": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: 1,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for DATA_RETENTION_TIME_IN_DAYS, must be an integer", value)
				}
				if v < 0 || v > 90 {
					return fmt.Errorf("%v is not a valid value for DATA_RETENTION_TIME_IN_DAYS, must be between 0 and 90", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeDatabase,
				ObjectTypeSchema,
				ObjectTypeTable,
			},
		},
		"DEFAULT_DDL_COLLATION": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: "",
			Validate: func(value string) (err error) {
				// todo: validate collation.
				if len(value) < 1 {
					return fmt.Errorf("%v is not a valid value for DEFAULT_DDL_COLLATION, must be a valid collation", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeDatabase,
				ObjectTypeSchema,
				ObjectTypeTable,
			},
		},
		"ENABLE_STREAM_TASK_REPLICATION": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: false,
			Validate:     validateBoolFunc,
			AllowedObjectTypes: []ObjectType{
				ObjectTypeDatabase,
				ObjectTypeReplicationGroup,
				ObjectTypeFailoverGroup,
			},
		},
		"MAX_CONCURRENCY_LEVEL": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: 0,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for MAX_CONCURRENCY_LEVEL, must be an integer", value)
				}
				if v < 0 {
					return fmt.Errorf("%v is not a valid value for MAX_CONCURRENCY_LEVEL, must be an integer greater than 0", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeWarehouse,
			},
		},
		"MAX_DATA_EXTENSION_TIME_IN_DAYS": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: 14,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for MAX_DATA_EXTENSION_TIME_IN_DAYS, must be an integer", value)
				}
				if v < 0 || v > 90 {
					return fmt.Errorf("%v is not a valid value for MAX_DATA_EXTENSION_TIME_IN_DAYS, must be between 0 and 90", value)
				}
				return nil
			},
		},
		"PIPE_EXECUTION_PAUSED": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: false,
			Validate:     validateBoolFunc,
			AllowedObjectTypes: []ObjectType{
				ObjectTypeSchema,
				ObjectTypePipe,
			},
		},
		"PREVENT_UNLOAD_TO_INTERNAL_STAGES": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: false,
			Validate:     validateBoolFunc,
			AllowedObjectTypes: []ObjectType{
				ObjectTypeUser,
			},
		},
		"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: 0,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for STATEMENT_QUEUED_TIMEOUT_IN_SECONDS, must be an integer", value)
				}
				if v < 0 {
					return fmt.Errorf("%v is not a valid value for STATEMENT_QUEUED_TIMEOUT_IN_SECONDS, must be an integer greater than 0", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeWarehouse,
			},
		},
		"SHARE_RESTRICTIONS": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: true,
			Validate:     validateBoolFunc,
			AllowedObjectTypes: []ObjectType{
				ObjectTypeShare,
			},
		},
		"SUSPEND_TASK_AFTER_NUM_FAILURES": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: 0,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for SUSPEND_TASK_AFTER_NUM_FAILURES, must be an integer", value)
				}
				if v < 0 {
					return fmt.Errorf("%v is not a valid value for SUSPEND_TASK_AFTER_NUM_FAILURES, must be an integer greater than 0", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeDatabase,
				ObjectTypeSchema,
				ObjectTypeTask,
			},
		},
		"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: "MEDIUM",
			Validate: func(value string) (err error) {
				if !slices.Contains([]string{"X-SMALL", "SMALL", "MEDIUM", "LARGE", "X-LARGE", "2X-LARGE", "3X-LARGE", "4X-LARGE", "5X-LARGE", "6X-LARGE"}, value) {
					return fmt.Errorf("%v is not a valid value for USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE, must be a valid warehouse size, such as \"SMALL\", \"MEDIUM\" or \"LARGE\"", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeDatabase,
				ObjectTypeSchema,
				ObjectTypeTask,
			},
		},
		"USER_TASK_TIMEOUT_MS": {
			TypeSet:      []ParameterType{ParameterTypeObject},
			DefaultValue: 3600000,
			Validate: func(value string) (err error) {
				v, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("%v is not a valid value for USER_TASK_TIMEOUT_MS, must be an integer", value)
				}
				if v < 0 || v > 86400000 {
					return fmt.Errorf("%v is not a valid value for USER_TASK_TIMEOUT_MS, must be an integer greater than 0 and less than 86400000 (1 day)", value)
				}
				return nil
			},
			AllowedObjectTypes: []ObjectType{
				ObjectTypeDatabase,
				ObjectTypeSchema,
				ObjectTypeTask,
			},
		},
	}
}

// GetParameterObjectTypeSetAsStrings returns a slice of all object types that can have parameters.
func GetParameterObjectTypeSetAsStrings() []string {
	objectTypeSet := []ObjectType{
		ObjectTypeDatabase,
		ObjectTypeSchema,
		ObjectTypePipe,
		ObjectTypeUser,
		ObjectTypeShare,
		ObjectTypeWarehouse,
		ObjectTypeTask,
		ObjectTypeReplicationGroup,
		ObjectTypeFailoverGroup,
		ObjectTypeTable,
	}
	result := make([]string, 0, len(objectTypeSet))
	for _, v := range objectTypeSet {
		result = append(result, string(v))
	}
	return result
}

// GetParameters returns a map of parameters that match the given type (e.g. Account, Session, Object).
func GetParameterDefaults(t ParameterType) map[string]ParameterDefault {
	parameters := ParameterDefaults()
	keys := maps.Keys(parameters)
	for _, key := range keys {
		typeSet := parameters[key].TypeSet
		if !slices.Contains(typeSet, t) {
			delete(parameters, key)
		}
	}
	return parameters
}

// GetParameter returns a parameter by key.
func GetParameterDefault(key string) ParameterDefault {
	return ParameterDefaults()[key]
}

// ParameterBuilder abstracts the creation of SQL queries for Snowflake parameters.
type ParameterBuilder struct {
	key           string
	value         string
	parameterType ParameterType
	objectType    ObjectType
	objectIdentifier    string
	db            *sql.DB
}

func NewParameter(key, value string, parameterType ParameterType, db *sql.DB) *ParameterBuilder {
	return &ParameterBuilder{
		key:           key,
		value:         value,
		parameterType: parameterType,
		db:            db,
	}
}

func (v *ParameterBuilder) WithObjectType(objectType ObjectType) *ParameterBuilder {
	v.objectType = objectType
	return v
}

func (v *ParameterBuilder) WithObjectIdentifier(objectIdentifier string) *ParameterBuilder {
	v.objectIdentifier = objectIdentifier
	return v
}

func (v *ParameterBuilder) SetParameter() error {
	if v.parameterType == ParameterTypeAccount || v.parameterType == ParameterTypeSession {
		// prepared statements do not work here for some reason. We already validate inputs so its okay
		stmt := fmt.Sprintf("ALTER ACCOUNT SET %s = %s", v.key, v.value)
		_, err := v.db.Exec(stmt)
		if err != nil {
			return err
		}
	} else if v.parameterType == ParameterTypeObject {
		stmt := fmt.Sprintf("ALTER %s %s SET %s = %s", v.objectType, v.objectIdentifier, v.key, v.value)
		_, err := v.db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

type Parameter struct {
	Key         sql.NullString `db:"key"`
	Value       sql.NullString `db:"value"`
	Default     sql.NullString `db:"default"`
	Level       sql.NullString `db:"level"`
	Description sql.NullString `db:"description"`
	PType       sql.NullString `db:"type"`
}

func ShowParameter(db *sql.DB, key string, parameterType ParameterType) (*Parameter, error) {
	var value Parameter
	var stmt string
	if parameterType == ParameterTypeAccount || parameterType == ParameterTypeSession {
		stmt = fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN ACCOUNT", key)
	} else {
		return nil, fmt.Errorf("unsupported parameter type %s", parameterType)
	}
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	params := []Parameter{}
	if err := sqlx.StructScan(rows, &params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	value = params[0]

	return &value, nil
}

func ShowObjectParameter(db *sql.DB, key string, objectType ObjectType, objectIdentifier string) (*Parameter, error) {
	var value Parameter
	stmt := fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN %s %s", key, objectType.String(), objectIdentifier)
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	params := []Parameter{}
	if err := sqlx.StructScan(rows, &params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	value = params[0]

	return &value, nil
}

func ListParameters(db *sql.DB, parameterType ParameterType, pattern string) ([]Parameter, error) {
	var stmt string
	if parameterType == ParameterTypeAccount || parameterType == ParameterTypeSession {
		if pattern != "" {
			stmt = fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN %v", pattern, parameterType)
		} else {
			stmt = fmt.Sprintf("SHOW PARAMETERS IN %v", parameterType)
		}
	} else {
		return nil, fmt.Errorf("unsupported parameter type %s", parameterType)
	}
	log.Printf("[DEBUG] query = %s", stmt)
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	params := []Parameter{}
	if err := sqlx.StructScan(rows, &params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return params, nil
}

func ListObjectParameters(db *sql.DB, objectType ObjectType, objectIdentifier, pattern string) ([]Parameter, error) {
	var stmt string
	if pattern != "" {
		stmt = fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN %s %s", pattern, objectType.String(), objectIdentifier)
	} else {
		stmt = fmt.Sprintf("SHOW PARAMETERS IN %s %s", objectType.String(), objectIdentifier)
	}
	log.Printf("[DEBUG] query = %s", stmt)
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	params := []Parameter{}
	if err := sqlx.StructScan(rows, &params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return params, nil
}
