package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type GrantOwnershipTargetRoleKind string

const (
	ToAccountGrantOwnershipTargetRoleKind  GrantOwnershipTargetRoleKind = "ToAccountRole"
	ToDatabaseGrantOwnershipTargetRoleKind GrantOwnershipTargetRoleKind = "ToDatabaseRole"
)

type OutboundPrivilegesBehavior string

const (
	CopyOutboundPrivilegesBehavior   OutboundPrivilegesBehavior = "COPY"
	RevokeOutboundPrivilegesBehavior OutboundPrivilegesBehavior = "REVOKE"
)

func (o OutboundPrivilegesBehavior) ToOwnershipCurrentGrantsOutboundPrivileges() *sdk.OwnershipCurrentGrantsOutboundPrivileges {
	switch o {
	case CopyOutboundPrivilegesBehavior:
		return sdk.Pointer(sdk.Copy)
	case RevokeOutboundPrivilegesBehavior:
		return sdk.Pointer(sdk.Revoke)
	default:
		return nil
	}
}

type GrantOwnershipKind string

const (
	OnObjectGrantOwnershipKind GrantOwnershipKind = "OnObject"
	OnAllGrantOwnershipKind    GrantOwnershipKind = "OnAll"
	OnFutureGrantOwnershipKind GrantOwnershipKind = "OnFuture"
)

type GrantOwnershipId struct {
	GrantOwnershipTargetRoleKind GrantOwnershipTargetRoleKind
	AccountRoleName              sdk.AccountObjectIdentifier
	DatabaseRoleName             sdk.DatabaseObjectIdentifier
	OutboundPrivilegesBehavior   *OutboundPrivilegesBehavior
	Kind                         GrantOwnershipKind
	Data                         fmt.Stringer
}

type OnObjectGrantOwnershipData struct {
	ObjectType sdk.ObjectType
	ObjectName sdk.ObjectIdentifier
}

func (g *OnObjectGrantOwnershipData) String() string {
	var parts []string
	parts = append(parts, g.ObjectType.String())
	parts = append(parts, g.ObjectName.FullyQualifiedName())
	return strings.Join(parts, helpers.IDDelimiter)
}

func (g *GrantOwnershipId) String() string {
	var parts []string
	parts = append(parts, string(g.GrantOwnershipTargetRoleKind))
	switch g.GrantOwnershipTargetRoleKind {
	case ToAccountGrantOwnershipTargetRoleKind:
		parts = append(parts, g.AccountRoleName.FullyQualifiedName())
	case ToDatabaseGrantOwnershipTargetRoleKind:
		parts = append(parts, g.DatabaseRoleName.FullyQualifiedName())
	}
	if g.OutboundPrivilegesBehavior != nil {
		parts = append(parts, string(*g.OutboundPrivilegesBehavior))
	} else {
		parts = append(parts, "")
	}
	parts = append(parts, string(g.Kind))
	data := g.Data.String()
	if len(data) > 0 {
		parts = append(parts, data)
	}
	return strings.Join(parts, helpers.IDDelimiter)
}

func ParseGrantOwnershipId(id string) (*GrantOwnershipId, error) {
	grantOwnershipId := new(GrantOwnershipId)

	parts := strings.Split(id, helpers.IDDelimiter)
	if len(parts) < 5 {
		return grantOwnershipId, sdk.NewError(`grant ownership identifier should hold at least 5 parts "<target_role_kind>|<role_name>|<outbound_privileges_behavior>|<grant_type>|<grant_data>"`)
	}

	grantOwnershipId.GrantOwnershipTargetRoleKind = GrantOwnershipTargetRoleKind(parts[0])
	switch grantOwnershipId.GrantOwnershipTargetRoleKind {
	case ToAccountGrantOwnershipTargetRoleKind:
		grantOwnershipId.AccountRoleName = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[1])
	case ToDatabaseGrantOwnershipTargetRoleKind:
		grantOwnershipId.DatabaseRoleName = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[1])
	default:
		return grantOwnershipId, sdk.NewError(fmt.Sprintf("unknown GrantOwnershipTargetRoleKind: %v, valid options are %v | %v", grantOwnershipId.GrantOwnershipTargetRoleKind, ToAccountGrantOwnershipTargetRoleKind, ToDatabaseGrantOwnershipTargetRoleKind))
	}

	if len(parts[2]) > 0 {
		switch outboundPrivilegesBehavior := OutboundPrivilegesBehavior(parts[2]); outboundPrivilegesBehavior {
		case CopyOutboundPrivilegesBehavior, RevokeOutboundPrivilegesBehavior:
			grantOwnershipId.OutboundPrivilegesBehavior = sdk.Pointer(outboundPrivilegesBehavior)
		default:
			return grantOwnershipId, sdk.NewError(fmt.Sprintf("unknown OutboundPrivilegesBehavior: %v, valid options are %v | %v", outboundPrivilegesBehavior, CopyOutboundPrivilegesBehavior, RevokeOutboundPrivilegesBehavior))
		}
	}

	grantOwnershipId.Kind = GrantOwnershipKind(parts[3])
	switch grantOwnershipId.Kind {
	case OnObjectGrantOwnershipKind:
		if len(parts) != 6 {
			return grantOwnershipId, sdk.NewError(`grant ownership identifier should consist of 6 parts "<target_role_kind>|<role_name>|<outbound_privileges_behavior>|OnObject|<object_type>|<object_name>"`)
		}
		// TODO: Custom type for OnObject grant - because ObjectName can be pretty much any of the possible identifiers
		grantOwnershipId.Data = &OnObjectGrantOwnershipData{
			ObjectType: sdk.ObjectType(parts[4]),
			ObjectName: sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[5]), // TODO: Fix should accept any identifier (most likely have to handle case by case for every object type)
		}
	case OnAllGrantOwnershipKind, OnFutureGrantOwnershipKind:
		bulkOperationGrantData := &BulkOperationGrantData{
			ObjectNamePlural: sdk.PluralObjectType(parts[4]),
		}
		if len(parts) != 7 {
			return grantOwnershipId, sdk.NewError(`grant ownership identifier should consist of 7 parts "<target_role_kind>|<role_name>|<outbound_privileges_behavior>|On[All or Future]|<object_type_plural>|In[Database or Schema]|<identifier>"`)
		}
		bulkOperationGrantData.Kind = BulkOperationGrantKind(parts[5])
		switch bulkOperationGrantData.Kind {
		case InDatabaseBulkOperationGrantKind:
			bulkOperationGrantData.Database = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[6]))
		case InSchemaBulkOperationGrantKind:
			bulkOperationGrantData.Schema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[6]))
		default:
			return grantOwnershipId, sdk.NewError(fmt.Sprintf("invalid BulkOperationGrantKind: %s, valid options are %v | %v", bulkOperationGrantData.Kind, InDatabaseBulkOperationGrantKind, InSchemaBulkOperationGrantKind))
		}
		grantOwnershipId.Data = bulkOperationGrantData
	default:
		return grantOwnershipId, sdk.NewError(fmt.Sprintf("unknown GrantOwnershipKind: %v", grantOwnershipId.Kind))
	}

	return grantOwnershipId, nil
}
