package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/go-misc/ver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

const ProviderAddr = "registry.terraform.io/Snowflake-Labs/snowflake"

func main() {
	version := flag.Bool("version", false, "spit out version for resources here")
	debug := flag.Bool("debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if *version {
		verString, err := ver.VersionStr()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(verString)
		return
	}

	plugin.Serve(&plugin.ServeOpts{
		Debug:        *debug,
		ProviderAddr: ProviderAddr,
		ProviderFunc: provider.Provider,
	})
}
