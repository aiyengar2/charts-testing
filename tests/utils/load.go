package utils

import (
	"fmt"
	"sync"

	"emperror.dev/errors"
	"github.com/gobuffalo/packr"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/instantiatedcheck"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"gopkg.in/yaml.v2"
)

var (
	box = packr.NewBox("./yamls")

	loadOnce   sync.Once
	instChecks []instantiatedcheck.InstantiatedCheck
	loadErr    error
)

// Load gets all options from yamls/ and creates the kubelinter templates and checks accordingly
func Load() (instChecks []instantiatedcheck.InstantiatedCheck, loadErr error) {
	loadOnce.Do(func() {
		// Load tests, register templates, and collect checks
		var checks []check.Check
		for _, fileName := range box.List() {
			contents, err := box.Find(fileName)
			if err != nil {
				loadErr = errors.Wrapf(err, "loading template from %s", fileName)
				return
			}
			var t Test
			if err := yaml.Unmarshal(contents, &t); err != nil {
				loadErr = errors.Wrapf(err, "unmarshalling template options from %s", fileName)
				return
			}
			templates.Register(GetTemplateFromTest(t))
			checks = append(checks, GetCheckFromTest(t)...)
		}

		// Validate and instantiate checks
		for _, check := range checks {
			instantiatedChk, err := instantiatedcheck.ValidateAndInstantiate(&check)
			if err != nil {
				loadErr = errors.Wrapf(err, "invalid check %s", check.Name)
				return
			}
			if instantiatedChk == nil {
				loadErr = fmt.Errorf("instantiated check for %q is nil", check.Name)
				return
			}
			instChecks = append(instChecks, *instantiatedChk)
		}
	})
	return instChecks, loadErr
}
