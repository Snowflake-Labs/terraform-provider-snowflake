package archtests

import (
	"fmt"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

// get all files from package
// filter all files ending with _acceptance_test.go
// check all exported methods start with TestAcc_
// list all failing methods
func TestArchCheck_AcceptanceTests_Resources(t *testing.T) {
	path := "../resources/"
	packagesDict, err := parser.ParseDir(token.NewFileSet(), path, nil, 0)
	require.NoError(t, err)
	fmt.Printf("%v", packagesDict)
	fmt.Printf("%v", packagesDict["resources"])
	fmt.Printf("%v", packagesDict["resources_test"])
}
