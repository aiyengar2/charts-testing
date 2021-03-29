package tests

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	invalidFuncSignatureWarning = "Function signature must match pattern func(ctx context.Context, objStruct MyStruct)"
	invalidStructTypeWarning    = "Each field of a provided struct must correspond to an object or slice of objects that implement v1.Object. " +
		"It is also expected that no two fields the struct correspond to the same underlying type."
)

var (
	objectInterface = reflect.TypeOf(new(v1.Object)).Elem()
	ctxInterface    = reflect.TypeOf(new(context.Context)).Elem()
)

type DoFunc func(ctx context.Context, objs []v1.Object) (pass bool)

func WrapFunc(fromFunc interface{}) DoFunc {
	caller := reflect.ValueOf(fromFunc)
	funcName := runtime.FuncForPC(caller.Pointer()).Name()
	funcType := reflect.TypeOf(fromFunc)

	err := validateFunctionSignature(funcType)
	if err != nil {
		logrus.Errorf(invalidFuncSignatureWarning)
		logrus.WithField("doFunc", funcName).Fatalf("could not wrap doFunc: %s", err)
	}

	inputType := funcType.In(1)
	supportedTypes, err := parseInputStruct(inputType)
	if err != nil {
		logrus.Errorf(invalidStructTypeWarning)
		logrus.WithField("doFunc", funcName).Fatalf("could not wrap doFunc: %s", err)
	}

	return func(ctx context.Context, objs []v1.Object) (pass bool) {
		in := reflect.New(inputType)

		singletonFieldToObjs := map[reflect.Value][]v1.Object{}
		for _, obj := range objs {
			objType := reflect.TypeOf(obj)
			fieldName, err := supportedTypes.getField(objType)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"resource":       fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName()),
					"supportedTypes": supportedTypes,
				}).Fatalf("Could not unmarshall %s into %s: %s", reflect.TypeOf(obj), inputType, err)
			}
			field := in.Elem().FieldByName(fieldName)
			if field.Kind() != reflect.Slice {
				// Singleton fields will be added in a separate loop to print violaters
				singletonFieldToObjs[field] = append(singletonFieldToObjs[field], obj)
				continue
			}
			objValue := reflect.ValueOf(obj)
			if field.Type().Elem().Kind() != reflect.Interface {
				objValue = objValue.Elem()
			}
			field.Set(reflect.Append(field, objValue))
		}

		// Ensure that you only find one resource if the struct is expecting a single resource
		foundMultipleSingletons := false
		for field, objs := range singletonFieldToObjs {
			if len(objs) > 0 {
				resources := make([]string, len(objs))
				for i, obj := range objs {
					resources[i] = fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
				}
				logrus.
					WithField("resources", resources).
					Errorf("Expected 1 resource of type %s, found %d", field.Type(), len(objs))
				foundMultipleSingletons = true
				continue
			}
			objValue := reflect.ValueOf(objs[0])
			if field.Type().Elem().Kind() != reflect.Interface {
				objValue = objValue.Elem()
			}
			field.Set(objValue)
		}

		if foundMultipleSingletons {
			logrus.
				WithField("doFunc", funcName).
				Fatalf("Failed to unmarshall objects into %s", inputType)
		}

		// Call the function with the provided context
		args := []reflect.Value{reflect.ValueOf(ctx), in.Elem()}
		return caller.Call(args)[0].Bool()
	}
}

func validateFunctionSignature(funcType reflect.Type) error {
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("expected function, received %v", funcType.Kind())
	}
	// Function should take in two arguments
	if funcType.NumIn() != 2 {
		return fmt.Errorf("expected function that takes in exactly 2 arguments, found %d", funcType.NumIn())
	}
	// First function argument should be a context
	ctxType := funcType.In(0)
	if !ctxType.Implements(ctxInterface) {
		return fmt.Errorf("expected first argument to implement %s, found %s", ctxInterface, ctxType.Kind())
	}
	// Function argument should be a struct
	inputType := funcType.In(1)
	if inputType.Kind() != reflect.Struct {
		return fmt.Errorf("expected second argument of the function to be a struct, found %d", inputType.Kind())
	}
	// Function should only output 1 argument
	if funcType.NumOut() != 1 {
		return fmt.Errorf("expected function that outputs exactly 1 return value, found %d", funcType.NumOut())
	}
	// Function should output a bool
	if funcType.Out(0).Kind() != reflect.Bool {
		return fmt.Errorf("expected function that outputs a bool, found %d", funcType.Out(1).Kind())
	}
	return nil
}

func parseInputStruct(inputType reflect.Type) (fieldTypeTracker, error) {
	// Each field on the input struct should either implement v1.Object or be a slice of objects that implement v1.Object
	supportedResources := fieldTypeTracker{
		types:      map[reflect.Type]string{},
		interfaces: map[reflect.Type]string{},
	}
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		fieldType := field.Type
		// Check if field is a struct or a slice
		var fieldPointerType reflect.Type
		var isInterface bool
		switch fieldType.Kind() {
		case reflect.Slice:
			fieldElemType := fieldType.Elem()
			switch fieldElemType.Kind() {
			case reflect.Struct:
				fieldPointerType = reflect.PtrTo(fieldElemType)
			case reflect.Interface:
				fieldPointerType = fieldElemType
				isInterface = true
			default:
				return supportedResources, fmt.Errorf("field %s must be a slice of structs", field.Name)
			}
		case reflect.Struct:
			fieldPointerType = reflect.PtrTo(fieldType)
		case reflect.Interface:
			fieldPointerType = fieldType
			isInterface = true
		default:
			return supportedResources, fmt.Errorf("field %s must be a struct or slice", field.Name)
		}
		if isInterface {
			err := supportedResources.addInterface(fieldPointerType, field.Name)
			if err != nil {
				return supportedResources, err
			}
		} else {
			err := supportedResources.addType(fieldPointerType, field.Name)
			if err != nil {
				return supportedResources, err
			}
		}
		// Field elem must implement v1.Object
		if !fieldPointerType.Implements(objectInterface) {
			return supportedResources, fmt.Errorf("field %s contains object(s) of type %s that do not implement %s", field.Name, fieldPointerType, objectInterface)
		}
	}
	return supportedResources, nil
}
