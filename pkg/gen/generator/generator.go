package generator

import (
	"fmt"
)

func GenerateAll(sb ...fmt.Stringer) {
	for _, b := range sb {
		fmt.Printf("%s\n", b.String())
	}
}
