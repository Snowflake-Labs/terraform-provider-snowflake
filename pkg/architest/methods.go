package architest

type MethodReceiver = func(method *Method)
type Methods []Method

func (methods Methods) All(receiver MethodReceiver) {
	for _, method := range methods {
		receiver(&method)
	}
}
