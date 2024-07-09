// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateAuthenticationPolicyRequest(
	name AccountObjectIdentifier,
) *CreateAuthenticationPolicyRequest {
	s := CreateAuthenticationPolicyRequest{}
	s.name = name
	return &s
}

func (s *CreateAuthenticationPolicyRequest) WithOrReplace(OrReplace bool) *CreateAuthenticationPolicyRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateAuthenticationPolicyRequest) WithAuthenticationMethods(AuthenticationMethods []SchemaObjectIdentifier) *CreateAuthenticationPolicyRequest {
	s.AuthenticationMethods = AuthenticationMethods
	return s
}

func (s *CreateAuthenticationPolicyRequest) WithMfaAuthenticationMethods(MfaAuthenticationMethods []SchemaObjectIdentifier) *CreateAuthenticationPolicyRequest {
	s.MfaAuthenticationMethods = MfaAuthenticationMethods
	return s
}

func (s *CreateAuthenticationPolicyRequest) WithMfaEnrollment(MfaEnrollment string) *CreateAuthenticationPolicyRequest {
	s.MfaEnrollment = &MfaEnrollment
	return s
}

func (s *CreateAuthenticationPolicyRequest) WithClientTypes(ClientTypes []SchemaObjectIdentifier) *CreateAuthenticationPolicyRequest {
	s.ClientTypes = ClientTypes
	return s
}

func (s *CreateAuthenticationPolicyRequest) WithSecurityIntegrations(SecurityIntegrations []SchemaObjectIdentifier) *CreateAuthenticationPolicyRequest {
	s.SecurityIntegrations = SecurityIntegrations
	return s
}

func (s *CreateAuthenticationPolicyRequest) WithComment(Comment string) *CreateAuthenticationPolicyRequest {
	s.Comment = &Comment
	return s
}

func NewAlterAuthenticationPolicyRequest(
	name AccountObjectIdentifier,
) *AlterAuthenticationPolicyRequest {
	s := AlterAuthenticationPolicyRequest{}
	s.name = name
	return &s
}

func (s *AlterAuthenticationPolicyRequest) WithIfExists(IfExists bool) *AlterAuthenticationPolicyRequest {
	s.IfExists = &IfExists
	return s
}

func (s *AlterAuthenticationPolicyRequest) WithSet(Set AuthenticationPolicySetRequest) *AlterAuthenticationPolicyRequest {
	s.Set = &Set
	return s
}

func (s *AlterAuthenticationPolicyRequest) WithUnset(Unset AuthenticationPolicyUnsetRequest) *AlterAuthenticationPolicyRequest {
	s.Unset = &Unset
	return s
}

func (s *AlterAuthenticationPolicyRequest) WithRenameTo(RenameTo AccountObjectIdentifier) *AlterAuthenticationPolicyRequest {
	s.RenameTo = &RenameTo
	return s
}

func NewAuthenticationPolicySetRequest() *AuthenticationPolicySetRequest {
	return &AuthenticationPolicySetRequest{}
}

func (s *AuthenticationPolicySetRequest) WithAuthenticationMethods(AuthenticationMethods []SchemaObjectIdentifier) *AuthenticationPolicySetRequest {
	s.AuthenticationMethods = AuthenticationMethods
	return s
}

func (s *AuthenticationPolicySetRequest) WithMfaAuthenticationMethods(MfaAuthenticationMethods []SchemaObjectIdentifier) *AuthenticationPolicySetRequest {
	s.MfaAuthenticationMethods = MfaAuthenticationMethods
	return s
}

func (s *AuthenticationPolicySetRequest) WithMfaEnrollment(MfaEnrollment []SchemaObjectIdentifier) *AuthenticationPolicySetRequest {
	s.MfaEnrollment = MfaEnrollment
	return s
}

func (s *AuthenticationPolicySetRequest) WithClientTypes(ClientTypes []SchemaObjectIdentifier) *AuthenticationPolicySetRequest {
	s.ClientTypes = ClientTypes
	return s
}

func (s *AuthenticationPolicySetRequest) WithSecurityIntegrations(SecurityIntegrations []SchemaObjectIdentifier) *AuthenticationPolicySetRequest {
	s.SecurityIntegrations = SecurityIntegrations
	return s
}

func (s *AuthenticationPolicySetRequest) WithComment(Comment string) *AuthenticationPolicySetRequest {
	s.Comment = &Comment
	return s
}

func NewAuthenticationPolicyUnsetRequest() *AuthenticationPolicyUnsetRequest {
	return &AuthenticationPolicyUnsetRequest{}
}

func (s *AuthenticationPolicyUnsetRequest) WithClientTypes(ClientTypes bool) *AuthenticationPolicyUnsetRequest {
	s.ClientTypes = &ClientTypes
	return s
}

func (s *AuthenticationPolicyUnsetRequest) WithAuthenticationMethods(AuthenticationMethods bool) *AuthenticationPolicyUnsetRequest {
	s.AuthenticationMethods = &AuthenticationMethods
	return s
}

func (s *AuthenticationPolicyUnsetRequest) WithSecurityIntegrations(SecurityIntegrations bool) *AuthenticationPolicyUnsetRequest {
	s.SecurityIntegrations = &SecurityIntegrations
	return s
}

func (s *AuthenticationPolicyUnsetRequest) WithMfaAuthenticationMethods(MfaAuthenticationMethods bool) *AuthenticationPolicyUnsetRequest {
	s.MfaAuthenticationMethods = &MfaAuthenticationMethods
	return s
}

func (s *AuthenticationPolicyUnsetRequest) WithMfaEnrollment(MfaEnrollment bool) *AuthenticationPolicyUnsetRequest {
	s.MfaEnrollment = &MfaEnrollment
	return s
}

func (s *AuthenticationPolicyUnsetRequest) WithComment(Comment bool) *AuthenticationPolicyUnsetRequest {
	s.Comment = &Comment
	return s
}

func NewDropAuthenticationPolicyRequest(
	name AccountObjectIdentifier,
) *DropAuthenticationPolicyRequest {
	s := DropAuthenticationPolicyRequest{}
	s.name = name
	return &s
}

func (s *DropAuthenticationPolicyRequest) WithIfExists(IfExists bool) *DropAuthenticationPolicyRequest {
	s.IfExists = &IfExists
	return s
}

func NewShowAuthenticationPolicyRequest() *ShowAuthenticationPolicyRequest {
	return &ShowAuthenticationPolicyRequest{}
}

func NewDescribeAuthenticationPolicyRequest(
	name AccountObjectIdentifier,
) *DescribeAuthenticationPolicyRequest {
	s := DescribeAuthenticationPolicyRequest{}
	s.name = name
	return &s
}
