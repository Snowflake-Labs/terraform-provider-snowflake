package main

import (
	"flag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

const ProviderAddr = "registry.terraform.io/Snowflake-Labs/snowflake"



func main() {
	debug := flag.Bool("debug", false, "set to true to run the provider with support for debuggers like delve")
	
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug:        *debug,
		ProviderAddr: ProviderAddr,
		ProviderFunc: provider.Provider,
	})
}
