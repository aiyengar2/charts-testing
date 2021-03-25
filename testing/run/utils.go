package run

import (
	"fmt"
	"path/filepath"

	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

func GetLintCtxs(globs []string) (lintCtxs []lintcontext.LintContext, err error) {
	files, err := getFiles(globs)
	if err != nil {
		return nil, err
	}
	lintCtxs, err = lintcontext.CreateContexts(files...)
	if err != nil {
		return nil, err
	}
	// Ensure that at least one lint object is found
	var atLeastOneObjectFound bool
	for _, lintCtx := range lintCtxs {
		if len(lintCtx.Objects()) > 0 {
			atLeastOneObjectFound = true
			break
		}
	}
	if !atLeastOneObjectFound {
		return nil, fmt.Errorf("Did not find any objects to run tests on within %v", globs)
	}
	return lintCtxs, nil
}

func getFiles(globs []string) (files []string, err error) {
	for _, glob := range globs {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse glob %s: %s", glob, err)
		}
		files = append(files, matches...)
	}
	return files, nil
}
