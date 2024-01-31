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
		if !exactlyOneValueSet(opts.AutomatedDataLoadsParams.GoogleAutoParams, opts.AutomatedDataLoadsParams.AzureAutoParams) {
			errs = append(errs, errExactlyOneOf("CreateNotificationIntegrationOptions.AutomatedDataLoadsParams", "GoogleAutoParams", "AzureAutoParams"))
		}
	}
	if valueSet(opts.PushNotificationParams) {
		if !exactlyOneValueSet(opts.PushNotificationParams.AmazonPushParams, opts.PushNotificationParams.GooglePushParams, opts.PushNotificationParams.AzurePushParams) {
			errs = append(errs, errExactlyOneOf("CreateNotificationIntegrationOptions.PushNotificationParams", "AmazonPushParams", "GooglePushParams", "AzurePushParams"))
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
	if !exactlyOneValueSet(opts.Set, opts.UnsetEmailParams, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterNotificationIntegrationOptions", "Set", "UnsetEmailParams", "SetTags", "UnsetTags"))
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
	if valueSet(opts.UnsetEmailParams) {
		if !anyValueSet(opts.UnsetEmailParams.AllowedRecipients, opts.UnsetEmailParams.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterNotificationIntegrationOptions.UnsetEmailParams", "AllowedRecipients", "Comment"))
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
