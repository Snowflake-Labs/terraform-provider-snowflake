// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type MaskingPolicyModel struct {
	Argument            tfconfig.Variable `json:"argument,omitempty"`
	Body                tfconfig.Variable `json:"body,omitempty"`
	Comment             tfconfig.Variable `json:"comment,omitempty"`
	Database            tfconfig.Variable `json:"database,omitempty"`
	ExemptOtherPolicies tfconfig.Variable `json:"exempt_other_policies,omitempty"`
	FullyQualifiedName  tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	Name                tfconfig.Variable `json:"name,omitempty"`
	ReturnDataType      tfconfig.Variable `json:"return_data_type,omitempty"`
	Schema              tfconfig.Variable `json:"schema,omitempty"`

	// added manually as a PoC
	DynamicBlock *config.DynamicBlock `json:"dynamic,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func MaskingPolicy(
	resourceName string,
	argument []sdk.TableColumnSignature,
	body string,
	database string,
	name string,
	returnDataType string,
	schema string,
) *MaskingPolicyModel {
	m := &MaskingPolicyModel{ResourceModelMeta: config.Meta(resourceName, resources.MaskingPolicy)}
	m.WithArgument(argument)
	m.WithBody(body)
	m.WithDatabase(database)
	m.WithName(name)
	m.WithReturnDataType(returnDataType)
	m.WithSchema(schema)
	return m
}

func MaskingPolicyWithDefaultMeta(
	argument []sdk.TableColumnSignature,
	body string,
	database string,
	name string,
	returnDataType string,
	schema string,
) *MaskingPolicyModel {
	m := &MaskingPolicyModel{ResourceModelMeta: config.DefaultMeta(resources.MaskingPolicy)}
	m.WithArgument(argument)
	m.WithBody(body)
	m.WithDatabase(database)
	m.WithName(name)
	m.WithReturnDataType(returnDataType)
	m.WithSchema(schema)
	return m
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (m *MaskingPolicyModel) MarshalJSON() ([]byte, error) {
	type Alias MaskingPolicyModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(m),
		DependsOn: m.DependsOn(),
	})
}

func (m *MaskingPolicyModel) WithDependsOn(values ...string) *MaskingPolicyModel {
	m.SetDependsOn(values...)
	return m
}

// added manually as a PoC
func (m *MaskingPolicyModel) WithDynamicBlock(dynamicBlock *config.DynamicBlock) *MaskingPolicyModel {
	m.DynamicBlock = dynamicBlock
	return m
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// argument attribute type is not yet supported, so WithArgument can't be generated

func (m *MaskingPolicyModel) WithBody(body string) *MaskingPolicyModel {
	m.Body = tfconfig.StringVariable(body)
	return m
}

func (m *MaskingPolicyModel) WithComment(comment string) *MaskingPolicyModel {
	m.Comment = tfconfig.StringVariable(comment)
	return m
}

func (m *MaskingPolicyModel) WithDatabase(database string) *MaskingPolicyModel {
	m.Database = tfconfig.StringVariable(database)
	return m
}

func (m *MaskingPolicyModel) WithExemptOtherPolicies(exemptOtherPolicies string) *MaskingPolicyModel {
	m.ExemptOtherPolicies = tfconfig.StringVariable(exemptOtherPolicies)
	return m
}

func (m *MaskingPolicyModel) WithFullyQualifiedName(fullyQualifiedName string) *MaskingPolicyModel {
	m.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return m
}

func (m *MaskingPolicyModel) WithName(name string) *MaskingPolicyModel {
	m.Name = tfconfig.StringVariable(name)
	return m
}

func (m *MaskingPolicyModel) WithReturnDataType(returnDataType string) *MaskingPolicyModel {
	m.ReturnDataType = tfconfig.StringVariable(returnDataType)
	return m
}

func (m *MaskingPolicyModel) WithSchema(schema string) *MaskingPolicyModel {
	m.Schema = tfconfig.StringVariable(schema)
	return m
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (m *MaskingPolicyModel) WithArgumentValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.Argument = value
	return m
}

func (m *MaskingPolicyModel) WithBodyValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.Body = value
	return m
}

func (m *MaskingPolicyModel) WithCommentValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.Comment = value
	return m
}

func (m *MaskingPolicyModel) WithDatabaseValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.Database = value
	return m
}

func (m *MaskingPolicyModel) WithExemptOtherPoliciesValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.ExemptOtherPolicies = value
	return m
}

func (m *MaskingPolicyModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.FullyQualifiedName = value
	return m
}

func (m *MaskingPolicyModel) WithNameValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.Name = value
	return m
}

func (m *MaskingPolicyModel) WithReturnDataTypeValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.ReturnDataType = value
	return m
}

func (m *MaskingPolicyModel) WithSchemaValue(value tfconfig.Variable) *MaskingPolicyModel {
	m.Schema = value
	return m
}
