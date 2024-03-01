package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var emailNotificationIntegrationSchema = map[string]*schema.Schema{
	// The first part of the schema is shared between all integration vendors
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
	},
	"allowed_recipients": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "List of email addresses that should receive notifications.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A comment for the email integration.",
	},
}

// EmailNotificationIntegration returns a pointer to the resource representing a notification integration.
func EmailNotificationIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateEmailNotificationIntegration,
		Read:   ReadEmailNotificationIntegration,
		Update: UpdateEmailNotificationIntegration,
		Delete: DeleteEmailNotificationIntegration,

		Schema: emailNotificationIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func toAllowedRecipients(emails []string) []sdk.NotificationIntegrationAllowedRecipient {
	allowedRecipients := make([]sdk.NotificationIntegrationAllowedRecipient, len(emails))
	for i, prefix := range emails {
		allowedRecipients[i] = sdk.NotificationIntegrationAllowedRecipient{Email: prefix}
	}
	return allowedRecipients
}

// CreateEmailNotificationIntegration implements schema.CreateFunc.
func CreateEmailNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	enabled := d.Get("enabled").(bool)

	createRequest := sdk.NewCreateNotificationIntegrationRequest(id, enabled)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	emailParamsRequest := sdk.NewEmailParamsRequest()
	if v, ok := d.GetOk("allowed_recipients"); ok {
		emailParamsRequest.WithAllowedRecipients(toAllowedRecipients(expandStringList(v.(*schema.Set).List())))
	}
	createRequest.WithEmailParams(emailParamsRequest)

	err := client.NotificationIntegrations.Create(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("error creating notification integration: %w", err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadEmailNotificationIntegration(d, meta)
}

// ReadEmailNotificationIntegration implements schema.ReadFunc.
func ReadEmailNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.NotificationIntegrations.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] notification integration (%s) not found", d.Id())
		d.SetId("")
		return err
	}

	if err := d.Set("name", integration.Name); err != nil {
		return err
	}

	if err := d.Set("enabled", integration.Enabled); err != nil {
		return err
	}

	if err := d.Set("comment", integration.Comment); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	integrationProperties, err := client.NotificationIntegrations.Describe(ctx, id)
	if err != nil {
		return fmt.Errorf("could not describe notification integration: %w", err)
	}
	for _, property := range integrationProperties {
		name := property.Name
		value := property.Value

		switch name {
		case "ALLOWED_RECIPIENTS":
			if value == "" {
				if err := d.Set("allowed_recipients", make([]string, 0)); err != nil {
					return err
				}
			} else {
				if err := d.Set("allowed_recipients", strings.Split(value, ",")); err != nil {
					return err
				}
			}
		default:
			log.Printf("[WARN] unexpected notification integration property %v returned from Snowflake", name)
		}
	}

	return err
}

// UpdateEmailNotificationIntegration implements schema.UpdateFunc.
func UpdateEmailNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	var runSetStatement bool
	var runUnsetStatement bool
	setRequest := sdk.NewNotificationIntegrationSetRequest()
	unsetRequest := sdk.NewNotificationIntegrationUnsetEmailParamsRequest()
	if d.HasChange("comment") {
		v := d.Get("comment").(string)
		if v == "" {
			runUnsetStatement = true
			unsetRequest.WithComment(sdk.Bool(true))
		} else {
			runSetStatement = true
			setRequest.WithComment(sdk.String(d.Get("comment").(string)))
		}
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		setRequest.WithEnabled(sdk.Bool(d.Get("enabled").(bool)))
	}

	if d.HasChange("allowed_recipients") {
		v := d.Get("allowed_recipients").(*schema.Set).List()
		if len(v) == 0 {
			runUnsetStatement = true
			unsetRequest.WithAllowedRecipients(sdk.Bool(true))
		} else {
			runSetStatement = true
			setRequest.WithSetEmailParams(sdk.NewSetEmailParamsRequest(toAllowedRecipients(expandStringList(v))))
		}
	}

	if runSetStatement {
		err := client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithSet(setRequest))
		if err != nil {
			return fmt.Errorf("error updating notification integration: %w", err)
		}
	}

	if runUnsetStatement {
		err := client.NotificationIntegrations.Alter(ctx, sdk.NewAlterNotificationIntegrationRequest(id).WithUnsetEmailParams(unsetRequest))
		if err != nil {
			return fmt.Errorf("error updating notification integration: %w", err)
		}
	}

	return ReadEmailNotificationIntegration(d, meta)
}

// DeleteEmailNotificationIntegration implements schema.DeleteFunc.
func DeleteEmailNotificationIntegration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.NotificationIntegrations.Drop(ctx, sdk.NewDropNotificationIntegrationRequest(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
