package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ObjectParameter(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: objectParameterBasic(prefix, "ENABLE_STREAM_TASK_REPLICATION", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", "ENABLE_STREAM_TASK_REPLICATION"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "true"),
				),
			},
		},
	})
}

func objectParameterBasic(prefix, key, value string) string {
	s := `
resource "snowflake_database" "d" {
	name = "%s"
}	
resource "snowflake_object_parameter" "p" {
	key = "%s"
	value = "%s"
	object_type = "DATABASE"
	object_name = snowflake_database.d.name
}
`
	return fmt.Sprintf(s,prefix, key, value)
}
