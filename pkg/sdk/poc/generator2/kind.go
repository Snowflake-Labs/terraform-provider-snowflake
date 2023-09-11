package generator2

import (
	"fmt"
	"reflect"
)

type Kind interface {
	Kind() string
}

type kindFormString struct {
	kind string
}

func (k kindFormString) Kind() string {
	return k.kind
}

func KindOf(kind string) Kind {
	return &kindFormString{
		kind: kind,
	}
}

func KindOfPointer(kind string) Kind {
	return KindOf("*" + kind)
}

func KindOfSlice(kind string) Kind {
	return KindOf("[]" + kind)
}

func KindOfT[T any]() Kind {
	return KindOf(getTypeName[T]())
}

func getTypeName[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	return fmt.Sprint(t)
}
