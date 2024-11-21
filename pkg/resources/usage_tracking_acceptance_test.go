package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/tracking"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_CompleteUsageTracking(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	userModel := model.User("test", id.Name())
	userModelWithComment := model.User("test", id.Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// Create + CustomDiff (parameters)
			{
				Config: config.FromModel(t, userModel),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(id.Name()).
						HasNoComment(),
				),
			},
			// Import
			{
				ResourceName: userModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedUserResource(t, id.Name()).
						HasNoComment(),
				),
			},
			// Update
			{
				Config: config.FromModel(t, userModelWithComment),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModelWithComment.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(comment),
				),
			},
		},
	}) // Delete

	assertQueryMetadataExists := func(t *testing.T, queryHistory []helpers.QueryHistory, operation tracking.Operation, query string) {
		t.Helper()
		found := false
		expectedMetadata := tracking.NewVersionedMetadata(resources.User, operation)
		for _, history := range queryHistory {
			if metadata, err := tracking.ParseMetadata(history.QueryText); err == nil {
				if expectedMetadata == metadata && strings.Contains(history.QueryText, query) {
					found = true
				}
			}
		}
		if !found {
			t.Fatalf("query history does not contain query metadata: %v with query containing: %s", expectedMetadata, query)
		}
	}

	queryHistory := acc.TestClient().InformationSchema.GetQueryHistory(t, 100)
	assertQueryMetadataExists(t, queryHistory, tracking.CreateOperation, fmt.Sprintf(`CREATE USER "%s"`, id.Name()))
	assertQueryMetadataExists(t, queryHistory, tracking.ReadOperation, fmt.Sprintf(`SHOW USERS LIKE '%s'`, id.Name()))
	assertQueryMetadataExists(t, queryHistory, tracking.UpdateOperation, fmt.Sprintf(`ALTER USER "%s" SET COMMENT = '%s'`, id.Name(), comment))
	assertQueryMetadataExists(t, queryHistory, tracking.DeleteOperation, fmt.Sprintf(`DROP USER IF EXISTS "%s"`, id.Name()))
	assertQueryMetadataExists(t, queryHistory, tracking.CustomDiffOperation, fmt.Sprintf(`SHOW PARAMETERS IN USER "%s"`, id.Name()))
	assertQueryMetadataExists(t, queryHistory, tracking.ImportOperation, fmt.Sprintf(`DESCRIBE USER "%s"`, id.Name()))
}
