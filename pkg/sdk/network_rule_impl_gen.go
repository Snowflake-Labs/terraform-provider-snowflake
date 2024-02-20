package sdk

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"strings"
)

var _ NetworkRules = (*networkRules)(nil)

type networkRules struct {
	client *Client
}

func (v *networkRules) Create(ctx context.Context, request *CreateNetworkRuleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkRules) Alter(ctx context.Context, request *AlterNetworkRuleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkRules) Drop(ctx context.Context, request *DropNetworkRuleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkRules) Show(ctx context.Context, request *ShowNetworkRuleRequest) ([]NetworkRule, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[ShowNetworkRulesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[ShowNetworkRulesRow, NetworkRule](dbRows)
	return resultList, nil
}

func (v *networkRules) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*NetworkRule, error) {
	networkRules, err := v.Show(ctx, NewShowNetworkRuleRequest().WithIn(&In{
		Schema: NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName()),
	}).WithLike(&Like{
		Pattern: String(id.Name()),
	}))
	if err != nil {
		return nil, err
	}
	return collections.FindOne(networkRules, func(r NetworkRule) bool { return r.Name == id.Name() })
}

func (v *networkRules) Describe(ctx context.Context, id SchemaObjectIdentifier) (*NetworkRuleDetails, error) {
	opts := &DescribeNetworkRuleOptions{
		name: id,
	}
	result, err := validateAndQueryOne[DescNetworkRulesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateNetworkRuleRequest) toOpts() *CreateNetworkRuleOptions {
	opts := &CreateNetworkRuleOptions{
		OrReplace: r.OrReplace,
		name:      r.name,
		Type:      r.Type,
		ValueList: r.ValueList,
		Mode:      r.Mode,
		Comment:   r.Comment,
	}
	return opts
}

func (r *AlterNetworkRuleRequest) toOpts() *AlterNetworkRuleOptions {
	opts := &AlterNetworkRuleOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.Set != nil {
		opts.Set = &NetworkRuleSet{
			ValueList: r.Set.ValueList,
			Comment:   r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &NetworkRuleUnset{
			ValueList: r.Unset.ValueList,
			Comment:   r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropNetworkRuleRequest) toOpts() *DropNetworkRuleOptions {
	opts := &DropNetworkRuleOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowNetworkRuleRequest) toOpts() *ShowNetworkRuleOptions {
	opts := &ShowNetworkRuleOptions{
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (row ShowNetworkRulesRow) convert() *NetworkRule {
	return &NetworkRule{
		CreatedOn:          row.CreatedOn,
		Name:               row.Name,
		DatabaseName:       row.DatabaseName,
		SchemaName:         row.SchemaName,
		Owner:              row.Owner,
		Comment:            row.Comment,
		Type:               row.Type,
		Mode:               row.Mode,
		EntriesInValueList: row.EntriesInValueList,
		OwnerRoleType:      row.OwnerRoleType,
	}
}

func (r *DescribeNetworkRuleRequest) toOpts() *DescribeNetworkRuleOptions {
	opts := &DescribeNetworkRuleOptions{
		name: r.name,
	}
	return opts
}

func (row DescNetworkRulesRow) convert() *NetworkRuleDetails {
	valueList := strings.Split(row.ValueList, ",")
	if len(valueList) == 1 && valueList[0] == "" {
		valueList = []string{}
	}
	return &NetworkRuleDetails{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Owner:        row.Owner,
		Comment:      row.Comment,
		Type:         row.Type,
		Mode:         row.Mode,
		ValueList:    valueList,
	}
}
