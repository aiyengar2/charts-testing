package tests

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DoFunc func(objs []v1.Object) (pass bool)

var (
	objectInterface = reflect.TypeOf(new(v1.Object)).Elem()
)

// TODO(aiyengar2): Allow doFunc to be able to take in a single Resource instead of slice. Fail on runtime
// if more than one of that type is received.

// TODO(aiyengar2): Allow user to define whether a test is strict on declared resources (e.g. at least 1 of every field must be registered)

// TODO(aiyengar2): Allow a user to define whether a test ignores unknown resources that are collected

func WrapFunc(fromFunc interface{}) DoFunc {
	caller := reflect.ValueOf(fromFunc)
	funcType := reflect.TypeOf(fromFunc)

	// Ensure the input is a function
	if funcType.Kind() != reflect.Func {
		logrus.Fatalf("could not wrap doFunc: expected function, received %v", funcType.Kind())
	}
	// Function should only take in one argument
	if funcType.NumIn() != 1 {
		logrus.Fatalf("could not wrap doFunc: expected function that takes in exactly 1 argument, found %d", funcType.NumIn())
	}
	// Function argument should be a struct
	inputType := funcType.In(0)
	if inputType.Kind() != reflect.Struct {
		logrus.Fatalf("could not wrap doFunc: expected function that takes in a struct, found %d", inputType.Kind())
	}
	// Each field on the input struct should be a slice of unique type whose objects implement v1.Object
	k8sTypes := make(map[reflect.Type]string, inputType.NumField())
	supportedTypes := make([]reflect.Type, inputType.NumField())
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		fieldType := field.Type
		// Field must be a slice
		if fieldType.Kind() != reflect.Slice {
			logrus.Fatalf("could not wrap doFunc: field %s must be a slice", field.Name)
		}
		fieldElemType := fieldType.Elem()
		fieldElemPointerType := reflect.PtrTo(fieldElemType)
		// The contents of the slice must be a struct
		if fieldElemType.Kind() != reflect.Struct {
			logrus.Fatalf("could not wrap doFunc: field %s must be a slice of structs", field.Name)
		}
		// Field elem must implement v1.Object
		if !fieldElemPointerType.Implements(objectInterface) {
			logrus.Fatalf("could not wrap doFunc: field %s has elements of type %s that do not implement %s", field.Name, fieldElemPointerType, objectInterface)
		}
		// Field elem type must be unique
		fieldName, seen := k8sTypes[fieldElemPointerType]
		if seen {
			logrus.Fatalf("could not wrap doFunc: field %s and %s track the same k8s object %s", fieldName, field.Name, fieldElemPointerType)
		}
		k8sTypes[fieldElemPointerType] = field.Name
		supportedTypes[i] = fieldElemPointerType
	}

	// Function should only output 1 argument
	if funcType.NumOut() != 1 {
		logrus.Fatalf("could not wrap doFunc: expected function that outputs exactly 1 return value, found %d", funcType.NumOut())
	}
	// Function should output a bool
	if funcType.Out(0).Kind() != reflect.Bool {
		logrus.Fatalf("could not wrap doFunc: expected function that outputs a bool, found %d", funcType.Out(1).Kind())
	}

	return func(objs []v1.Object) (pass bool) {
		in := reflect.New(inputType)
		for _, obj := range objs {
			objType := reflect.TypeOf(obj)
			fieldName, ok := k8sTypes[objType]
			if !ok {
				logrus.WithFields(logrus.Fields{
					"resource":       fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName()),
					"supportedTypes": supportedTypes,
				}).Fatalf("Could not unmarshall %s into %s", reflect.TypeOf(obj), inputType)
			}
			field := in.Elem().FieldByName(fieldName)
			field.Set(reflect.Append(field, reflect.ValueOf(obj).Elem()))
		}
		args := []reflect.Value{in.Elem()}
		return caller.Call(args)[0].Bool()
	}
}
