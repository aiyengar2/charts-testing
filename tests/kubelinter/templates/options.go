package templates

import (
	"golang.stackrox.io/kube-linter/pkg/check"
)

type Options struct {
	// Name is the value used to refer to this template in a check
	Name string `json:"name" yaml:"name"`
	// DisplayName is the human-readable name for this name that shows up on test runs
	DisplayName string `json:"displayName" yaml:"displayName"`
	// Description describes what this template represents
	Description string `json:"description" yaml:"description"`
	// Scope is the objectKinds that this template will be applied on
	Scope check.ObjectKindsDesc `json:"scope" yaml:"scope"`
	// Parameters has a list of specifications for parameters that this template takes in
	ParameterSpecs []check.ParameterDesc `json:"paramSpecs" yaml:"paramSpecs"`
}

func GetTemplateFromOptions(opts Options) check.Template {
	return check.Template{
		Key:                    opts.Name,
		HumanName:              opts.DisplayName,
		Description:            opts.Description,
		SupportedObjectKinds:   opts.Scope,
		Parameters:             opts.ParameterSpecs,
		ParseAndValidateParams: GetParseAndValidateParamsFuncFromOptions(opts),
		Instantiate:            GetInstantiateFuncFromOptions(opts),
	}
}

// ParseAndValidateParamsFunc tkaes in a list of params and returns an object that can be parsed by InstantiateFunc
type ParseAndValidateParamsFunc func(params map[string]interface{}) (interface{}, error)

func GetParseAndValidateParamsFuncFromOptions(opts Options) ParseAndValidateParamsFunc {
	return nil
}

// InstantiateFunc takes in the object from ParseAndValidateParamsFunc and returns a function that executes the check
type InstantiateFunc func(parsedParams interface{}) (check.Func, error)

func GetInstantiateFuncFromOptions(opts Options) InstantiateFunc {
	return nil
}
