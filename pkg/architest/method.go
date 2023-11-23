package architest

type Method struct {
	name string
	file *File
}

func (method *Method) Name() string {
	return method.name
}

func (method *Method) FileName() string {
	return method.file.FileName()
}

func NewMethod(name string, file *File) *Method {
	return &Method{
		name: name,
		file: file,
	}
}
