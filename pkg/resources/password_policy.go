package resources

import (
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var passwordPolicySchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database this password policy belongs to.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema this password policy belongs to.",
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Identifier for the password policy; must be unique for your account.",
	},
	"or_replace": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether to override a previous password policy with the same name.",
	},
	"if_not_exists": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Prevent overwriting a previous password policy with the same name.",
	},
	"min_length": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      8,
		Description:  "Specifies the minimum number of characters the password must contain. Supported range: 8 to 256, inclusive. Default: 8",
		ValidateFunc: validation.IntBetween(8, 256),
	},
	"max_length": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      256,
		Description:  "Specifies the maximum number of characters the password must contain. This number must be greater than or equal to the sum of PASSWORD_MIN_LENGTH, PASSWORD_MIN_UPPER_CASE_CHARS, and PASSWORD_MIN_LOWER_CASE_CHARS. Supported range: 8 to 256, inclusive. Default: 256",
		ValidateFunc: validation.IntBetween(8, 256),
	},
	"min_upper_case_chars": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the minimum number of uppercase characters the password must contain. Supported range: 0 to 256, inclusive. Default: 1",
		ValidateFunc: validation.IntBetween(0, 256),
	},
	"min_lower_case_chars": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the minimum number of lowercase characters the password must contain. Supported range: 0 to 256, inclusive. Default: 1",
		ValidateFunc: validation.IntBetween(0, 256),
	},
	"min_numeric_chars": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the minimum number of numeric characters the password must contain. Supported range: 0 to 256, inclusive. Default: 1",
		ValidateFunc: validation.IntBetween(0, 256),
	},
	"min_special_chars": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the minimum number of special characters the password must contain. Supported range: 0 to 256, inclusive. Default: 1",
		ValidateFunc: validation.IntBetween(0, 256),
	},
	"max_age_days": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      90,
		Description:  "Specifies the maximum number of days before the password must be changed. Supported range: 0 to 999, inclusive. A value of zero (i.e. 0) indicates that the password does not need to be changed. Snowflake does not recommend choosing this value for a default account-level password policy or for any user-level policy. Instead, choose a value that meets your internal security guidelines. Default: 90, which means the password must be changed every 90 days.",
		ValidateFunc: validation.IntBetween(0, 999),
	},
	"max_retries": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      5,
		Description:  "Specifies the maximum number of attempts to enter a password before being locked out. Supported range: 1 to 10, inclusive. Default: 5",
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"lockout_time_mins": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      15,
		Description:  "Specifies the number of minutes the user account will be locked after exhausting the designated number of password retries (i.e. PASSWORD_MAX_RETRIES). Supported range: 1 to 999, inclusive. Default: 15",
		ValidateFunc: validation.IntBetween(1, 999),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Adds a comment or overwrites an existing comment for the password policy.",
	},
}

func PasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "A password policy specifies the requirements that must be met to create and reset a password to authenticate to Snowflake.",
		Create:      CreatePasswordPolicy,
		Read:        ReadPasswordPolicy,
		Update:      UpdatePasswordPolicy,
		Delete:      DeletePasswordPolicy,

		Schema: passwordPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreatePasswordPolicy implements schema.CreateFunc.
func CreatePasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	builder, err := snowflake.NewPasswordPolicyBuilder()
	if err != nil {
		return fmt.Errorf("couldn't create password policy builder: %w", err)
	}

	props := &snowflake.PasswordPolicyProps{
		Database: d.Get("database").(string),
		Schema:   d.Get("schema").(string),
		Name:     d.Get("name").(string),

		OrReplace:           d.Get("or_replace").(bool),
		OrReplaceOk:         builder.Ok(d.GetOk("or_replace")),
		IfNotExists:         d.Get("if_not_exists").(bool),
		IfNotExistsOk:       builder.Ok(d.GetOk("if_not_exists")),
		MinLength:           d.Get("min_length").(int),
		MinLengthOk:         builder.Ok(d.GetOk("min_length")),
		MaxLength:           d.Get("max_length").(int),
		MaxLengthOk:         builder.Ok(d.GetOk("max_length")),
		MinUpperCaseChars:   d.Get("min_upper_case_chars").(int),
		MinUpperCaseCharsOk: builder.Ok(d.GetOk("min_upper_case_chars")),
		MinLowerCaseChars:   d.Get("min_lower_case_chars").(int),
		MinLowerCaseCharsOk: builder.Ok(d.GetOk("min_lower_case_chars")),
		MinNumericChars:     d.Get("min_numeric_chars").(int),
		MinNumericCharsOk:   builder.Ok(d.GetOk("min_numeric_chars")),
		MinSpecialChars:     d.Get("min_special_chars").(int),
		MinSpecialCharsOk:   builder.Ok(d.GetOk("min_special_chars")),
		MaxAgeDays:          d.Get("max_age_days").(int),
		MaxAgeDaysOk:        builder.Ok(d.GetOk("max_age_days")),
		MaxRetries:          d.Get("max_retries").(int),
		MaxRetriesOk:        builder.Ok(d.GetOk("max_retries")),
		LockoutTimeMins:     d.Get("lockout_time_mins").(int),
		LockoutTimeMinsOk:   builder.Ok(d.GetOk("lockout_time_mins")),
		Comment:             d.Get("comment").(string),
		CommentOk:           builder.Ok(d.GetOk("comment")),
	}

	stmt, err := builder.Create(props)
	if err != nil {
		return fmt.Errorf("couldn't generate create statement: %w", err)
	}

	db := meta.(*sql.DB)
	_, err = db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error executing create statement: %w", err)
	}

	d.SetId(props.Id())

	return nil
}

// ReadPasswordPolicy implements schema.ReadFunc.
func ReadPasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	builder, err := snowflake.NewPasswordPolicyBuilder()
	if err != nil {
		return fmt.Errorf("couldn't create password policy builder: %w", err)
	}

	props := &snowflake.PasswordPolicyProps{
		Database: d.Get("database").(string),
		Schema:   d.Get("schema").(string),
		Name:     d.Get("name").(string),
	}

	stmt, err := builder.Describe(props)
	if err != nil {
		return fmt.Errorf("couldn't generate describe statement: %w", err)
	}

	db := meta.(*sql.DB)
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("error querying password policy: %w", err)
	}

	defer rows.Close()
	err = builder.ParseDescribe(rows, props)
	if err != nil {
		return fmt.Errorf("failed to parse result of describe: %w", err)
	}

	err = d.Set("min_length", props.MinLength)
	if err != nil {
		return fmt.Errorf("error setting min_length: %w", err)
	}
	err = d.Set("max_length", props.MaxLength)
	if err != nil {
		return fmt.Errorf("error setting max_length: %w", err)
	}
	err = d.Set("min_upper_case_chars", props.MinUpperCaseChars)
	if err != nil {
		return fmt.Errorf("error setting min_upper_case_chars: %w", err)
	}
	err = d.Set("min_lower_case_chars", props.MinLowerCaseChars)
	if err != nil {
		return fmt.Errorf("error setting min_lower_case_chars: %w", err)
	}
	err = d.Set("min_numeric_chars", props.MinNumericChars)
	if err != nil {
		return fmt.Errorf("error setting min_numeric_chars: %w", err)
	}
	err = d.Set("min_special_chars", props.MinSpecialChars)
	if err != nil {
		return fmt.Errorf("error setting min_special_chars: %w", err)
	}
	err = d.Set("max_age_days", props.MaxAgeDays)
	if err != nil {
		return fmt.Errorf("error setting max_age_days: %w", err)
	}
	err = d.Set("max_retries", props.MaxRetries)
	if err != nil {
		return fmt.Errorf("error setting max_retries: %w", err)
	}
	err = d.Set("lockout_time_mins", props.LockoutTimeMins)
	if err != nil {
		return fmt.Errorf("error setting lockout_time_mins: %w", err)
	}

	return nil
}

// UpdatePasswordPolicy implements schema.UpdateFunc.
func UpdatePasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	builder, err := snowflake.NewPasswordPolicyBuilder()
	if err != nil {
		return fmt.Errorf("couldn't create password policy builder: %w", err)
	}

	props := &snowflake.PasswordPolicyProps{
		Database: d.Get("database").(string),
		Schema:   d.Get("schema").(string),
		Name:     d.Get("name").(string),
	}

	if d.HasChange("min_length") {
		props.MinLength = d.Get("min_length").(int)
		props.MinLengthOk = true
	}
	if d.HasChange("max_length") {
		props.MaxLength = d.Get("max_length").(int)
		props.MaxLengthOk = true
	}
	if d.HasChange("min_upper_case_chars") {
		props.MinUpperCaseChars = d.Get("min_upper_case_chars").(int)
		props.MinUpperCaseCharsOk = true
	}
	if d.HasChange("min_lower_case_chars") {
		props.MinLowerCaseChars = d.Get("min_lower_case_chars").(int)
		props.MinLowerCaseCharsOk = true
	}
	if d.HasChange("min_numeric_chars") {
		props.MinNumericChars = d.Get("min_numeric_chars").(int)
		props.MinNumericCharsOk = true
	}
	if d.HasChange("min_special_chars") {
		props.MinSpecialChars = d.Get("min_special_chars").(int)
		props.MinSpecialCharsOk = true
	}
	if d.HasChange("max_age_days") {
		props.MaxAgeDays = d.Get("max_age_days").(int)
		props.MaxAgeDaysOk = true
	}
	if d.HasChange("max_retries") {
		props.MaxRetries = d.Get("max_retries").(int)
		props.MaxRetriesOk = true
	}
	if d.HasChange("lockout_time_mins") {
		props.LockoutTimeMins = d.Get("lockout_time_mins").(int)
		props.LockoutTimeMinsOk = true
	}
	if d.HasChange("comment") {
		props.Comment = d.Get("comment").(string)
		props.CommentOk = true
	}

	stmt, err := builder.Alter(props)
	if err != nil {
		return fmt.Errorf("couldn't generate alter statement for password policy: %w", err)
	}

	db := meta.(*sql.DB)
	_, err = db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error executing alter statement: %w", err)
	}

	return nil
}

// DeletePasswordPolicy implements schema.DeleteFunc.
func DeletePasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	builder, err := snowflake.NewPasswordPolicyBuilder()
	if err != nil {
		return fmt.Errorf("couldn't create password policy builder: %w", err)
	}

	props := &snowflake.PasswordPolicyProps{
		Database: d.Get("database").(string),
		Schema:   d.Get("schema").(string),
		Name:     d.Get("name").(string),
	}

	stmt, err := builder.Drop(props)
	if err != nil {
		return fmt.Errorf("couldn't generate drop statement: %w", err)
	}

	db := meta.(*sql.DB)
	_, err = db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error executing drop statement: %w", err)
	}

	return nil
}
