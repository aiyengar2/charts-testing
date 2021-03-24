package utils

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
)

type Test struct {
	// Name is the name of the test
	Name string `json:"name" yaml:"name"`
	// Info contains human-readable information about this test
	Info Info `json:"info" yaml:"info"`
	// Spec describes the input that is expected for this test
	Spec Spec `json:"spec" yaml:"spec"`
	// Cases is a list of test cases that implement the test Spec
	Cases []Case `json:"cases" yaml:"cases"`
}

type Info struct {
	// DisplayName is the human-readable name of the test
	DisplayName string `json:"display_name" yaml:"displayName"`
	// Description describes what this test checks for
	Description string `json:"description" yaml:"description"`
}

type Spec struct {
	// Kinds is a list of resource kinds that this test will be applied on
	Kinds []string `json:"kinds" yaml:"kinds"`
	// Flags are a list of specifications for flags that can be submitted by each test Case
	Flags []check.ParameterDesc `json:"flags" yaml:"flags"`
}

func (f FlagSpec) ToParameterDesc() (paramDesc check.ParameterDesc, err error) {
	var flagType check.ParameterType
	var isPointer bool
	switch f.Type {
	case "string":
		flagType = check.StringType
		isPointer = false
	case "integer":
		flagType = check.IntegerType
		isPointer = false
	case "boolean":
		flagType = check.BooleanType
		isPointer = false
	case "number":
		flagType = check.NumberType
		isPointer = false
	case "object":
		flagType = check.ObjectType
		isPointer = true
	case "array":
		flagType = check.ArrayType
		isPointer = true
	default:
		return paramDesc, fmt.Errorf("Invalid flag type provided: %s", f.Type)
	}
	var flagArrayType check.ParameterType
	switch f.ArrayElemType {
	case "string":
		flagArrayType = check.StringType
	case "integer":
		flagArrayType = check.IntegerType
	case "boolean":
		flagArrayType = check.BooleanType
	case "number":
		flagArrayType = check.NumberType
	case "object":
		flagArrayType = check.ObjectType
	case "array":
		flagArrayType = check.ArrayType
	default:
		return paramDesc, fmt.Errorf("Invalid flag type provided: %s", f.Type)
	}
	flagSubParams := make([]check.ParameterDesc, len(f.SubFlags))
	for i, subFlag := range f.SubFlags {
		flagSubParams[i], err = subFlag.ToParameterDesc()
		if err != nil {
			return paramDesc, err
		}
	}
	paramDesc = check.ParameterDesc{
		Name:               f.Name,
		Type:               flagType,
		Description:        f.Description,
		Examples:           f.Examples,
		Enum:               f.Enum,
		SubParameters:      flagSubParams,
		ArrayElemType:      flagArrayType,
		Required:           f.Required,
		NoRegex:            f.NoRegex,
		NotNegatable:       f.NotNegatable,
		XXXStructFieldName: strings.ToTitle(f.Name),
		XXXIsPointer:       isPointer,
	}
	return paramDesc, nil
}

type FlagSpec struct {
	// Name is the name of the flag
	Name string `json:"name" yaml:"name"`
	// Type is the type of input expected for this flag, e.g. "string", "integer", "boolean", "number", "object", "array"
	Type string `json:"type" yaml:"type"`
	// Description is a description of the flag
	Description string
	// Examples are examples that can be given for this specific flag
	Examples []string
	// Enum is set if the object is always going to be one of a specified set of values.
	// Only relevant if Type is "string"
	Enum []string
	// SubFlags are the child flags of the given flag
	// Only relevant if Type is "object".
	SubFlags []FlagSpec
	// ArrayElemType is only set when the object is of type array, and it describes the type
	// of the element of the array.
	ArrayElemType string
	// Required denotes whether the parameter is required.
	Required bool
	// NoRegex is set if the parameter does not support regexes.
	// Only relevant if Type is "string".
	NoRegex bool
	// NotNegatable is set if the parameter does not support negation via a leading !.
	// OnlyRelevant if Type is "string".
	NotNegatable bool
}

type Case struct {
	// Name is the list of this specific test case
	Name string `json:"name" yaml:"name"`
	// Description describes what this test case checks for
	Description string `json:"description" yaml:"description"`
	// Remediation describes how to fix an issue introduced by this test Case
	Remediation string `json:"remediation" yaml:"remediation"`
	// Kinds are the kinds of resources this test case is scoped to. Must be a subset of the kinds provided in the test Spec
	Kinds []string `json:"kinds" yaml:"kinds"`
	// Flags implement the FlagSpec defined in the test's Spec
	Flags map[string]interface{} `json:"flags" yaml:"flags"`
}

func GetTemplateFromTest(t Test) (template check.Template) {
	return check.Template{
		Key:                    fmt.Sprintf("%s-template", t.Name),
		HumanName:              t.Info.DisplayName,
		Description:            t.Info.Description,
		SupportedObjectKinds:   check.ObjectKindsDesc{ObjectKinds: t.Spec.Kinds},
		Parameters:             t.Spec.Flags,
		ParseAndValidateParams: GetParseAndValidateParamsFuncFromTest(t),
		Instantiate:            GetInstantiateFuncFromTest(t),
	}
}

func GetCheckFromTest(t Test) (checks []check.Check) {
	checks = make([]check.Check, len(t.Cases))
	for i, testCase := range t.Cases {
		checks[i] = check.Check{
			Name:        fmt.Sprintf("%s-%s", t.Name, testCase.Name),
			Description: testCase.Description,
			Remediation: testCase.Remediation,
			Template:    fmt.Sprintf("template-%s", t.Name),
			Scope:       &check.ObjectKindsDesc{ObjectKinds: testCase.Kinds},
			Params:      testCase.Flags,
		}
	}
	return checks
}

// ParseAndValidateParamsFunc tkaes in a list of params and returns an object that can be parsed by InstantiateFunc
type ParseAndValidateParamsFunc func(params map[string]interface{}) (interface{}, error)

func GetParseAndValidateParamsFuncFromTest(t Test) ParseAndValidateParamsFunc {
	return nil
}

// InstantiateFunc takes in the object from ParseAndValidateParamsFunc and returns a function that executes the check
type InstantiateFunc func(parsedParams interface{}) (check.Func, error)

func GetInstantiateFuncFromTest(t Test) InstantiateFunc {
	return nil
}
