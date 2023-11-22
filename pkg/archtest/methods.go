package archtest

type MethodReceiver = func(method *Method)

func IterateMethods(methods []Method, receiver MethodReceiver) {
	for _, method := range methods {
		receiver(&method)
	}
}
