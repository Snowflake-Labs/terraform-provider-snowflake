package sdk

import "fmt"

var (
	_ validatable = new(CreateViewOptions)
	_ validatable = new(AlterViewOptions)
	_ validatable = new(DropViewOptions)
	_ validatable = new(ShowViewOptions)
	_ validatable = new(DescribeViewOptions)
)

func (opts *CreateViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateViewOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.RowAccessPolicy) {
		if !ValidObjectIdentifier(opts.RowAccessPolicy.RowAccessPolicy) {
			errs = append(errs, errInvalidIdentifier("CreateViewOptions", "RowAccessPolicy"))
		}
		if !valueSet(opts.RowAccessPolicy.On) {
			errs = append(errs, errNotSet("CreateViewOptions.RowAccessPolicy", "On"))
		}
	}
	if valueSet(opts.AggregationPolicy) {
		if !ValidObjectIdentifier(opts.AggregationPolicy.AggregationPolicy) {
			errs = append(errs, errInvalidIdentifier("CreateViewOptions", "AggregationPolicy"))
		}
	}
	if valueSet(opts.Columns) {
		for i, columnOption := range opts.Columns {
			if valueSet(columnOption.MaskingPolicy) {
				if !ValidObjectIdentifier(columnOption.MaskingPolicy.MaskingPolicy) {
					errs = append(errs, errInvalidIdentifier(fmt.Sprintf("CreateViewOptions.Columns[%d]", i), "MaskingPolicy"))
				}
			}
			if valueSet(columnOption.ProjectionPolicy) {
				if !ValidObjectIdentifier(columnOption.ProjectionPolicy.ProjectionPolicy) {
					errs = append(errs, errInvalidIdentifier(fmt.Sprintf("CreateViewOptions.Columns[%d]", i), "ProjectionPolicy"))
				}
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetSecure, opts.SetChangeTracking, opts.UnsetSecure, opts.SetTags, opts.UnsetTags, opts.AddDataMetricFunction, opts.DropDataMetricFunction, opts.SetDataMetricSchedule, opts.UnsetDataMetricSchedule, opts.AddRowAccessPolicy, opts.DropRowAccessPolicy, opts.DropAndAddRowAccessPolicy, opts.DropAllRowAccessPolicies, opts.SetAggregationPolicy, opts.UnsetAggregationPolicy, opts.SetMaskingPolicyOnColumn, opts.UnsetMaskingPolicyOnColumn, opts.SetProjectionPolicyOnColumn, opts.UnsetProjectionPolicyOnColumn, opts.SetTagsOnColumn, opts.UnsetTagsOnColumn) {
		errs = append(errs, errExactlyOneOf("AlterViewOptions", "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking", "UnsetSecure", "SetTags", "UnsetTags", "AddDataMetricFunction", "DropDataMetricFunction", "SetDataMetricSchedule", "UnsetDataMetricSchedule", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetProjectionPolicyOnColumn", "UnsetProjectionPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn"))
	}
	if everyValueSet(opts.IfExists, opts.SetSecure) {
		errs = append(errs, errOneOf("AlterViewOptions", "IfExists", "SetSecure"))
	}
	if everyValueSet(opts.IfExists, opts.UnsetSecure) {
		errs = append(errs, errOneOf("AlterViewOptions", "IfExists", "UnsetSecure"))
	}
	if valueSet(opts.SetDataMetricSchedule) {
		if !exactlyOneValueSet(opts.SetDataMetricSchedule.Minutes, opts.SetDataMetricSchedule.UsingCron, opts.SetDataMetricSchedule.TriggerOnChanges) {
			errs = append(errs, errExactlyOneOf("AlterViewOptions.SetDataMetricSchedule", "Minutes", "UsingCron", "TriggerOnChanges"))
		}
	}
	if valueSet(opts.AddRowAccessPolicy) {
		if !ValidObjectIdentifier(opts.AddRowAccessPolicy.RowAccessPolicy) {
			errs = append(errs, errInvalidIdentifier("AlterViewOptions.AddRowAccessPolicy", "RowAccessPolicy"))
		}
		if !valueSet(opts.AddRowAccessPolicy.On) {
			errs = append(errs, errNotSet("AlterViewOptions.AddRowAccessPolicy", "On"))
		}
	}
	if valueSet(opts.DropRowAccessPolicy) {
		if !ValidObjectIdentifier(opts.DropRowAccessPolicy.RowAccessPolicy) {
			errs = append(errs, errInvalidIdentifier("AlterViewOptions.DropRowAccessPolicy", "RowAccessPolicy"))
		}
	}
	if valueSet(opts.DropAndAddRowAccessPolicy) {
		if valueSet(opts.DropAndAddRowAccessPolicy.Drop) {
			if !ValidObjectIdentifier(opts.DropAndAddRowAccessPolicy.Drop.RowAccessPolicy) {
				errs = append(errs, errInvalidIdentifier("AlterViewOptions.DropAndAddRowAccessPolicy.Drop", "RowAccessPolicy"))
			}
		}
		if valueSet(opts.DropAndAddRowAccessPolicy.Add) {
			if !ValidObjectIdentifier(opts.DropAndAddRowAccessPolicy.Add.RowAccessPolicy) {
				errs = append(errs, errInvalidIdentifier("AlterViewOptions.DropAndAddRowAccessPolicy.Add", "RowAccessPolicy"))
			}
			if !valueSet(opts.DropAndAddRowAccessPolicy.Add.On) {
				errs = append(errs, errNotSet("AlterViewOptions.DropAndAddRowAccessPolicy.Add", "On"))
			}
		}
	}
	if valueSet(opts.SetAggregationPolicy) {
		if !ValidObjectIdentifier(opts.SetAggregationPolicy.AggregationPolicy) {
			errs = append(errs, errInvalidIdentifier("AlterViewOptions.SetAggregationPolicy", "AggregationPolicy"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
