package sdk

import (
	"context"
	"log"
	"strings"
)

func Sweep(client *Client, prefix string) error {
	sweepers := []func() error{
		getFailoverGroupSweeper(client, prefix),
		getShareSweeper(client, prefix),
		getDatabaseSweeper(client, prefix),
		getWarehouseSweeper(client, prefix),
	}
	for _, sweeper := range sweepers {
		if err := sweeper(); err != nil {
			return err
		}
	}
	return nil
}

func SweepAll(client *Client) error {
	return Sweep(client, "")
}

func getFailoverGroupSweeper(client *Client, prefix string) func() error {
	return func() error {
		if prefix == "" {
			log.Printf("[DEBUG] Sweeping all failover groups")
		} else {
			log.Printf("[DEBUG] Sweeping all failover groups with prefix %s", prefix)
		}
		ctx := context.Background()
		currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
		if err != nil {
			return err
		}
		opts := &FailoverGroupShowOptions{
			InAccount: NewAccountIdentifierFromAccountLocator(currentAccount),
		}
		fgs, err := client.FailoverGroups.Show(ctx, opts)
		if err != nil {
			return err
		}
		for _, fg := range fgs {
			if (prefix == "" || strings.HasPrefix(fg.Name, prefix)) && fg.AccountLocator == currentAccount {
				if err := client.FailoverGroups.Drop(ctx, fg.ID(), nil); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func getShareSweeper(client *Client, prefix string) func() error {
	return func() error {
		if prefix == "" {
			log.Printf("[DEBUG] Sweeping all shares")
		} else {
			log.Printf("[DEBUG] Sweeping all shares with prefix %s", prefix)
		}
		ctx := context.Background()
		shares, err := client.Shares.Show(ctx, nil)
		if err != nil {
			return err
		}
		for _, share := range shares {
			if share.Kind == ShareKindOutbound {
				if prefix == "" || strings.HasPrefix(share.Name.Name(), prefix) {
					if err := client.Shares.Drop(ctx, share.ID()); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

func getDatabaseSweeper(client *Client, prefix string) func() error {
	return func() error {
		if prefix == "" {
			log.Printf("[DEBUG] Sweeping all databases")
		} else {
			log.Printf("[DEBUG] Sweeping all databases with prefix %s", prefix)
		}
		ctx := context.Background()
		dbs, err := client.Databases.Show(ctx, nil)
		if err != nil {
			return err
		}
		for _, db := range dbs {
			if prefix == "" || strings.HasPrefix(db.Name, prefix) {
				if db.Name != "SNOWFLAKE" {
					if err := client.Databases.Drop(ctx, db.ID(), nil); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

func getWarehouseSweeper(client *Client, prefix string) func() error {
	return func() error {
		if prefix == "" {
			log.Printf("[DEBUG] Sweeping all warehouses")
		} else {
			log.Printf("[DEBUG] Sweeping all warehouses with prefix %s", prefix)
		}
		ctx := context.Background()
		whs, err := client.Warehouses.Show(ctx, nil)
		if err != nil {
			return err
		}
		for _, wh := range whs {
			if prefix == "" || strings.HasPrefix(wh.Name, prefix) {
				if err := client.Warehouses.Drop(ctx, wh.ID(), nil); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
