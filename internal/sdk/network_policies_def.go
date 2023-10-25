// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	ip = g.QueryStruct("IP").
		Text("IP", g.KeywordOptions().SingleQuotes().Required())

	NetworkPoliciesDef = g.NewInterface(
		"NetworkPolicies",
		"NetworkPolicy",
		g.KindOfT[AccountObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-network-policy",
			g.QueryStruct("CreateNetworkPolicies").
				Create().
				OrReplace().
				SQL("NETWORK POLICY").
				Name().
				ListQueryStructField("AllowedIpList", ip, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
				ListQueryStructField("BlockedIpList", ip, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy",
			g.QueryStruct("AlterNetworkPolicy").
				Alter().
				SQL("NETWORK POLICY").
				IfExists().
				Name().
				OptionalQueryStructField(
					"Set",
					g.QueryStruct("NetworkPolicySet").
						ListQueryStructField("AllowedIpList", ip, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
						ListQueryStructField("BlockedIpList", ip, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
						OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
						WithValidation(g.AtLeastOneValueSet, "AllowedIpList", "BlockedIpList", "Comment"),
					g.KeywordOptions().SQL("SET"),
				).
				OptionalSQL("UNSET COMMENT").
				Identifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ExactlyOneValueSet, "Set", "UnsetComment", "RenameTo").
				WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-network-policy",
			g.QueryStruct("DropNetworkPolicy").
				Drop().
				SQL("NETWORK POLICY").
				IfExists().
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		).
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-network-policies",
			g.DbStruct("showNetworkPolicyDBRow").
				Field("created_on", "string").
				Field("name", "string").
				Field("comment", "string").
				Field("entries_in_allowed_ip_list", "int").
				Field("entries_in_blocked_ip_list", "int"),
			g.PlainStruct("NetworkPolicy").
				Field("CreatedOn", "string").
				Field("Name", "string").
				Field("Comment", "string").
				Field("EntriesInAllowedIpList", "int").
				Field("EntriesInBlockedIpList", "int"),
			g.QueryStruct("ShowNetworkPolicies").
				Show().
				SQL("NETWORK POLICIES"),
		).
		DescribeOperation(
			g.DescriptionMappingKindSlice,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-network-policy",
			g.DbStruct("describeNetworkPolicyDBRow").
				Field("name", "string").
				Field("value", "string"),
			g.PlainStruct("NetworkPolicyDescription").
				Field("Name", "string").
				Field("Value", "string"),
			g.QueryStruct("DescribeNetworkPolicy").
				Describe().
				SQL("NETWORK POLICY").
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
