package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/jmoiron/sqlx"
)

// ParameterType is the type of parameter.
type ParameterType string

const (
	ParameterTypeAccount ParameterType = "ACCOUNT"
	ParameterTypeSession ParameterType = "SESSION"
	ParameterTypeObject  ParameterType = "OBJECT"
)

type ParameterExecutor struct {
	db *sql.DB
}

func NewParameterExecutor(db *sql.DB) *ParameterExecutor {
	return &ParameterExecutor{
		db: db,
	}
}

func (v *ParameterExecutor) Execute(stmt string, args ...interface{}) error {
	_, err := v.db.Exec(stmt, args...)
	return err
}

func (v *ParameterExecutor) Query(stmt string) ([]Parameter, error) {
	rows, err := v.db.Query(stmt)
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

func (v *ParameterExecutor) QueryOne(stmt string) (*Parameter, error) {
	params, err := v.Query(stmt)
	if err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, nil
	}
	return &params[0], nil
}

// AccountParameterBuilder abstracts the creation of SQL queries for Snowflake account parameters.
type AccountParameterBuilder struct {
	key      string
	value    string
	executor *ParameterExecutor
}

func NewAccountParameter(key, value string, db *sql.DB) *AccountParameterBuilder {
	return &AccountParameterBuilder{
		key:      key,
		value:    value,
		executor: NewParameterExecutor(db),
	}
}

func (v *AccountParameterBuilder) SetParameter() error {
	stmt := fmt.Sprintf("ALTER ACCOUNT SET %s = %s", v.key, v.value)
	return v.executor.Execute(stmt)
}

// SessionParameterBuilder abstracts the creation of SQL queries for Snowflake session parameters.
type SessionParameterBuilder struct {
	key       string
	value     string
	onAccount bool
	user      string
	executor  *ParameterExecutor
}

func NewSessionParameter(key, value string, db *sql.DB) *SessionParameterBuilder {
	return &SessionParameterBuilder{
		key:      key,
		value:    value,
		executor: NewParameterExecutor(db),
	}
}

func (v *SessionParameterBuilder) SetOnAccount(onAccount bool) *SessionParameterBuilder {
	v.onAccount = onAccount
	return v
}

func (v *SessionParameterBuilder) SetUser(user string) *SessionParameterBuilder {
	v.user = user
	return v
}

func (v *SessionParameterBuilder) SetParameter() error {
	if v.onAccount {
		stmt := fmt.Sprintf("ALTER ACCOUNT SET %s = %s", v.key, v.value)
		return v.executor.Execute(stmt)
	}
	if v.user == "" {
		return fmt.Errorf("user is required when setting session parameters on a user")
	}
	stmt := fmt.Sprintf("ALTER USER %s SET %s = %s", v.user, v.key, v.value)
	return v.executor.Execute(stmt)
}

// ObjectParameterBuilder abstracts the creation of SQL queries for Snowflake object parameters.
type ObjectParameterBuilder struct {
	key              string
	value            string
	onAccount        bool
	objectType       sdk.ObjectType
	objectIdentifier string
	executor         *ParameterExecutor
}

func NewObjectParameter(key, value string, db *sql.DB) *ObjectParameterBuilder {
	return &ObjectParameterBuilder{
		key:      key,
		value:    value,
		executor: NewParameterExecutor(db),
	}
}

func (v *ObjectParameterBuilder) SetOnAccount(onAccount bool) *ObjectParameterBuilder {
	v.onAccount = onAccount
	return v
}

func (v *ObjectParameterBuilder) WithObjectType(objectType sdk.ObjectType) *ObjectParameterBuilder {
	v.objectType = objectType
	return v
}

func (v *ObjectParameterBuilder) WithObjectIdentifier(objectIdentifier string) *ObjectParameterBuilder {
	v.objectIdentifier = objectIdentifier
	return v
}

func (v *ObjectParameterBuilder) SetParameter() error {
	if v.onAccount {
		stmt := fmt.Sprintf("ALTER ACCOUNT SET %s = %s", v.key, v.value)
		return v.executor.Execute(stmt)
	}
	if v.objectType == "" {
		return fmt.Errorf("object type is required when setting object parameters")
	}
	if v.objectIdentifier == "" {
		return fmt.Errorf("object identifier is required when setting object parameters")
	}

	stmt := fmt.Sprintf("ALTER %s %s SET %s = %s", v.objectType, v.objectIdentifier, v.key, v.value)
	return v.executor.Execute(stmt)
}

type Parameter struct {
	Key         sql.NullString `db:"key"`
	Value       sql.NullString `db:"value"`
	Default     sql.NullString `db:"default"`
	Level       sql.NullString `db:"level"`
	Description sql.NullString `db:"description"`
	PType       sql.NullString `db:"type"`
}

func ShowAccountParameter(db *sql.DB, key string) (*Parameter, error) {
	stmt := fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN ACCOUNT", key)
	executor := NewParameterExecutor(db)
	params, err := executor.Query(stmt)
	if err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, nil
	}
	return &params[0], nil
}

func ShowSessionParameter(db *sql.DB, key string, user string) (*Parameter, error) {
	stmt := fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN USER %s", key, user)
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

	return &params[0], nil
}

func ShowObjectParameter(db *sql.DB, key string, objectType sdk.ObjectType, objectIdentifier string) (*Parameter, error) {
	stmt := fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN %s %s", key, objectType.String(), objectIdentifier)
	executor := NewParameterExecutor(db)
	return executor.QueryOne(stmt)
}

func ListAccountParameters(db *sql.DB, pattern string) ([]Parameter, error) {
	var stmt string
	if pattern != "" {
		stmt = fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN ACCOUNT", pattern)
	} else {
		stmt = "SHOW PARAMETERS IN ACCOUNT"
	}
	executor := NewParameterExecutor(db)
	return executor.Query(stmt)
}

func ListSessionParameters(db *sql.DB, pattern string, user string) ([]Parameter, error) {
	var stmt string
	if pattern != "" {
		stmt = fmt.Sprintf("SHOW PARAMETERS LIKE '%s' FOR USER %s", pattern, user)
	} else {
		stmt = fmt.Sprintf("SHOW PARAMETERS FOR USER %s", user)
	}
	executor := NewParameterExecutor(db)
	return executor.Query(stmt)
}

func ListObjectParameters(db *sql.DB, objectType sdk.ObjectType, objectIdentifier, pattern string) ([]Parameter, error) {
	var stmt string
	if pattern != "" {
		stmt = fmt.Sprintf("SHOW PARAMETERS LIKE '%s' IN %s %s", pattern, objectType.String(), objectIdentifier)
	} else {
		stmt = fmt.Sprintf("SHOW PARAMETERS IN %s %s", objectType.String(), objectIdentifier)
	}
	executor := NewParameterExecutor(db)
	return executor.Query(stmt)
}
