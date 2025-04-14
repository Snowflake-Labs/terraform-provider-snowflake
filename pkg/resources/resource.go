package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceValueSetter interface {
	Set(string, any) error
}

type DropSafelyFunc[T sdk.ObjectIdentifierConstraint] func(context.Context, T) error

func ResourceDeleteContextFunc[ID sdk.ObjectIdentifierConstraint](
	parseFunc func(string) (ID, error),
	dropFunc func(*sdk.Client) DropSafelyFunc[ID],
) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := parseFunc(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		err = dropFunc(client)(ctx, id)
		if !errors.Is(err, sdk.ErrSkippable) {
			return diag.FromErr(err)
		}

		d.SetId("")
		return nil
	}
}
