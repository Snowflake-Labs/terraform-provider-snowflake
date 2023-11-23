package architest

type Method struct {
	name string
	file *File
}

func NewMethod(name string, file *File) *Method {
	return &Method{
		name: name,
		file: file,
	}
}
