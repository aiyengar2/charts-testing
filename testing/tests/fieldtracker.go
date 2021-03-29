package tests

import (
	"fmt"
	"reflect"
)

type fieldTypeTracker struct {
	types      map[reflect.Type]string
	interfaces map[reflect.Type]string
}

func (r *fieldTypeTracker) addType(typeToAdd reflect.Type, fieldName string) error {
	if currFieldName, exists := r.types[typeToAdd]; exists {
		return fmt.Errorf("field %s and %s track the same object type %s", currFieldName, fieldName, typeToAdd)
	}
	r.types[typeToAdd] = fieldName
	return nil
}

func (r *fieldTypeTracker) addInterface(interfaceToAdd reflect.Type, fieldName string) error {
	if currFieldName, exists := r.types[interfaceToAdd]; exists {
		return fmt.Errorf("field %s and %s track the same object interface %s", currFieldName, fieldName, interfaceToAdd)
	}
	r.interfaces[interfaceToAdd] = fieldName
	return nil
}

func (r *fieldTypeTracker) getField(fieldForType reflect.Type) (fieldName string, err error) {
	for supportedType, field := range r.types {
		if fieldForType == supportedType {
			return field, nil
		}
	}
	implementsFields := []string{}
	for supportedInterface, field := range r.interfaces {
		if fieldForType.Implements(supportedInterface) {
			implementsFields = append(implementsFields, field)
		}
	}
	if len(implementsFields) == 0 {
		return "", fmt.Errorf("no existing fields support %s", fieldForType)
	}
	if len(implementsFields) > 1 {
		return "", fmt.Errorf("placement of %s is ambiguous, can be marshalled into multiple interface fields: %s", fieldForType, implementsFields)
	}
	return implementsFields[0], nil
}

func (r fieldTypeTracker) String() string {
	supportedTypes := make([]reflect.Type, len(r.types))
	i := 0
	for supportedType := range r.types {
		supportedTypes[i] = supportedType
		i++
	}
	supportedInterfaces := make([]reflect.Type, len(r.interfaces))
	i = 0
	for supportedInterface := range r.interfaces {
		supportedInterfaces[i] = supportedInterface
		i++
	}
	return fmt.Sprintf("{types: %s, interfaces: %s}", supportedTypes, supportedInterfaces)
}
