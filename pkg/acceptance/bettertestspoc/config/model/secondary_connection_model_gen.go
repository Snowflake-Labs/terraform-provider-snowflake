// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type SecondaryConnectionModel struct {
	AsReplicaOf        tfconfig.Variable `json:"as_replica_of,omitempty"`
	Comment            tfconfig.Variable `json:"comment,omitempty"`
	FullyQualifiedName tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	IsPrimary          tfconfig.Variable `json:"is_primary,omitempty"`
	Name               tfconfig.Variable `json:"name,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func SecondaryConnection(
	resourceName string,
	asReplicaOf string,
	name string,
) *SecondaryConnectionModel {
	s := &SecondaryConnectionModel{ResourceModelMeta: config.Meta(resourceName, resources.SecondaryConnection)}
	s.WithAsReplicaOf(asReplicaOf)
	s.WithName(name)
	return s
}

func SecondaryConnectionWithDefaultMeta(
	asReplicaOf string,
	name string,
) *SecondaryConnectionModel {
	s := &SecondaryConnectionModel{ResourceModelMeta: config.DefaultMeta(resources.SecondaryConnection)}
	s.WithAsReplicaOf(asReplicaOf)
	s.WithName(name)
	return s
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (s *SecondaryConnectionModel) WithAsReplicaOf(asReplicaOf string) *SecondaryConnectionModel {
	s.AsReplicaOf = tfconfig.StringVariable(asReplicaOf)
	return s
}

func (s *SecondaryConnectionModel) WithComment(comment string) *SecondaryConnectionModel {
	s.Comment = tfconfig.StringVariable(comment)
	return s
}

func (s *SecondaryConnectionModel) WithFullyQualifiedName(fullyQualifiedName string) *SecondaryConnectionModel {
	s.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return s
}

func (s *SecondaryConnectionModel) WithIsPrimary(isPrimary bool) *SecondaryConnectionModel {
	s.IsPrimary = tfconfig.BoolVariable(isPrimary)
	return s
}

func (s *SecondaryConnectionModel) WithName(name string) *SecondaryConnectionModel {
	s.Name = tfconfig.StringVariable(name)
	return s
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (s *SecondaryConnectionModel) WithAsReplicaOfValue(value tfconfig.Variable) *SecondaryConnectionModel {
	s.AsReplicaOf = value
	return s
}

func (s *SecondaryConnectionModel) WithCommentValue(value tfconfig.Variable) *SecondaryConnectionModel {
	s.Comment = value
	return s
}

func (s *SecondaryConnectionModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *SecondaryConnectionModel {
	s.FullyQualifiedName = value
	return s
}

func (s *SecondaryConnectionModel) WithIsPrimaryValue(value tfconfig.Variable) *SecondaryConnectionModel {
	s.IsPrimary = value
	return s
}

func (s *SecondaryConnectionModel) WithNameValue(value tfconfig.Variable) *SecondaryConnectionModel {
	s.Name = value
	return s
}
