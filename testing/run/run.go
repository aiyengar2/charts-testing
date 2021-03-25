package run

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rancher/charts/testing/tests"
	"github.com/sirupsen/logrus"
)

var (
	bold = color.New(color.Bold)
)

func RunTests(templateGlobs []string) error {
	lintCtxs, err := GetLintCtxs(templateGlobs)
	if err != nil {
		return err
	}
	testsToRun := tests.List()
	numTests := len(testsToRun)
	if numTests == 0 {
		logrus.Errorf("Could not find any tests to run")
		return nil
	}
	passedAllTests := true
	for _, lintCtx := range lintCtxs {
		lintObjs := lintCtx.Objects()
		invalidLinObjs := lintCtx.InvalidObjects()
		if len(invalidLinObjs) > 0 {
			for _, invalidObj := range invalidLinObjs {
				logrus.Errorf("%v", invalidObj.LoadErr)
			}
			logrus.Fatalf("Failed to parse k8s manifest")
		}
		for _, test := range testsToRun {
			objs, err := test.FilterObjects(lintObjs)
			if err != nil {
				err = fmt.Errorf("Unable to get resources for test %s: %s", test.ID, err)
				return err
			}
			logrus.Infof("%s Running test...", bold.Sprintf("(test %s)", test.ID))
			if test.Run(objs) {
				logrus.Infof("%s %s", bold.Sprintf("(test %s)", test.ID), color.GreenString("pass"))
			} else {
				passedAllTests = false
				logrus.Errorf("%s %s", bold.Sprintf("(test %s)", test.ID), color.RedString("fail"))
			}
		}
	}
	if passedAllTests {
		logrus.Infof("All tests passed.")
	} else {
		logrus.Errorf("Failed some tests.")
	}
	return nil
}
