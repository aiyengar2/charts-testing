package kubelinter

import (
	"fmt"
	"os"

	"emperror.dev/errors"
	"golang.org/x/crypto/ssh/terminal"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/ignore"
	"golang.stackrox.io/kube-linter/pkg/instantiatedcheck"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"

	"github.com/rancher/charts/tests/kubelinter/checks"
)

func RunTests(globs []string) (err error) {
	instChecks, err := checks.GetInstantiatedChecks()
	if err != nil {
		return err
	}
	lintCtxs, err := getLintCtxs(globs)
	if err != nil {
		return err
	}
	result, err := RunChecks(lintCtxs, instChecks)
	if err != nil {
		return err
	}
	if len(result.Reports) == 0 {
		fmt.Fprintln(os.Stderr, "All tests passed")
		return nil
	}
	stderrIsTerminal := terminal.IsTerminal(int(os.Stderr.Fd()))
	for _, report := range result.Reports {
		if stderrIsTerminal {
			report.FormatToTerminal(os.Stderr)
		} else {
			report.FormatPlain(os.Stderr)
		}
	}
	return errors.Errorf("Some tests failed", len(result.Reports))
}

func RunChecks(lintCtxs []lintcontext.LintContext, instChecks []instantiatedcheck.InstantiatedCheck) (result run.Result, err error) {
	for _, lintCtx := range lintCtxs {
		for _, obj := range lintCtx.Objects() {
			for _, check := range instChecks {
				if !check.Matcher.Matches(obj.K8sObject.GetObjectKind().GroupVersionKind()) {
					continue
				}
				if ignore.ObjectForCheck(obj.K8sObject.GetAnnotations(), check.Spec.Name) {
					continue
				}
				diagnostics := check.Func(lintCtx, obj)
				for _, d := range diagnostics {
					result.Reports = append(result.Reports, diagnostic.WithContext{
						Diagnostic:  d,
						Check:       check.Spec.Name,
						Remediation: check.Spec.Remediation,
						Object:      obj,
					})
				}
			}
		}
	}
	return result, err
}
