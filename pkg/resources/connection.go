package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*       --- DESIGN PROP ---

// basic create
resource "snowflake_connection" "basic" {
    name = "example_connection"
}

// create replica of
resource "snowflake_connection" "replica" {
    name = "example_connection"
    replica_of = snowflake_connection.example.name
}

// enable failover
resource "snowflake_connection" "enable_failover" {
    name = "example_connection"
    enable_failover {
        to_accounts = [
            { account_identifier = "example_org.example_account" },
            { account_identifier = "sec_example_org.sec_example_account" },
        ]
        ignore_edition_check = true
    }
}

// promote to primary
resource "snowflake_connection" "promote_to_primary" {
    name = "example_connection"
    primary = true
}
*/

var connectionSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the connection. Must start with an alphabetic character and may only contain letters, decimal digits (0-9), and underscores (_). For a primary connection, the name must be unique across connection names and account names in the organization. For a secondary connection, the name must match the name of its primary connection."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"replica_of": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the identifier for a primary connection from which to create a replica (i.e. a secondary connection).",
	},
	"enable_failover": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Enables failover for given connection.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"to_accounts": {
					Type:        schema.TypeList,
					Required:    true,
					Description: "Specifies a list of accounts in your organization where a secondary connection for this primary connection can be promoted to serve as the primary connection. Include your organization name for each account in the list.",
					MinItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"account_identifier": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies account identifier for which replication should be enabled. The account identifiers should be in the form of `\"<organization_name>\".\"<account_name>\"`.",
							},
						},
					},
				},
				"ignore_edition_check": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Allows replicating data to accounts on lower editions in either of the following scenarios: " +
						"1. The primary database is in a Business Critical (or higher) account but one or more of the accounts approved for replication are on lower editions. Business Critical Edition is intended for Snowflake accounts with extremely sensitive data. " +
						"2. The primary database is in a Business Critical (or higher) account and a signed business associate agreement is in place to store PHI data in the account per HIPAA and HITRUST regulations, but no such agreement is in place for one or more of the accounts approved for replication, regardless if they are Business Critical (or higher) accounts. " +
						"Both scenarios are prohibited by default in an effort to help prevent account administrators for Business Critical (or higher) accounts from inadvertently replicating sensitive data to accounts on lower editions.",
				},
			},
		},
	},
	"primary": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Promote connection to serve as primary connection.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the connection.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CONNECTIONS` for the given secret.",
		Elem: &schema.Resource{
			Schema: schemas.ShowConnectionSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}


