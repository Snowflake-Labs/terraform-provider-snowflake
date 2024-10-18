// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type SecretWithGenericStringModel struct {
	Comment            tfconfig.Variable `json:"comment,omitempty"`
	Database           tfconfig.Variable `json:"database,omitempty"`
	FullyQualifiedName tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	Name               tfconfig.Variable `json:"name,omitempty"`
	Schema             tfconfig.Variable `json:"schema,omitempty"`
	SecretString       tfconfig.Variable `json:"secret_string,omitempty"`
	SecretType         tfconfig.Variable `json:"secret_type,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func SecretWithGenericString(
	resourceName string,
	database string,
	name string,
	schema string,
	secretString string,
) *SecretWithGenericStringModel {
	s := &SecretWithGenericStringModel{ResourceModelMeta: config.Meta(resourceName, resources.SecretWithGenericString)}
	s.WithDatabase(database)
	s.WithName(name)
	s.WithSchema(schema)
	s.WithSecretString(secretString)
	return s
}

func SecretWithGenericStringWithDefaultMeta(
	database string,
	name string,
	schema string,
	secretString string,
) *SecretWithGenericStringModel {
	s := &SecretWithGenericStringModel{ResourceModelMeta: config.DefaultMeta(resources.SecretWithGenericString)}
	s.WithDatabase(database)
	s.WithName(name)
	s.WithSchema(schema)
	s.WithSecretString(secretString)
	return s
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (s *SecretWithGenericStringModel) WithComment(comment string) *SecretWithGenericStringModel {
	s.Comment = tfconfig.StringVariable(comment)
	return s
}

func (s *SecretWithGenericStringModel) WithDatabase(database string) *SecretWithGenericStringModel {
	s.Database = tfconfig.StringVariable(database)
	return s
}

func (s *SecretWithGenericStringModel) WithFullyQualifiedName(fullyQualifiedName string) *SecretWithGenericStringModel {
	s.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return s
}

func (s *SecretWithGenericStringModel) WithName(name string) *SecretWithGenericStringModel {
	s.Name = tfconfig.StringVariable(name)
	return s
}

func (s *SecretWithGenericStringModel) WithSchema(schema string) *SecretWithGenericStringModel {
	s.Schema = tfconfig.StringVariable(schema)
	return s
}

func (s *SecretWithGenericStringModel) WithSecretString(secretString string) *SecretWithGenericStringModel {
	s.SecretString = tfconfig.StringVariable(secretString)
	return s
}

func (s *SecretWithGenericStringModel) WithSecretType(secretType string) *SecretWithGenericStringModel {
	s.SecretType = tfconfig.StringVariable(secretType)
	return s
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (s *SecretWithGenericStringModel) WithCommentValue(value tfconfig.Variable) *SecretWithGenericStringModel {
	s.Comment = value
	return s
}

func (s *SecretWithGenericStringModel) WithDatabaseValue(value tfconfig.Variable) *SecretWithGenericStringModel {
	s.Database = value
	return s
}

func (s *SecretWithGenericStringModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *SecretWithGenericStringModel {
	s.FullyQualifiedName = value
	return s
}

func (s *SecretWithGenericStringModel) WithNameValue(value tfconfig.Variable) *SecretWithGenericStringModel {
	s.Name = value
	return s
}

func (s *SecretWithGenericStringModel) WithSchemaValue(value tfconfig.Variable) *SecretWithGenericStringModel {
	s.Schema = value
	return s
}

func (s *SecretWithGenericStringModel) WithSecretStringValue(value tfconfig.Variable) *SecretWithGenericStringModel {
	s.SecretString = value
	return s
}

func (s *SecretWithGenericStringModel) WithSecretTypeValue(value tfconfig.Variable) *SecretWithGenericStringModel {
	s.SecretType = value
	return s
}
