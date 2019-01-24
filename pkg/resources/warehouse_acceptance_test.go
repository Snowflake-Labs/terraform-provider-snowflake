package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccWarehouse(t *testing.T) {
	prefix := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: wConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", strings.ToUpper(prefix)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_warehouse.w", "warehouse_size"),
				),
			},
			// RENAME
			{
				Config: wConfig(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", strings.ToUpper(prefix2)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_warehouse.w", "warehouse_size"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: wConfig2(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "name", strings.ToUpper(prefix2)),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "comment", "test comment 2"),
					resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_size", "Small"),
				),
			},
		},
	})
}

func wConfig(prefix string) string {
	s := `
resource "snowflake_warehouse" "w" {
	name = "%s"
	comment = "test comment"
}
`
	return fmt.Sprintf(s, prefix)
}

func wConfig2(prefix string) string {
	s := `
resource "snowflake_warehouse" "w" {
	name = "%s"
	comment = "test comment 2"
	warehouse_size = "small"
}
`
	return fmt.Sprintf(s, prefix)
}
