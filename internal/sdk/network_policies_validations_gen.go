// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import "errors"

var (
	_ validatable = new(CreateNetworkPolicyOptions)
	_ validatable = new(AlterNetworkPolicyOptions)
	_ validatable = new(DropNetworkPolicyOptions)
	_ validatable = new(ShowNetworkPolicyOptions)
	_ validatable = new(DescribeNetworkPolicyOptions)
)

func (opts *CreateNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *AlterNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.Set, opts.UnsetComment, opts.RenameTo); !ok {
		errs = append(errs, errExactlyOneOf("Set", "UnsetComment", "RenameTo"))
	}
	if valueSet(opts.RenameTo) && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Set) {
		if ok := anyValueSet(opts.Set.AllowedIpList, opts.Set.BlockedIpList, opts.Set.Comment); !ok {
			errs = append(errs, errAtLeastOneOf("AllowedIpList", "BlockedIpList", "Comment"))
		}
	}
	return errors.Join(errs...)
}

func (opts *DropNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *ShowNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *DescribeNetworkPolicyOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
