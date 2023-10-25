name: Feature Request
description: Something is missing or could be improved.
labels: ["feature-request"]
body:
  - type: markdown
    attributes:
      value: |
        Thank you for taking the time to fill out this feature request! Please note that this issue tracker is only used for bug reports and feature requests. Other issues will be closed.

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
    id: use-case
    attributes:
      label: Use Cases or Problem Statement
      description: What use cases or problems are you trying to solve?
      placeholder: Description of use cases or problems.
    validations:
      required: true
  - type: textarea
    id: proposal
    attributes:
      label: Proposal
      description: What solutions would you prefer?
      placeholder: Description of proposed solutions.
    validations:
      required: true
  - type: dropdown
    id: impact
    attributes:
      label: How much impact is this issue causing?
      description: High represents completely not able to use the provider without this. Medium represents unable to solve a specific problem or understand something. Low represents minor provider code, configuration, or documentation issues.
      options:
        - High
        - Medium
        - Low
    validations:
      required: true
  - type: textarea
    id: additional-information
    attributes:
      label: Additional Information
      description: Are there any additional details about your environment, workflow, or recent changes that might be relevant? Have you discovered a workaround? Are there links to other related issues?
    validations:
      required: false
