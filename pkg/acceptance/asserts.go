package acceptance

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestCheckResourceAttrNumberAtLeast checks if specified field (key) in resource (name) is equal or greater than atLeast value.
func TestCheckResourceAttrNumberAtLeast(name string, key string, atLeast int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		is, err := getPrimaryModuleInstanceState(s, name)
		if err != nil {
			return err
		}

		v, ok := is.Attributes[key]
		if !ok {
			return fmt.Errorf("attribute %s not found", key)
		}

		actualValue, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("failed to parse attribute %s, err: %w", key, err)
		}

		if actualValue < atLeast {
			return fmt.Errorf("expected attribute %s to be at least %d, but was %d", key, atLeast, actualValue)
		}

		return nil
	}
}

func getPrimaryModuleInstanceState(s *terraform.State, name string) (*terraform.InstanceState, error) {
	ms := s.RootModule()

	rs, ok := ms.Resources[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s in %s", name, ms.Path)
	}

	is := rs.Primary
	if is == nil {
		return nil, fmt.Errorf("no primary instance: %s in %s", name, ms.Path)
	}

	return is, nil
}
