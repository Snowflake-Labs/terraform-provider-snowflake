package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/chanzuckerberg/go-misc/ver"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	version := flag.Bool("version", false, "spit out version for resources here")
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
		ProviderFunc: func() terraform.ResourceProvider {
			return provider.Provider()
		},
	})

}
