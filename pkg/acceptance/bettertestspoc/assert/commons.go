package assert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestCheckFuncProvider is an interface with just one method providing resource.TestCheckFunc.
// It allows using it as input the "Check:" in resource.TestStep.
// It should be used with AssertThat.
type TestCheckFuncProvider interface {
	ToTerraformTestCheckFunc(t *testing.T) resource.TestCheckFunc
}

// AssertThat should be used for "Check:" input in resource.TestStep instead of e.g. resource.ComposeTestCheckFunc.
// It allows performing all the checks implementing the TestCheckFuncProvider interface.
func AssertThat(t *testing.T, fs ...TestCheckFuncProvider) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		var result []error

		for i, f := range fs {
			if err := f.ToTerraformTestCheckFunc(t)(s); err != nil {
				result = append(result, fmt.Errorf("check %d/%d error:\n%w", i+1, len(fs), err))
			}
		}

		return errors.Join(result...)
	}
}

var _ TestCheckFuncProvider = (*testCheckFuncWrapper)(nil)

type testCheckFuncWrapper struct {
	f resource.TestCheckFunc
}

func (w *testCheckFuncWrapper) ToTerraformTestCheckFunc(_ *testing.T) resource.TestCheckFunc {
	return w.f
}

// Check allows using the basic terraform checks while using AssertThat.
// To use, just simply wrap the check in Check.
func Check(f resource.TestCheckFunc) TestCheckFuncProvider {
	return &testCheckFuncWrapper{f}
}

// ImportStateCheckFuncProvider is an interface with just one method providing resource.ImportStateCheckFunc.
// It allows using it as input the "ImportStateCheck:" in resource.TestStep for import tests.
// It should be used with AssertThatImport.
type ImportStateCheckFuncProvider interface {
	ToTerraformImportStateCheckFunc(t *testing.T) resource.ImportStateCheckFunc
}

// AssertThatImport should be used for "ImportStateCheck:" input in resource.TestStep instead of e.g. importchecks.ComposeImportStateCheck.
// It allows performing all the checks implementing the ImportStateCheckFuncProvider interface.
func AssertThatImport(t *testing.T, fs ...ImportStateCheckFuncProvider) resource.ImportStateCheckFunc {
	t.Helper()
	return func(s []*terraform.InstanceState) error {
		var result []error

		for i, f := range fs {
			if err := f.ToTerraformImportStateCheckFunc(t)(s); err != nil {
				result = append(result, fmt.Errorf("check %d/%d error:\n%w", i+1, len(fs), err))
			}
		}

		return errors.Join(result...)
	}
}

var _ ImportStateCheckFuncProvider = (*importStateCheckFuncWrapper)(nil)

type importStateCheckFuncWrapper struct {
	f resource.ImportStateCheckFunc
}

func (w *importStateCheckFuncWrapper) ToTerraformImportStateCheckFunc(_ *testing.T) resource.ImportStateCheckFunc {
	return w.f
}

// CheckImport allows using the basic terraform import checks while using AssertThatImport.
// To use, just simply wrap the check in CheckImport.
func CheckImport(f resource.ImportStateCheckFunc) ImportStateCheckFuncProvider {
	return &importStateCheckFuncWrapper{f}
}

// InPlaceAssertionVerifier is an interface providing a method allowing verifying all the prepared assertions in place.
// It does not return function like TestCheckFuncProvider or ImportStateCheckFuncProvider; it runs all the assertions in place instead.
type InPlaceAssertionVerifier interface {
	VerifyAll(t *testing.T)
}

// AssertThatObject should be used in the SDK tests for created object validation.
// It verifies all the prepared assertions in place.
func AssertThatObject(t *testing.T, objectAssert InPlaceAssertionVerifier) {
	objectAssert.VerifyAll(t)
}
