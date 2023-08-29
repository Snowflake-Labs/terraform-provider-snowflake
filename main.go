// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/hashicorp/terraform-plugin-go/tfprotov6"
    "github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
    "github.com/hashicorp/terraform-plugin-mux/tf5to6server"
    "github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

    "example.com/terraform-provider-examplecloud/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// these will be set by the goreleaser configuration
// to appropriate values for the compiled binary.
var version string = "dev" // goreleaser can pass other information to the main package, such as the specific commit
// https://goreleaser.com/cookbooks/using-main.version/

func main() {
	ctx := context.Background()

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	upgradedSdkServer, err := tf5to6server.UpgradeServer(
        ctx,
        provider.Provider().GRPCProvider, // Example terraform-plugin-sdk provider
    )

    if err != nil {
        log.Fatal(err)
    }


	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		Address: "registry.terraform.io/Snowflake-Labs/snowflake",
		Debug:   debug,
	}
	tflog.Info(ctx, "Hello!")

	defer func() {
		tflog.Info(ctx, "Goodbye!")
	}()
	err := providerserver.Serve(ctx, provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
