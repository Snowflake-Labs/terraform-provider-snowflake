//go:build exclude

package main

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

type (
	Set    struct{}
	TagSet struct {
		Tags []string
	}
	NoAction struct{}
)

type OneOfAlterSet interface {
	~*Set | ~*TagSet | ~*NoAction
}

type Alter[T OneOfAlterSet] struct {
	AlterSet T
}

func Call() {
	alter := Alter[*TagSet]{
		AlterSet: &TagSet{
			Tags: []string{"one", "two"},
		},
	}
	_ = alter
}
