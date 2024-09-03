# Manual tests

This directory is dedicated to hold steps for manual tests that are not possible to re-recreate in automated tests (or very hard to set up).
Every test should be placed in the subfolder representing a particular test (mostly because multiple files can be necessary for a single test)
and should contain a file describing the manual steps to perform the test.

Here's the list of cases we currently cannot reproduce and write acceptance tests for:
- `user_default_database_and_role`: Setting up a user with default_namespace and default_role, then logging into that user to see what happens with those values in various scenarios (e.g. insufficient privileges on the role).