package checks

import (
	"fmt"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/instantiatedcheck"

	"github.com/rancher/charts/tests/kubelinter/templates"
)

var (
	box = packr.NewBox("./yamls")

	loadOnce sync.Once
	list     []instantiatedcheck.InstantiatedCheck
	loadErr  error
)

// GetInstantiatedChecks gets all the checks from checks/, validates them, and instantiates them
func GetInstantiatedChecks() (list []instantiatedcheck.InstantiatedCheck, loadErr error) {
	loadOnce.Do(func() {
		if err := templates.RegisterTemplates(); err != nil {
			loadErr = err
			return
		}
		for _, fileName := range box.List() {
			contents, err := box.Find(fileName)
			if err != nil {
				loadErr = errors.Wrapf(err, "loading check from %s", fileName)
				return
			}
			var chk check.Check
			if err := yaml.Unmarshal(contents, &chk); err != nil {
				loadErr = errors.Wrapf(err, "unmarshalling check from %s", fileName)
				return
			}
			instantiatedChk, err := instantiatedcheck.ValidateAndInstantiate(&chk)
			if err != nil {
				loadErr = errors.Wrapf(err, "invalid check %s", chk.Name)
				return
			}
			if instantiatedChk == nil {
				loadErr = fmt.Errorf("instantiated check for %q is nil", chk.Name)
				return
			}

			list = append(list, *instantiatedChk)
		}
	})
	return list, loadErr
}
