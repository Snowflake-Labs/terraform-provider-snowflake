package generator

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/gen/builder"
)

func GenerateAll(api any, sb ...*builder.StructBuilder) {
	for _, b := range sb {
		fmt.Printf("%s", b.String())
	}
}
