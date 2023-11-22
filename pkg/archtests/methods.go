package archtests

type MethodReceiver = func(method *Method)

func iterateMethods(methods []Method, receiver MethodReceiver) {
	for _, method := range methods {
		receiver(&method)
	}
}
