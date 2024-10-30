package model

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func (c *ConnectionModel) WithAsReplicaOfIdentifier(asReplicaOf sdk.ExternalObjectIdentifier) *ConnectionModel {
	asReplicaOfString := strings.ReplaceAll(asReplicaOf.FullyQualifiedName(), `"`, "")
	c.AsReplicaOf = config.StringVariable(asReplicaOfString)
	return c
}

func (c *ConnectionModel) WithEnableFailover(toAccount ...sdk.AccountIdentifier) *ConnectionModel {
	variables := make([]config.Variable, 0)
	for _, v := range toAccount {
		variables = append(variables, config.StringVariable(v.Name()))
	}

	c.EnableFailoverToAccounts = config.ListVariable(variables...)

	return c
}
