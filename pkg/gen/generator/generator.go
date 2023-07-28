package generator

import (
	"fmt"
)

func GenerateAll(api any, sb ...fmt.Stringer) {
	for _, b := range sb {
		fmt.Printf("%s\n", b.String())
	}
}
