package provider_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform/helper/schema"
	_ "github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := provider.Provider()
	err := p.InternalValidate()
	a.NoError(err)
}

// func TestConfigureProvider(t *testing.T) {
// 	// a := assert.New(t)
// }

func TestDSN(t *testing.T) {
	t.Parallel()
	type args struct {
		s *schema.ResourceData
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"simple", args{resourceData(t, "acct", "user", "pass", "region", "role")},
			"user:pass@acct.region.snowflakecomputing.com:443?region=region&role=role", false},
		{"us-west-2 special case", args{resourceData(t, "acct2", "user2", "pass2", "us-west-2", "role2")},
			"user2:pass2@acct2.snowflakecomputing.com:443?role=role2", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.DSN(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("DSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func resourceData(t *testing.T, account, username, password, region, role string) *schema.ResourceData {
	a := assert.New(t)

	in := map[string]interface{}{
		"account":  account,
		"username": username,
		"password": password,
		"region":   region,
		"role":     role,
	}

	d := schema.TestResourceDataRaw(t, provider.Provider().Schema, in)
	a.NotNil(d)
	return d
}
