package resources

import (
	"context"
	"errors"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceValueSetter interface {
	Set(string, any) error
}

type DropFunc = func(context.Context, sdk.SchemaObjectIdentifier) error

func CommonDelete(dropFunc func(*sdk.Client) DropFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

		err := dropFunc(client)(ctx, id)
		if !errors.Is(err, sdk.ErrSkippable) {
			return diag.FromErr(err)
		}

		d.SetId("")
		return nil
	}
}
