package sdk

var (
	_ validatable = new(CreateNotificationIntegrationOptions)
	_ validatable = new(AlterNotificationIntegrationOptions)
	_ validatable = new(DropNotificationIntegrationOptions)
	_ validatable = new(ShowNotificationIntegrationOptions)
	_ validatable = new(DescribeNotificationIntegrationOptions)
)

func (opts *CreateNotificationIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateNotificationIntegrationOptions", "IfNotExists", "OrReplace"))
	}
	if !exactlyOneValueSet(opts.AutomatedDataLoadsParams, opts.PushNotificationParams, opts.EmailParams) {
		errs = append(errs, errExactlyOneOf("CreateNotificationIntegrationOptions", "AutomatedDataLoadsParams", "PushNotificationParams", "EmailParams"))
	}
	if valueSet(opts.AutomatedDataLoadsParams) {
		if !exactlyOneValueSet(opts.AutomatedDataLoadsParams.GoogleAutomatedDataLoad, opts.AutomatedDataLoadsParams.AzureAutomatedDataLoad) {
			errs = append(errs, errExactlyOneOf("CreateNotificationIntegrationOptions.AutomatedDataLoadsParams", "GoogleAutomatedDataLoad", "AzureAutomatedDataLoad"))
		}
	}
	if valueSet(opts.PushNotificationParams) {
		if !exactlyOneValueSet(opts.PushNotificationParams.AmazonPush, opts.PushNotificationParams.GooglePush, opts.PushNotificationParams.AzurePush) {
			errs = append(errs, errExactlyOneOf("CreateNotificationIntegrationOptions.PushNotificationParams", "AmazonPush", "GooglePush", "AzurePush"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterNotificationIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterNotificationIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if everyValueSet(opts.Set.SetPushParams, opts.Set.SetEmailParams) {
			errs = append(errs, errOneOf("AlterNotificationIntegrationOptions.Set", "SetPushParams", "SetEmailParams"))
		}
		if !anyValueSet(opts.Set.Enabled, opts.Set.SetPushParams, opts.Set.SetEmailParams, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterNotificationIntegrationOptions.Set", "Enabled", "SetPushParams", "SetEmailParams", "Comment"))
		}
		if valueSet(opts.Set.SetPushParams) {
			if !exactlyOneValueSet(opts.Set.SetPushParams.SetAmazonPush, opts.Set.SetPushParams.SetGooglePush, opts.Set.SetPushParams.SetAzurePush) {
				errs = append(errs, errExactlyOneOf("AlterNotificationIntegrationOptions.Set.SetPushParams", "SetAmazonPush", "SetGooglePush", "SetAzurePush"))
			}
		}
		if valueSet(opts.Set.SetEmailParams) {
			if !valueSet(opts.Set.SetEmailParams.AllowedRecipients) {
				errs = append(errs, errNotSet("AlterNotificationIntegrationOptions.Set.SetEmailParams", "AllowedRecipients"))
			}
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.AllowedRecipients, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterNotificationIntegrationOptions.Unset", "AllowedRecipients", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropNotificationIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowNotificationIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeNotificationIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
