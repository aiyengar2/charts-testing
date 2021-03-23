package templates

import (
	"sync"

	"emperror.dev/errors"
	"github.com/gobuffalo/packr"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"gopkg.in/yaml.v2"
)

var (
	box = packr.NewBox("./yamls")

	loadOnce sync.Once
	loadErr  error
)

// RegisterTemplates registers all the templates stored in yamls/ to the templates registry
func RegisterTemplates() (loadErr error) {
	loadOnce.Do(func() {
		for _, fileName := range box.List() {
			contents, err := box.Find(fileName)
			if err != nil {
				loadErr = errors.Wrapf(err, "loading template from %s", fileName)
				return
			}
			var opts Options
			if err := yaml.Unmarshal(contents, &opts); err != nil {
				loadErr = errors.Wrapf(err, "unmarshalling template options from %s", fileName)
				return
			}
			templates.Register(GetTemplateFromOptions(opts))
		}
	})
	return loadErr
}
