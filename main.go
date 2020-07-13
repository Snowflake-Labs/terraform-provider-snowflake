package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/chanzuckerberg/go-misc/ver"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/olekukonko/tablewriter"
)

func main() {
	doc := flag.Bool("doc", false, "spit out docs for resources here")
	version := flag.Bool("version", false, "spit out version for resources here")
	flag.Parse()

	if *doc {
		generateDocs()
		return
	}

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

func generateDocs() {
	// schema := provider.Provider().Schema
	resources := provider.Provider().ResourcesMap

	names := make([]string, 0)
	for k := range resources {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, name := range names {
		resource := resources[name]
		fmt.Printf("\n### %s\n\n", name)
		if strings.HasSuffix(name, "_grant") {
			grant_resource_name := strings.Replace(name, "_grant", "", -1)
			granted_to_name := strings.Replace(grant_resource_name, "snowflake_", "", -1)
			fmt.Printf(
				`**Note**: The %s resource creates exclusive attachments of grants.
Across the entire Snowflake account, all of the %ss to which a single grant is attached must be declared
by a single %s resource. This means that even any %s that have the attached
grant via any other mechanism (including other Terraform resources) will have that attached grant revoked by this resource.
These resources do not enforce exclusive attachment of a grant, it is the user's responsibility to enforce this.
`, name, granted_to_name, name, grant_resource_name)
			fmt.Println("")
		}
		fmt.Printf("#### properties\n\n")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(false)
		table.SetHeader([]string{"name", "type", "description", "optional", " required", "computed", "default"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		properties := make([]string, 0)
		for k := range resource.Schema {
			properties = append(properties, k)
		}
		sort.Strings(properties)
		for _, property := range properties {
			s := resource.Schema[property]
			table.Append([]string{property, typeString(s.Type), s.Description, boolString(s.Optional), boolString(s.Required), boolString(s.Computed), interfaceString(s.Default)})
		}
		table.Render()
	}

}

func typeString(t schema.ValueType) string {
	switch t {
	case schema.TypeBool:
		return "bool"
	case schema.TypeInt:
		return "int"
	case schema.TypeFloat:
		return "float"
	case schema.TypeString:
		return "string"
	case schema.TypeList:
		return "list"
	case schema.TypeMap:
		return "map"
	case schema.TypeSet:
		return "set"
	}
	return "?"
}

func boolString(t bool) string { return fmt.Sprintf("%t", t) }

func interfaceString(t interface{}) string { return fmt.Sprintf("%#v", t) }
