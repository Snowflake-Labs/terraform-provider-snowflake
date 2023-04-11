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
	manager, err := snowflake.NewPasswordPolicyManager()
	if err != nil {
		return fmt.Errorf("couldn't create password policy builder: %w", err)
	}

	input := &snowflake.PasswordPolicyCreateInput{
		PasswordPolicy: snowflake.PasswordPolicy{
			SchemaObjectIdentifier: snowflake.SchemaObjectIdentifier{
				Database:   d.Get("database").(string),
				Schema:     d.Get("schema").(string),
				ObjectName: d.Get("name").(string),
			},

			MinLength:           d.Get("min_length").(int),
			MinLengthOk:         manager.Ok(d.GetOk("min_length")),
			MaxLength:           d.Get("max_length").(int),
			MaxLengthOk:         manager.Ok(d.GetOk("max_length")),
			MinUpperCaseChars:   d.Get("min_upper_case_chars").(int),
			MinUpperCaseCharsOk: manager.Ok(d.GetOk("min_upper_case_chars")),
			MinLowerCaseChars:   d.Get("min_lower_case_chars").(int),
			MinLowerCaseCharsOk: manager.Ok(d.GetOk("min_lower_case_chars")),
			MinNumericChars:     d.Get("min_numeric_chars").(int),
			MinNumericCharsOk:   manager.Ok(d.GetOk("min_numeric_chars")),
			MinSpecialChars:     d.Get("min_special_chars").(int),
			MinSpecialCharsOk:   manager.Ok(d.GetOk("min_special_chars")),
			MaxAgeDays:          d.Get("max_age_days").(int),
			MaxAgeDaysOk:        manager.Ok(d.GetOk("max_age_days")),
			MaxRetries:          d.Get("max_retries").(int),
			MaxRetriesOk:        manager.Ok(d.GetOk("max_retries")),
			LockoutTimeMins:     d.Get("lockout_time_mins").(int),
			LockoutTimeMinsOk:   manager.Ok(d.GetOk("lockout_time_mins")),
			Comment:             d.Get("comment").(string),
			CommentOk:           manager.Ok(d.GetOk("comment")),
		},

		OrReplace:     d.Get("or_replace").(bool),
		OrReplaceOk:   manager.Ok(d.GetOk("or_replace")),
		IfNotExists:   d.Get("if_not_exists").(bool),
		IfNotExistsOk: manager.Ok(d.GetOk("if_not_exists")),
	}

	stmt, err := manager.Create(input)
	if err != nil {
		return fmt.Errorf("couldn't generate create statement: %w", err)
	}

	db := meta.(*sql.DB)
	_, err = db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error executing create statement: %w", err)
	}

	d.SetId(PasswordPolicyID(&input.PasswordPolicy))

	return nil
}

// ReadPasswordPolicy implements schema.ReadFunc.
func ReadPasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	manager, err := snowflake.NewPasswordPolicyManager()
	if err != nil {
		return fmt.Errorf("couldn't create password policy builder: %w", err)
	}

	input := &snowflake.PasswordPolicyReadInput{
		Database:   d.Get("database").(string),
		Schema:     d.Get("schema").(string),
		ObjectName: d.Get("name").(string),
	}

	stmt, err := manager.Read(input)
	if err != nil {
		return fmt.Errorf("couldn't generate describe statement: %w", err)
	}

	db := meta.(*sql.DB)
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("error querying password policy: %w", err)
	}

	defer rows.Close()
	output, err := manager.Parse(rows)
	if err != nil {
		return fmt.Errorf("failed to parse result of describe: %w", err)
	}

	if err = d.Set("min_length", output.MinLength); err != nil {
		return fmt.Errorf("error setting min_length: %w", err)
	}
	if err = d.Set("max_length", output.MaxLength); err != nil {
		return fmt.Errorf("error setting max_length: %w", err)
	}
	if err = d.Set("min_upper_case_chars", output.MinUpperCaseChars); err != nil {
		return fmt.Errorf("error setting min_upper_case_chars: %w", err)
	}
	if err = d.Set("min_lower_case_chars", output.MinLowerCaseChars); err != nil {
		return fmt.Errorf("error setting min_lower_case_chars: %w", err)
	}
	if err = d.Set("min_numeric_chars", output.MinNumericChars); err != nil {
		return fmt.Errorf("error setting min_numeric_chars: %w", err)
	}
	if err = d.Set("min_special_chars", output.MinSpecialChars); err != nil {
		return fmt.Errorf("error setting min_special_chars: %w", err)
	}
	if err = d.Set("max_age_days", output.MaxAgeDays); err != nil {
		return fmt.Errorf("error setting max_age_days: %w", err)
	}
	if err = d.Set("max_retries", output.MaxRetries); err != nil {
		return fmt.Errorf("error setting max_retries: %w", err)
	}
	if err = d.Set("lockout_time_mins", output.LockoutTimeMins); err != nil {
		return fmt.Errorf("error setting lockout_time_mins: %w", err)
	}

	return nil
}

// UpdatePasswordPolicy implements schema.UpdateFunc.
func UpdatePasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	manager, err := snowflake.NewPasswordPolicyManager()
	if err != nil {
		return fmt.Errorf("couldn't create password policy builder: %w", err)
	}

	runAlter := false
	alterInput := &snowflake.PasswordPolicyUpdateInput{
		PasswordPolicy: snowflake.PasswordPolicy{
			SchemaObjectIdentifier: snowflake.SchemaObjectIdentifier{
				Database:   d.Get("database").(string),
				Schema:     d.Get("schema").(string),
				ObjectName: d.Get("name").(string),
			},
		},
	}
	runUnset := false
	unsetInput := &snowflake.PasswordPolicyUpdateInput{
		PasswordPolicy: snowflake.PasswordPolicy{
			SchemaObjectIdentifier: snowflake.SchemaObjectIdentifier{
				Database:   d.Get("database").(string),
				Schema:     d.Get("schema").(string),
				ObjectName: d.Get("name").(string),
			},
		},
	}

	if d.HasChange("min_length") {
		val, ok := d.GetOk("min_length")
		if ok {
			alterInput.MinLength = val.(int)
			alterInput.MinLengthOk = true
			runAlter = true
		} else {
			unsetInput.MinLengthOk = true
			runUnset = true
		}
	}
	if d.HasChange("max_length") {
		val, ok := d.GetOk("max_length")
		if ok {
			alterInput.MaxLength = val.(int)
			alterInput.MaxLengthOk = true
			runAlter = true
		} else {
			unsetInput.MaxLengthOk = true
			runUnset = true
		}
	}
	if d.HasChange("min_upper_case_chars") {
		val, ok := d.GetOk("min_upper_case_chars")
		if ok {
			alterInput.MinUpperCaseChars = val.(int)
			alterInput.MinUpperCaseCharsOk = true
			runAlter = true
		} else {
			unsetInput.MinUpperCaseCharsOk = true
			runUnset = true
		}
	}
	if d.HasChange("min_lower_case_chars") {
		val, ok := d.GetOk("min_lower_case_chars")
		if ok {
			alterInput.MinLowerCaseChars = val.(int)
			alterInput.MinLowerCaseCharsOk = true
			runAlter = true
		} else {
			unsetInput.MinLowerCaseCharsOk = true
			runUnset = true
		}
	}
	if d.HasChange("min_numeric_chars") {
		val, ok := d.GetOk("min_numeric_chars")
		if ok {
			alterInput.MinNumericChars = val.(int)
			alterInput.MinNumericCharsOk = true
			runAlter = true
		} else {
			unsetInput.MinNumericCharsOk = true
			runUnset = true
		}
	}
	if d.HasChange("min_special_chars") {
		val, ok := d.GetOk("min_special_chars")
		if ok {
			alterInput.MinSpecialChars = val.(int)
			alterInput.MinSpecialCharsOk = true
			runAlter = true
		} else {
			unsetInput.MinSpecialCharsOk = true
			runUnset = true
		}
	}
	if d.HasChange("max_age_days") {
		val, ok := d.GetOk("max_age_days")
		if ok {
			alterInput.MaxAgeDays = val.(int)
			alterInput.MaxAgeDaysOk = true
			runAlter = true
		} else {
			unsetInput.MaxAgeDaysOk = true
			runUnset = true
		}
	}
	if d.HasChange("max_retries") {
		val, ok := d.GetOk("max_retries")
		if ok {
			alterInput.MaxRetries = val.(int)
			alterInput.MaxRetriesOk = true
			runAlter = true
		} else {
			unsetInput.MaxRetriesOk = true
			runUnset = true
		}
	}
	if d.HasChange("lockout_time_mins") {
		val, ok := d.GetOk("lockout_time_mins")
		if ok {
			alterInput.LockoutTimeMins = val.(int)
			alterInput.LockoutTimeMinsOk = true
			runAlter = true
		} else {
			unsetInput.LockoutTimeMinsOk = true
			runUnset = true
		}
	}
	if d.HasChange("comment") {
		val, ok := d.GetOk("comment")
		if ok {
			alterInput.Comment = val.(string)
			alterInput.CommentOk = true
			runAlter = true
		} else {
			unsetInput.CommentOk = true
			runUnset = true
		}
	}

	db := meta.(*sql.DB)

	if runAlter {
		stmt, err := manager.Update(alterInput)
		if err != nil {
			return fmt.Errorf("couldn't generate alter statement for password policy: %w", err)
		}

		_, err = db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing alter statement: %w", err)
		}
	}

	if runUnset {
		stmt, err := manager.Unset(unsetInput)
		if err != nil {
			return fmt.Errorf("couldn't generate unset statement for password policy: %w", err)
		}

		_, err = db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing unset statement: %w", err)
		}
	}

	return nil
}

// DeletePasswordPolicy implements schema.DeleteFunc.
func DeletePasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	manager, err := snowflake.NewPasswordPolicyManager()
	if err != nil {
		return fmt.Errorf("couldn't create password policy builder: %w", err)
	}

	input := &snowflake.PasswordPolicyDeleteInput{
		SchemaObjectIdentifier: snowflake.SchemaObjectIdentifier{
			Database:   d.Get("database").(string),
			Schema:     d.Get("schema").(string),
			ObjectName: d.Get("name").(string),
		},
	}

	stmt, err := manager.Delete(input)
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

func PasswordPolicyID(pp *snowflake.PasswordPolicy) string {
	return pp.QualifiedName()
}
