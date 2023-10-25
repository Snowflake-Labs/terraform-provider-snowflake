name: Bug Report
description: Something is incorrect or not working as expected.
labels: ["bug"]
body:
  - type: markdown
    attributes:
      value: |
        Thank you for taking the time to fill out this bug report! Please note that this issue tracker is only used for bug reports and feature requests. Other issues will be closed.

        If you have a configuration, workflow, or other question, please go back to the issue chooser and select one of the question links.
  - type: textarea
    id: versions
    attributes:
      label: Terraform CLI and Provider Versions
      description: What versions of Terraform CLI and the provider?
      placeholder: Output of `terraform version` from configuration directory
    validations:
      required: true
  - type: textarea
    id: terraform-configuration
    attributes:
      label: Terraform Configuration
      description: Please copy and paste any relevant Terraform configuration. This will be automatically formatted into code, so no need for backticks.
      render: terraform
    validations:
      required: true
  - type: textarea
    id: expected-behavior
    attributes:
      label: Expected Behavior
      description: What did you expect to happen?
      placeholder: Description of what should have happened.
    validations:
      required: true
  - type: textarea
    id: actual-behavior
    attributes:
      label: Actual Behavior
      description: What actually happened?
      placeholder: Description of what actually happened.
    validations:
      required: true
  - type: textarea
    id: reproduction-steps
    attributes:
      label: Steps to Reproduce
      description: List of steps to reproduce the issue.
      value: |
        1. `terraform apply`
    validations:
      required: true
  - type: dropdown
    id: impact
    attributes:
      label: How much impact is this issue causing?
      description: High represents completely not able to use the provider or unexpected destruction of data/infrastructure. Medium represents unable to upgrade provider version or an issue with potential workaround. Low represents minor provider code, configuration, or documentation issues.
      options:
        - High
        - Medium
        - Low
    validations:
      required: true
  - type: input
    id: logs
    attributes:
      label: Logs
      description: Please provide a link to a [GitHub Gist](https://gist.github.com) containing TRACE log output. [Terraform Debugging Documentation](https://www.terraform.io/internals/debugging)
      placeholder: https://gist.github.com/example/12345678
    validations:
      required: false
  - type: textarea
    id: additional-information
    attributes:
      label: Additional Information
      description: Are there any additional details about your environment, workflow, or recent changes that might be relevant? Have you discovered a workaround? Are there links to other related issues?
    validations:
      required: false
