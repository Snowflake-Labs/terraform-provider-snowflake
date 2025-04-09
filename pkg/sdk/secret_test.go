package sdk

import (
	"fmt"
	"log"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/stretchr/testify/assert"
)

func Test_TestSecret(t *testing.T) {
	t.Logf("t.Log: %s", random.Password())
	fmt.Printf("fmt.Printf: %s\n", random.Password())
	log.Printf("log.Printf: %s", random.Password())
	assert.True(t, false, fmt.Sprintf("Test failed: %s", random.Password()))
	assert.Equal(t, random.Password(), random.Password())
}
