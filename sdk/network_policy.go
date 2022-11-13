package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	ResourceNetworkPolicy   = "NETWORK POLICY"
	ResourceNetworkPolicies = "NETWORK POLICIES"
	DescAllowedIPList       = "ALLOWED_IP_LIST"
	DescBlockedIPList       = "BLOCKED_IP_LIST"
)

// Compile-time proof of interface implementation.
var _ NetworkPolicies = (*networkPolicies)(nil)

// NetworkPolicies describes all the network policies related methods that the
// Snowflake API supports.
type NetworkPolicies interface {
	// List all the network policies.
	List(ctx context.Context) ([]*NetworkPolicy, error)
	// Create a new network policy with the given options.
	Create(ctx context.Context, options NetworkPolicyCreateOptions) (*NetworkPolicy, error)
	// Read a network policy by its name.
	Read(ctx context.Context, policy string) (*NetworkPolicy, error)
	// Update attributes of an existing network policy.
	Update(ctx context.Context, policy string, options NetworkPolicyUpdateOptions) (*NetworkPolicy, error)
	// Delete a network policy by its name.
	Delete(ctx context.Context, policy string) error
	// Rename a network policy name.
	Rename(ctx context.Context, old string, new string) error
}

// networkPolicies implements NetworkPolicies
type networkPolicies struct {
	client *Client
}

// NetworkPolicy represents a Snowflake network policy.
type NetworkPolicy struct {
	Name          string
	Comment       string
	CreatedOn     time.Time
	AllowedIPList []string
	BlockedIPList []string
}

type networkPolicyEntity struct {
	Name      sql.NullString `db:"name"`
	Comment   sql.NullString `db:"comment"`
	CreatedOn sql.NullTime   `db:"created_on"`
}

type networkPolicyDesc struct {
	Name  sql.NullString `db:"name"`
	Value sql.NullString `db:"value"`
}

func (np *networkPolicyEntity) toNetworkPolicy() *NetworkPolicy {
	return &NetworkPolicy{
		Name:      np.Name.String,
		Comment:   np.Comment.String,
		CreatedOn: np.CreatedOn.Time,
	}
}

type NetworkPolicyProperties struct {
	// Optional: Specifies a list of IPv4 addresses that are denied access to your Snowflake account
	BlockedIPList *[]string

	// Optional: Specifies a comment for the network policy.
	Comment *string
}

// NetworkPolicyCreateOptions represents the options for creating a network policy.
type NetworkPolicyCreateOptions struct {
	*NetworkPolicyProperties

	// Required: Identifier for the network policy; must be unique for your account.
	Name string

	// Required: Specifies a list of IPv4 addresses that are allowed access to your Snowflake account.
	AllowedIPList []string
}

func (o NetworkPolicyCreateOptions) validate() error {
	if o.Name == "" {
		return errors.New("user name must not be empty")
	}
	if len(o.AllowedIPList) == 0 {
		return errors.New("allowed ip list must not be empty")
	}
	return nil
}

// NetworkPolicyUpdateOptions represents the options for updating a network policy.
type NetworkPolicyUpdateOptions struct {
	*NetworkPolicyProperties

	// Optional: Specifies a list of IPv4 addresses that are allowed access to your Snowflake account.
	AllowedIPList *[]string
}

func (np *networkPolicies) formatNetworkPolicyProperties(properties *NetworkPolicyProperties) string {
	var s string
	if properties.Comment != nil {
		s = s + " comment='" + *properties.Comment + "'"
	}
	if properties.BlockedIPList != nil && len(*properties.BlockedIPList) > 0 {
		ips := addQuote(*properties.BlockedIPList)
		s = s + " blocked_ip_list=(" + strings.Join(ips, ",") + ")"
	}
	return s
}

// Update attributes of an existing network policy.
func (np *networkPolicies) Update(ctx context.Context, policy string, options NetworkPolicyUpdateOptions) (*NetworkPolicy, error) {
	if policy == "" {
		return nil, errors.New("network policy name must not be empty")
	}
	sql := fmt.Sprintf("ALTER %s %s SET", ResourceNetworkPolicy, policy)
	if options.AllowedIPList != nil && len(*options.AllowedIPList) > 0 {
		allowed := addQuote(*options.AllowedIPList)
		sql = sql + " allowed_ip_list=(" + strings.Join(allowed, ",") + ")"
	}
	if options.NetworkPolicyProperties != nil {
		sql = sql + np.formatNetworkPolicyProperties(options.NetworkPolicyProperties)
	}
	if _, err := np.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}

	entities, err := np.list(ctx)
	if err != nil {
		return nil, err
	}
	for _, entity := range entities {
		if entity.Name == policy {
			return entity, nil
		}
	}
	return nil, ErrNoRecord
}

// Create a new network policy with the given options.
func (np *networkPolicies) Create(ctx context.Context, options NetworkPolicyCreateOptions) (*NetworkPolicy, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate create options: %w", err)
	}

	sql := fmt.Sprintf("CREATE %s %s", ResourceNetworkPolicy, options.Name)
	if len(options.AllowedIPList) > 0 {
		allowed := addQuote(options.AllowedIPList)
		sql = sql + " allowed_ip_list=(" + strings.Join(allowed, ",") + ")"
	}
	if options.NetworkPolicyProperties != nil {
		sql = sql + np.formatNetworkPolicyProperties(options.NetworkPolicyProperties)
	}
	if _, err := np.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}

	entities, err := np.list(ctx)
	if err != nil {
		return nil, err
	}
	for _, entity := range entities {
		if entity.Name == options.Name {
			return entity, nil
		}
	}
	return nil, ErrNoRecord
}

func (np *networkPolicies) list(ctx context.Context) ([]*NetworkPolicy, error) {
	sql := fmt.Sprintf("SHOW %s", ResourceNetworkPolicies)
	rows, err := np.client.query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	entities := []*NetworkPolicy{}
	for rows.Next() {
		var entity networkPolicyEntity
		if err := rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		networkPolicy := entity.toNetworkPolicy()

		values, err := np.desc(ctx, entity.Name.String)
		if err != nil {
			return nil, fmt.Errorf("desc: %w", err)
		}
		for _, item := range values {
			if item.Name.String == DescAllowedIPList {
				networkPolicy.AllowedIPList = strings.Split(item.Value.String, ",")
			}
			if item.Name.String == DescBlockedIPList {
				networkPolicy.BlockedIPList = strings.Split(item.Value.String, ",")
			}
		}
		entities = append(entities, networkPolicy)
	}
	return entities, nil
}

// List all the network policies.
func (np *networkPolicies) List(ctx context.Context) ([]*NetworkPolicy, error) {
	return np.list(ctx)
}

func (np *networkPolicies) desc(ctx context.Context, policy string) ([]*networkPolicyDesc, error) {
	sql := fmt.Sprintf("DESC %s %s", ResourceNetworkPolicy, policy)
	rows, err := np.client.query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	entities := []*networkPolicyDesc{}
	for rows.Next() {
		var entity networkPolicyDesc
		if err := rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		entities = append(entities, &entity)
	}
	return entities, nil
}

// Read a network policy by its name.
func (np *networkPolicies) Read(ctx context.Context, policy string) (*NetworkPolicy, error) {
	entities, err := np.list(ctx)
	if err != nil {
		return nil, err
	}
	for _, entity := range entities {
		if entity.Name == policy {
			return entity, nil
		}
	}
	return nil, ErrNoRecord
}

// Delete a network policy by its name.
func (np *networkPolicies) Delete(ctx context.Context, policy string) error {
	return np.client.drop(ctx, ResourceNetworkPolicy, policy)
}

// Rename a network policy name.
func (np *networkPolicies) Rename(ctx context.Context, old string, new string) error {
	return np.client.rename(ctx, ResourceNetworkPolicy, old, new)
}
