package sdk

import (
	"context"
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

func Sweep(client *Client, prefix string) error {
	sweepers := []func() error{
		getFailoverGroupSweeper(client, prefix),
		getShareSweeper(client, prefix),
		getDatabaseSweeper(client, prefix),
		getWarehouseSweeper(client, prefix),
		getRoleSweeper(client, prefix),
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
		opts := &ShowFailoverGroupOptions{
			InAccount: NewAccountIdentifierFromAccountLocator(currentAccount),
		}
		fgs, err := client.FailoverGroups.Show(ctx, opts)
		if err != nil {
			return err
		}
		for _, fg := range fgs {
			if (prefix == "" || strings.HasPrefix(fg.Name, prefix)) && fg.AccountLocator == currentAccount {
				log.Printf("[DEBUG] Dropping failover group %s", fg.Name)
				if err := client.FailoverGroups.Drop(ctx, fg.ID(), nil); err != nil {
					return err
				}
			} else {
				log.Printf("[DEBUG] Skipping failover group %s", fg.Name)
			}
		}
		return nil
	}
}

func getRoleSweeper(client *Client, prefix string) func() error {
	return func() error {
		if prefix == "" {
			log.Printf("[DEBUG] Sweeping all roles")
		} else {
			log.Printf("[DEBUG] Sweeping all roles with prefix %s", prefix)
		}
		ctx := context.Background()
		roles, err := client.Roles.Show(ctx, nil)
		if err != nil {
			return err
		}
		for _, role := range roles {
			if (prefix == "" || strings.HasPrefix(role.Name, prefix)) && !slices.Contains([]string{"ACCOUNTADMIN", "SECURITYADMIN", "SYSADMIN", "ORGADMIN", "USERADMIN", "PUBLIC"}, role.Name) {
				log.Printf("[DEBUG] Dropping role %s", role.Name)
				if err := client.Roles.Drop(ctx, role.ID(), nil); err != nil {
					return err
				}
			} else {
				log.Printf("[DEBUG] Skipping role %s", role.Name)
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
			if (share.Kind == ShareKindOutbound) && (prefix == "" || strings.HasPrefix(share.Name.Name(), prefix)) {
				log.Printf("[DEBUG] Dropping share %s", share.Name.Name())
				if err := client.Shares.Drop(ctx, share.ID()); err != nil {
					return err
				}
			} else {
				log.Printf("[DEBUG] Skipping share %s", share.Name.Name())
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
			if (prefix == "" || strings.HasPrefix(db.Name, prefix)) && db.Name != "SNOWFLAKE" {
				log.Printf("[DEBUG] Dropping database %s", db.Name)
				if err := client.Databases.Drop(ctx, db.ID(), nil); err != nil {
					return err
				}
			} else {
				log.Printf("[DEBUG] Skipping database %s", db.Name)
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
			if (prefix == "" || strings.HasPrefix(wh.Name, prefix)) && wh.Name != "SNOWFLAKE" {
				log.Printf("[DEBUG] Dropping warehouse %s", wh.Name)
				if err := client.Warehouses.Drop(ctx, wh.ID(), nil); err != nil {
					return err
				}
			} else {
				log.Printf("[DEBUG] Skipping warehouse %s", wh.Name)
			}
		}
		return nil
	}
}
