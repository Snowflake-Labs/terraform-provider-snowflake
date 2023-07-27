//go:build exclude

package main

/*
type (
	OptionOne struct {
		b int
	}
	OptionTwo   struct{}
	OptionThree struct{}
)

type OneOf interface {
	~*OptionOne | ~*OptionTwo | ~*OptionThree
}

func OneOfFunc[T OneOf](one T) {
}

func OneOfFunc2[T constraints.Ordered](one T) {
}

func Call() {
	// Won't compile -> OneOfFunc[*OptionOne](&OptionTwo{})
	OneOfFunc[*OptionOne](&OptionOne{})
}
*/
