package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var storageIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "EXTERNAL_STAGE",
		ValidateFunc: validation.StringInSlice([]string{"EXTERNAL_STAGE"}, true),
		ForceNew:     true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"storage_allowed_locations": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Explicitly limits external stages that use the integration to reference one or more storage locations.",
		MinItems:    1,
	},
	"storage_blocked_locations": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Explicitly prohibits external stages that use the integration from referencing one or more storage locations.",
	},
	// TODO (SNOW-1015282): Remove S3gov option before going into V1
	"storage_provider": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"S3", "S3gov", "GCS", "AZURE", "S3GOV"}, false),
	},
	"storage_aws_external_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role.",
	},
	"storage_aws_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
	},
	"storage_aws_object_acl": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"bucket-owner-full-control"}, false),
		Description:  "\"bucket-owner-full-control\" Enables support for AWS access control lists (ACLs) to grant the bucket owner full control.",
	},
	"storage_aws_role_arn": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"azure_tenant_id": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"azure_consent_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The consent URL that is used to create an Azure Snowflake service principle inside your tenant.",
	},
	"azure_multi_tenant_app_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "This is the name of the Snowflake client application created for your account.",
	},
	"storage_gcp_service_account": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "This is the name of the Snowflake Google Service Account created for your account.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the storage integration was created.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// StorageIntegration returns a pointer to the resource representing a storage integration.
func StorageIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateStorageIntegration,
		Read:   ReadStorageIntegration,
		Update: UpdateStorageIntegration,
		Delete: DeleteStorageIntegration,

		Schema: storageIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateStorageIntegration(d *schema.ResourceData, meta any) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	name := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("name").(string))
	enabled := d.Get("enabled").(bool)
	stringStorageAllowedLocations := expandStringList(d.Get("storage_allowed_locations").([]any))
	storageAllowedLocations := make([]sdk.StorageLocation, len(stringStorageAllowedLocations))
	for i, loc := range stringStorageAllowedLocations {
		storageAllowedLocations[i] = sdk.StorageLocation{
			Path: loc,
		}
	}

	req := sdk.NewCreateStorageIntegrationRequest(name, enabled, storageAllowedLocations)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(sdk.String(v.(string)))
	}

	if _, ok := d.GetOk("storage_blocked_locations"); ok {
		stringStorageBlockedLocations := expandStringList(d.Get("storage_blocked_locations").([]any))
		storageBlockedLocations := make([]sdk.StorageLocation, len(stringStorageBlockedLocations))
		for i, loc := range stringStorageBlockedLocations {
			storageBlockedLocations[i] = sdk.StorageLocation{
				Path: loc,
			}
		}
	}

	storageProvider := d.Get("storage_provider").(string)

	switch storageProvider {
	case "S3", "S3GOV", "S3gov":
		v, ok := d.GetOk("storage_aws_role_arn")
		if !ok {
			return fmt.Errorf("if you use the S3 storage provider you must specify a storage_aws_role_arn")
		}

		s3Params := sdk.NewS3StorageParamsRequest(v.(string))
		if _, ok := d.GetOk("storage_aws_object_acl"); ok {
			s3Params.WithStorageAwsObjectAcl(sdk.String(d.Get("storage_aws_object_acl").(string)))
		}
		req.WithS3StorageProviderParams(s3Params)
	case "AZURE":
		v, ok := d.GetOk("azure_tenant_id")
		if !ok {
			return fmt.Errorf("if you use the Azure storage provider you must specify an azure_tenant_id")
		}
		req.WithAzureStorageProviderParams(sdk.NewAzureStorageParamsRequest(sdk.String(v.(string))))
	case "GCS":
		req.WithGCSStorageProviderParams(sdk.NewGCSStorageParamsRequest())
	default:
		return fmt.Errorf("unexpected provider %v", storageProvider)
	}

	if err := client.StorageIntegrations.Create(ctx, req); err != nil {
		return fmt.Errorf("error creating storage integration: %w", err)
	}

	d.SetId(helpers.EncodeSnowflakeID(name))
	return ReadStorageIntegration(d, meta)
}

func ReadStorageIntegration(d *schema.ResourceData, meta any) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id, ok := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	if !ok {
		return fmt.Errorf("storage integration read, error decoding id: %s as sdk.AccountObjectIdentifier, got: %T", d.Id(), id)
	}

	s, err := client.StorageIntegrations.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] storage integration (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return err
	}

	if s.Category != "STORAGE" {
		return fmt.Errorf("expected %v to be a STORAGE integration, got %v", d.Id(), s.Category)
	}
	if err := d.Set("name", s.Name); err != nil {
		return err
	}
	if err := d.Set("type", s.StorageType); err != nil {
		return err
	}
	if err := d.Set("created_on", s.CreatedOn.String()); err != nil {
		return err
	}
	if err := d.Set("enabled", s.Enabled); err != nil {
		return err
	}
	if err := d.Set("comment", s.Comment); err != nil {
		return err
	}

	storageIntegrationProps, err := client.StorageIntegrations.Describe(ctx, id)
	if err != nil {
		return fmt.Errorf("could not describe storage integration (%s), err = %w", d.Id(), err)
	}

	for _, prop := range storageIntegrationProps {
		switch prop.Name {
		case "STORAGE_PROVIDER":
			if err := d.Set("storage_provider", prop.Value); err != nil {
				return err
			}
		case "STORAGE_ALLOWED_LOCATIONS":
			if err := d.Set("storage_allowed_locations", strings.Split(prop.Value, ",")); err != nil {
				return err
			}
		case "STORAGE_BLOCKED_LOCATIONS":
			if prop.Value != "" {
				if err := d.Set("storage_blocked_locations", strings.Split(prop.Value, ",")); err != nil {
					return err
				}
			}
		case "STORAGE_AWS_IAM_USER_ARN":
			if err := d.Set("storage_aws_iam_user_arn", prop.Value); err != nil {
				return err
			}
		case "STORAGE_AWS_OBJECT_ACL":
			if prop.Value != "" {
				if err := d.Set("storage_aws_object_acl", prop.Value); err != nil {
					return err
				}
			}
		case "STORAGE_AWS_ROLE_ARN":
			if err := d.Set("storage_aws_role_arn", prop.Value); err != nil {
				return err
			}
		case "STORAGE_AWS_EXTERNAL_ID":
			if err := d.Set("storage_aws_external_id", prop.Value); err != nil {
				return err
			}
		case "STORAGE_GCP_SERVICE_ACCOUNT":
			if err := d.Set("storage_gcp_service_account", prop.Value); err != nil {
				return err
			}
		case "AZURE_CONSENT_URL":
			if err := d.Set("azure_consent_url", prop.Value); err != nil {
				return err
			}
		case "AZURE_MULTI_TENANT_APP_NAME":
			if err := d.Set("azure_multi_tenant_app_name", prop.Value); err != nil {
				return err
			}
		}
	}

	return err
}

func UpdateStorageIntegration(d *schema.ResourceData, meta any) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id, ok := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	if !ok {
		return fmt.Errorf("storage integration update, error decoding id: %s as sdk.AccountObjectIdentifier, got: %T", d.Id(), id)
	}

	var runSetStatement bool
	setReq := sdk.NewStorageIntegrationSetRequest()

	if d.HasChange("comment") {
		runSetStatement = true
		setReq.WithComment(sdk.String(d.Get("comment").(string)))
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		setReq.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("storage_allowed_locations") {
		runSetStatement = true
		stringStorageAllowedLocations := expandStringList(d.Get("storage_allowed_locations").([]any))
		storageAllowedLocations := make([]sdk.StorageLocation, len(stringStorageAllowedLocations))
		for i, loc := range stringStorageAllowedLocations {
			storageAllowedLocations[i] = sdk.StorageLocation{
				Path: loc,
			}
		}
		setReq.WithStorageAllowedLocations(storageAllowedLocations)
	}

	// We need to UNSET this if we remove all storage blocked locations, because Snowflake won't accept an empty list
	if d.HasChange("storage_blocked_locations") {
		v := d.Get("storage_blocked_locations").([]interface{})
		if len(v) == 0 {
			if err := client.StorageIntegrations.Alter(ctx, sdk.NewAlterStorageIntegrationRequest(id).
				WithUnset(sdk.NewStorageIntegrationUnsetRequest().WithStorageBlockedLocations(sdk.Bool(true)))); err != nil {
				return fmt.Errorf("error unsetting storage_blocked_locations, err = %w", err)
			}
		} else {
			runSetStatement = true
			stringStorageBlockedLocations := expandStringList(d.Get("storage_allowed_locations").([]any))
			storageBlockedLocations := make([]sdk.StorageLocation, len(stringStorageBlockedLocations))
			for i, loc := range stringStorageBlockedLocations {
				storageBlockedLocations[i] = sdk.StorageLocation{
					Path: loc,
				}
			}
			setReq.WithStorageBlockedLocations(storageBlockedLocations)
		}
	}

	if d.HasChange("storage_aws_role_arn") || d.HasChange("storage_aws_object_acl") {
		runSetStatement = true
		s3SetParams := sdk.NewSetS3StorageParamsRequest(d.Get("storage_aws_role_arn").(string))

		if d.HasChange("storage_aws_object_acl") {
			if v, ok := d.GetOk("storage_aws_object_acl"); ok {
				s3SetParams.WithStorageAwsObjectAcl(sdk.String(v.(string)))
			} else {
				if err := client.StorageIntegrations.Alter(ctx, sdk.NewAlterStorageIntegrationRequest(id).
					WithUnset(sdk.NewStorageIntegrationUnsetRequest().WithStorageAwsObjectAcl(sdk.Bool(true)))); err != nil {
					return fmt.Errorf("error unsetting storage_aws_object_acl, err = %w", err)
				}
			}
		}

		setReq.WithS3Params(s3SetParams)
	}

	if d.HasChange("azure_tenant_id") {
		runSetStatement = true
		setReq.WithAzureParams(sdk.NewSetAzureStorageParamsRequest(d.Get("azure_tenant_id").(string)))
	}

	if runSetStatement {
		if err := client.StorageIntegrations.Alter(ctx, sdk.NewAlterStorageIntegrationRequest(id).WithSet(setReq)); err != nil {
			return fmt.Errorf("error updating storage integration, err = %w", err)
		}
	}

	return ReadStorageIntegration(d, meta)
}

func DeleteStorageIntegration(d *schema.ResourceData, meta any) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id, ok := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	if !ok {
		return fmt.Errorf("storage integration delete, error decoding id: %s as sdk.AccountObjectIdentifier, got: %T", d.Id(), id)
	}
	if err := client.StorageIntegrations.Drop(ctx, sdk.NewDropStorageIntegrationRequest(id)); err != nil {
		return fmt.Errorf("error dropping storage integration (%s), err = %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}
