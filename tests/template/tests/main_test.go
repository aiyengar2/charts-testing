package tests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/rancher/charts/testing/builder"
	"github.com/sirupsen/logrus"
	"k8s.io/kubectl/pkg/scheme"
)

var (
	suite = builder.NewTestSuite(packageName)

	strict         = flag.Bool("strict", true, "fail if any object cannot be decoded")
	chartDirectory = flag.String("chart", fmt.Sprintf("packages/%s/charts", packageName), "dir pointing to Helm chart")
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	// If you are using CRDs, you will need to call the corresponding AddToScheme calls here
	s := scheme.Scheme
	suite.SetCustomScheme(s)

	if err := parseValuesForChart(suite, *chartDirectory, *strict); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}
