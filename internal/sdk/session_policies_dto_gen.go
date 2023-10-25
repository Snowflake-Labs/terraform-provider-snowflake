// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSessionPolicyOptions]   = new(CreateSessionPolicyRequest)
	_ optionsProvider[AlterSessionPolicyOptions]    = new(AlterSessionPolicyRequest)
	_ optionsProvider[DropSessionPolicyOptions]     = new(DropSessionPolicyRequest)
	_ optionsProvider[ShowSessionPolicyOptions]     = new(ShowSessionPolicyRequest)
	_ optionsProvider[DescribeSessionPolicyOptions] = new(DescribeSessionPolicyRequest)
)

type CreateSessionPolicyRequest struct {
	OrReplace                *bool
	IfNotExists              *bool
	name                     SchemaObjectIdentifier // required
	SessionIdleTimeoutMins   *int
	SessionUiIdleTimeoutMins *int
	Comment                  *string
}

type AlterSessionPolicyRequest struct {
	IfExists  *bool
	name      SchemaObjectIdentifier // required
	RenameTo  *SchemaObjectIdentifier
	Set       *SessionPolicySetRequest
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
	Unset     *SessionPolicyUnsetRequest
}

type SessionPolicySetRequest struct {
	SessionIdleTimeoutMins   *int
	SessionUiIdleTimeoutMins *int
	Comment                  *string
}

type SessionPolicyUnsetRequest struct {
	SessionIdleTimeoutMins   *bool
	SessionUiIdleTimeoutMins *bool
	Comment                  *bool
}

type DropSessionPolicyRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowSessionPolicyRequest struct{}

type DescribeSessionPolicyRequest struct {
	name SchemaObjectIdentifier // required
}
