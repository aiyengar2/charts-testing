package builder

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

const (
	notK8sObjErrPrefix = "failed to decode: Object 'Kind' is missing in"
)

func parseTemplate(glob string, options ParseOptions) ([]lintcontext.Object, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse glob %s: %s", glob, err)
	}
	lintCtxs, err := lintcontext.CreateContextsWithOptions(options.Options, files...)
	if err != nil {
		return nil, err
	}
	var atLeastOneObjectFound bool
	objs := []lintcontext.Object{}
	for _, lintCtx := range lintCtxs {
		// Add valid objects to template
		lintObjs := lintCtx.Objects()
		if len(lintCtx.Objects()) > 0 {
			atLeastOneObjectFound = true
			objs = append(objs, lintObjs...)
		}
		invalidObjs := lintCtx.InvalidObjects()
		if len(invalidObjs) > 0 {
			// Handle invalid objects
			if options.Strict {
				var err *multierror.Error
				for _, obj := range invalidObjs {
					loadErr := obj.LoadErr
					if strings.Contains(loadErr.Error(), notK8sObjErrPrefix) {
						logrus.Warn(loadErr)
						continue
					}
					err = multierror.Append(err, loadErr)
				}
				return nil, fmt.Errorf("Unable to parse template with invalid objects: %s", err)
			} else {
				logrus.Warnf("Encountered errors loading k8s objects from template %s", glob)
				for _, obj := range invalidObjs {
					logrus.Warn(obj.LoadErr)
				}
			}
		}
	}
	if !atLeastOneObjectFound {
		return nil, fmt.Errorf("%v does not contain a template with k8s objects", glob)
	}
	return objs, nil
}
