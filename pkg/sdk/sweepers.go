package sdk

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
)

func SweepAfterIntegrationTests(client *Client, suffix string) error {
	return sweep(client, suffix)
}

func SweepAfterAcceptanceTests(client *Client, suffix string) error {
	return sweep(client, suffix)
}

// TODO [SNOW-867247]: move this to test code
// TODO [SNOW-867247]: use if exists/use method from helper for dropping
// TODO [SNOW-867247]: sweep all missing account-level objects (like users, integrations, replication groups, network policies, ...)
// TODO [SNOW-867247]: extract sweepers to a separate dir
// TODO [SNOW-867247]: rework the sweepers (funcs -> objects)
// TODO [SNOW-867247]: consider generalization (almost all the sweepers follow the same pattern: show, drop if matches)
// TODO [SNOW-867247]: consider failing after all sweepers and not with the first error
// TODO [SNOW-867247]: consider showing only objects with the given suffix (in almost every sweeper)
func sweep(client *Client, suffix string) error {
	if suffix == "" {
		return fmt.Errorf("suffix is required to run sweepers")
	}
	sweepers := []func() error{
		getAccountPolicyAttachmentsSweeper(client),
		getResourceMonitorSweeper(client, suffix),
		getFailoverGroupSweeper(client, suffix),
		getShareSweeper(client, suffix),
		getDatabaseSweeper(client, suffix),
		getWarehouseSweeper(client, suffix),
		getRoleSweeper(client, suffix),
	}
	for _, sweeper := range sweepers {
		if err := sweeper(); err != nil {
			return err
		}
	}
	return nil
}

func getAccountPolicyAttachmentsSweeper(client *Client) func() error {
	return func() error {
		log.Printf("[DEBUG] Unsetting password and session policies set on the account level")
		ctx := context.Background()
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				PasswordPolicy: Bool(true),
			},
		}
		_ = client.Accounts.Alter(ctx, opts)
		opts = &AlterAccountOptions{
			Unset: &AccountUnset{
				SessionPolicy: Bool(true),
			},
		}
		_ = client.Accounts.Alter(ctx, opts)
		return nil
	}
}

func getResourceMonitorSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping resource monitors with suffix %s", suffix)
		ctx := context.Background()

		rms, err := client.ResourceMonitors.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping resource monitor ended with error, err = %w", err)
		}
		for _, rm := range rms {
			if strings.HasSuffix(rm.Name, suffix) {
				log.Printf("[DEBUG] Dropping resource monitor %s", rm.ID().FullyQualifiedName())
				if err := client.ResourceMonitors.Drop(ctx, rm.ID(), &DropResourceMonitorOptions{IfExists: Bool(true)}); err != nil {
					return fmt.Errorf("sweeping resource monitor %s ended with error, err = %w", rm.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping resource monitor %s", rm.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getFailoverGroupSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping failover groups with suffix %s", suffix)
		ctx := context.Background()

		currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
		if err != nil {
			return fmt.Errorf("sweeping failover groups ended with error, err = %w", err)
		}
		opts := &ShowFailoverGroupOptions{
			InAccount: NewAccountIdentifierFromAccountLocator(currentAccount),
		}
		fgs, err := client.FailoverGroups.Show(ctx, opts)
		if err != nil {
			return fmt.Errorf("sweeping failover groups ended with error, err = %w", err)
		}
		for _, fg := range fgs {
			if strings.HasSuffix(fg.Name, suffix) && fg.AccountLocator == currentAccount {
				log.Printf("[DEBUG] Dropping failover group %s", fg.ID().FullyQualifiedName())
				if err := client.FailoverGroups.Drop(ctx, fg.ID(), nil); err != nil {
					return fmt.Errorf("sweeping failover group %s ended with error, err = %w", fg.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping failover group %s", fg.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getShareSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping shares with suffix %s", suffix)
		ctx := context.Background()

		shares, err := client.Shares.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping shares ended with error, err = %w", err)
		}
		for _, share := range shares {
			if share.Kind == ShareKindOutbound && strings.HasSuffix(share.Name.Name(), suffix) {
				log.Printf("[DEBUG] Dropping share %s", share.ID().FullyQualifiedName())
				if err := client.Shares.Drop(ctx, share.ID(), &DropShareOptions{IfExists: Bool(true)}); err != nil {
					return fmt.Errorf("sweeping share %s ended with error, err = %w", share.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping share %s", share.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getDatabaseSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping databases with suffix %s", suffix)
		ctx := context.Background()

		dbs, err := client.Databases.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping databases ended with error, err = %w", err)
		}
		for _, db := range dbs {
			if strings.HasSuffix(db.Name, suffix) && db.Name != "SNOWFLAKE" {
				log.Printf("[DEBUG] Dropping database %s", db.ID().FullyQualifiedName())
				if err := client.Databases.Drop(ctx, db.ID(), nil); err != nil {
					if strings.Contains(err.Error(), "Object found is of type 'APPLICATION', not specified type 'DATABASE'") {
						log.Printf("[DEBUG] Skipping database %s", db.ID().FullyQualifiedName())
					} else {
						return fmt.Errorf("sweeping database %s ended with error, err = %w", db.ID().FullyQualifiedName(), err)
					}
				}
			} else {
				log.Printf("[DEBUG] Skipping database %s", db.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getWarehouseSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping warehouses with suffix %s", suffix)
		ctx := context.Background()

		whs, err := client.Warehouses.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping warehouses ended with error, err = %w", err)
		}
		for _, wh := range whs {
			if strings.HasSuffix(wh.Name, suffix) && wh.Name != "SNOWFLAKE" {
				log.Printf("[DEBUG] Dropping warehouse %s", wh.ID().FullyQualifiedName())
				if err := client.Warehouses.Drop(ctx, wh.ID(), nil); err != nil {
					return fmt.Errorf("sweeping warehouse %s ended with error, err = %w", wh.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping warehouse %s", wh.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getRoleSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping roles with suffix %s", suffix)
		ctx := context.Background()

		roles, err := client.Roles.Show(ctx, NewShowRoleRequest())
		if err != nil {
			return fmt.Errorf("sweeping roles ended with error, err = %w", err)
		}
		for _, role := range roles {
			if strings.HasSuffix(role.Name, suffix) && !slices.Contains([]string{"ACCOUNTADMIN", "SECURITYADMIN", "SYSADMIN", "ORGADMIN", "USERADMIN", "PUBLIC"}, role.Name) {
				log.Printf("[DEBUG] Dropping role %s", role.ID().FullyQualifiedName())
				if err := client.Roles.Drop(ctx, NewDropRoleRequest(role.ID())); err != nil {
					return fmt.Errorf("sweeping role %s ended with error, err = %w", role.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping role %s", role.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}
