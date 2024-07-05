package bettertestspoc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type TestCheckFuncProvider interface {
	ToTerraformTestCheckFunc(t *testing.T) resource.TestCheckFunc
}

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

func Check(f resource.TestCheckFunc) TestCheckFuncProvider {
	return &testCheckFuncWrapper{f}
}

type ImportStateCheckFuncProvider interface {
	ToTerraformImportStateCheckFunc(t *testing.T) resource.ImportStateCheckFunc
}

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

func CheckImport(f resource.ImportStateCheckFunc) ImportStateCheckFuncProvider {
	return &importStateCheckFuncWrapper{f}
}
