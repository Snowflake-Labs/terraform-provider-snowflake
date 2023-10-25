// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Comment(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("set", func(t *testing.T) {
		comment := random.Comment()
		err := client.Comments.Set(ctx, &sdk.SetCommentOptions{
			ObjectType: sdk.ObjectTypeWarehouse,
			ObjectName: testWarehouse(t).ID(),
			Value:      sdk.String(comment),
		})
		require.NoError(t, err)
		wh, err := client.Warehouses.ShowByID(ctx, testWarehouse(t).ID())
		require.NoError(t, err)
		assert.Equal(t, comment, wh.Comment)
	})

	// TODO: uncomment once we can create tables/columns
	// t.Run("set column", func(t *testing.T) {
	// 	comment := randomComment(t)
	// 	err := client.Comments.SetColumn(ctx, &SetColumnCommentOpts{
	// 		Column: testWarehouse.ID(),
	// 		Value:  String(comment),
	// 	})
	// 	require.NoError(t, err)
	// 	wh, err := client.Warehouses.ShowByID(ctx, testWarehouse.ID())
	// 	require.NoError(t, err)
	// 	assert.Equal(t, comment, wh.Comment)
	// })
}
