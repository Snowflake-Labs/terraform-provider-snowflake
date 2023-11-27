package sdk

import "errors"

var (
	_ validatable = new(CreateEventTableOptions)
	_ validatable = new(ShowEventTableOptions)
	_ validatable = new(DescribeEventTableOptions)
	_ validatable = new(AlterEventTableOptions)
)

func (opts *CreateEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) && *opts.OrReplace && *opts.IfNotExists {
		errs = append(errs, errOneOf("OrReplace", "IfNotExists"))
	}
	if valueSet(opts.CopyGrants) && !valueSet(opts.OrReplace) {
		errs = append(errs, errors.New("CopyGrants requires OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *ShowEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (v *EventTableClusteringAction) validate() error {
	if v == nil {
		return ErrNilOptions
	}
	var errs []error
	if ok := anyValueSet(
		v.ClusterBy,
		v.ResumeRecluster,
		v.SuspendRecluster,
		v.DropClusteringKey,
	); !ok {
		errs = append(errs, errAtLeastOneOf("ClusterBy", "ResumeRecluster", "SuspendRecluster", "DropClusteringKey"))
	}
	if ok := exactlyOneValueSet(
		v.ClusterBy,
		v.ResumeRecluster,
		v.SuspendRecluster,
		v.DropClusteringKey,
	); !ok {
		errs = append(errs, errExactlyOneOf("ClusterBy", "ResumeRecluster", "SuspendRecluster", "DropClusteringKey"))
	}
	return errors.Join(errs...)
}

func (v *EventTableSearchOptimizationAction) validate() error {
	if v == nil {
		return ErrNilOptions
	}
	var errs []error
	if ok := anyValueSet(
		v.Add,
		v.Drop,
	); !ok {
		errs = append(errs, errAtLeastOneOf("Add", "Drop"))
	}
	if ok := exactlyOneValueSet(
		v.Add,
		v.Drop,
	); !ok {
		errs = append(errs, errExactlyOneOf("Add", "Drop"))
	}
	return errors.Join(errs...)
}

func (v *EventTableSet) validate() error {
	if v == nil {
		return ErrNilOptions
	}
	var errs []error
	if ok := anyValueSet(
		v.DataRetentionTimeInDays,
		v.MaxDataExtensionTimeInDays,
		v.ChangeTracking,
		v.Comment,
	); !ok {
		errs = append(errs, errAtLeastOneOf("DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "Comment"))
	}
	if ok := exactlyOneValueSet(
		v.DataRetentionTimeInDays,
		v.MaxDataExtensionTimeInDays,
		v.ChangeTracking,
		v.Comment,
	); !ok {
		errs = append(errs, errExactlyOneOf("DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "Comment"))
	}
	return errors.Join(errs...)
}

func (v *EventTableUnset) validate() error {
	if v == nil {
		return ErrNilOptions
	}
	var errs []error
	if ok := anyValueSet(
		v.DataRetentionTimeInDays,
		v.MaxDataExtensionTimeInDays,
		v.ChangeTracking,
		v.Comment,
	); !ok {
		errs = append(errs, errAtLeastOneOf("DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "Comment"))
	}
	if ok := exactlyOneValueSet(
		v.DataRetentionTimeInDays,
		v.MaxDataExtensionTimeInDays,
		v.ChangeTracking,
		v.Comment,
	); !ok {
		errs = append(errs, errExactlyOneOf("DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "Comment"))
	}
	return errors.Join(errs...)
}

func (opts *DescribeEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterEventTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := anyValueSet(
		opts.Set,
		opts.Unset,
		opts.SetTags,
		opts.UnsetTags,
		opts.AddRowAccessPolicy,
		opts.DropRowAccessPolicy,
		opts.DropAllRowAccessPolicies,
		opts.SearchOptimizationAction,
		opts.ClusteringAction,
		opts.RenameTo,
	); !ok {
		errs = append(errs, errAtLeastOneOf("Set", "Unset", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAllRowAccessPolicies", "SearchOptimizationAction", "ClusteringAction", "RenameTo"))
	}
	if ok := exactlyOneValueSet(
		opts.Set,
		opts.Unset,
		opts.SetTags,
		opts.UnsetTags,
		opts.AddRowAccessPolicy,
		opts.DropRowAccessPolicy,
		opts.DropAllRowAccessPolicies,
		opts.SearchOptimizationAction,
		opts.ClusteringAction,
		opts.RenameTo,
	); !ok {
		errs = append(errs, errExactlyOneOf("Set", "Unset", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAllRowAccessPolicies", "SearchOptimizationAction", "ClusteringAction", "RenameTo"))
	}
	if valueSet(opts.ClusteringAction) {
		if err := opts.ClusteringAction.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.SearchOptimizationAction) {
		if err := opts.SearchOptimizationAction.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return JoinErrors(errs...)
}
