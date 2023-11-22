package archtests

type Method struct {
	name string
}

func NewMethod(name string) *Method {
	return &Method{
		name: name,
	}
}
