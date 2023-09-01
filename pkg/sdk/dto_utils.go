package sdk

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"unsafe"
)

func copyByFieldNames[Opt, Req any](opt *Opt, req *Req) {
	optType := reflect.TypeOf(opt)
	optValue := reflect.ValueOf(opt)
	reqType := reflect.TypeOf(req)
	reqValue := reflect.ValueOf(req)
	copyByFieldNamesImpl(optType, optValue, reqType, reqValue)
}

// TODO Rename to Opt and Req
func copyByFieldNamesImpl(outType reflect.Type, outValue reflect.Value, inType reflect.Type, inValue reflect.Value) {
	outElem := outValue.Elem()
	inElem := inValue.Elem()

	if outElem.Kind() != reflect.Struct || inElem.Kind() != reflect.Struct {
		panic(fmt.Sprintf("One of parameters [out (kind: %v), in (kind: %v)] is not a struct type\n", outElem.Kind(), inElem.Kind()))
	}

	for i := 0; i < outElem.NumField(); i++ {
		outFieldMeta := outElem.Type().Field(i)
		outFieldValue := outElem.Field(i)
		for j := 0; j < inElem.NumField(); j++ {
			inFieldMeta := inElem.Type().Field(j)
			inFieldValue := inElem.Field(j)
			if strings.EqualFold(outFieldMeta.Name, inFieldMeta.Name) {
				// this is a trick that lets us access inValue's unexported fields, otherwise we would get an error
				inValuePtr := reflect.NewAt(inFieldMeta.Type, unsafe.Pointer(inFieldValue.UnsafeAddr())).Elem()
				if inFieldMeta.Type != outFieldMeta.Type {
					// by DTO conventions we have right now, this mean we've encountered nested request e.g.
					// type Options {
					// 		Field1 int
					// 		NestedField NestedField - we assume NestedField and NestedRequest have the same fields
					// 		Field2 string
					// }
					// type Request struct {
					// 		field1 int
					// 		nestedField NestedRequest - this will map to Options.NestedField
					// 		field2 string
					// }
					log.Printf("Setting nested request %s.%s from %s.%s\n", outType.String(), outFieldMeta.Name, inType.String(), inFieldMeta.Name)
					// create instance of option's field type
					outValueInstance := reflect.New(outFieldMeta.Type.Elem())
					// recursive call to copy all the request's fields into instance of option's field type
					copyByFieldNamesImpl(outFieldMeta.Type, outValueInstance, inFieldMeta.Type, inValuePtr)
					// set option's field with instance of field type with copied fields from request's field
					outFieldValue.Set(outValueInstance)
				} else {
					log.Printf("Setting %s.%s from %s.%s\n", outType.String(), outFieldMeta.Name, inType.String(), inFieldMeta.Name)
					// TODO: Copy (right now it might be referencing value from opts - not sure - have to check)
					// set option's field with request's field - they have same type
					reflect.NewAt(outFieldMeta.Type, unsafe.Pointer(outFieldValue.UnsafeAddr())).Elem().Set(inValuePtr)
				}
				break
			}
		}
	}
}

func printStruct(s any) {
	elem := reflect.ValueOf(s).Elem()

	if elem.Kind() != reflect.Struct {
		panic("Not struct")
	}

	fmt.Println(elem.Type().Name() + " {")
	printStructImpl(elem, 2)
	fmt.Println("}")
}

func printStructImpl(value reflect.Value, indent int) {
	for i := 0; i < value.NumField(); i++ {
		fieldMeta := value.Type().Field(i)
		fieldValue := value.Field(i)

		if fieldValue.Kind() == reflect.Slice {
			fmt.Print(strings.Repeat(" ", indent) + fieldMeta.Name + ": [")
			for j := 0; j < fieldValue.Len(); j++ {
				s, _ := json.Marshal(fieldValue.Index(j).Interface())
				fmt.Printf("%s, ", string(s))
			}
			fmt.Println("]")
		} else if fieldValue.Kind() == reflect.Struct {
			fmt.Println(strings.Repeat(" ", indent) + fieldMeta.Name + ": {")
			printStructImpl(fieldValue, indent+2)
			fmt.Println(strings.Repeat(" ", indent) + "}")
		} else if fieldValue.Kind() == reflect.Pointer && !fieldValue.IsNil() {
			fmt.Println(strings.Repeat(" ", indent) + fieldValue.Elem().Type().Name() + ": {")
			printStructImpl(fieldValue.Elem(), indent+2)
			fmt.Println(strings.Repeat(" ", indent) + "}")
		} else {
			val := reflect.NewAt(fieldMeta.Type, unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
			fmt.Printf("%s%s: %v\n", strings.Repeat(" ", indent), fieldMeta.Name, val.Interface())
		}
	}
}
