package tests

import (
	"os"
	"testing"

	"github.com/rancher/charts/testing/builder"
	cisv1 "github.com/rancher/cis-operator/pkg/apis/cis.cattle.io/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/kubectl/pkg/scheme"
)

const (
	testPrefix = "example"
)

var (
	suite = builder.NewTestSuite(testPrefix)
)

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// If you are using CRDs, you will need to call the corresponding AddToScheme calls here
	s := scheme.Scheme
	cisv1.AddToScheme(s)
	suite.SetCustomScheme(s)

	// Add templates that you plan to test here that may have errors on parsing
	templates := []string{"../../../packages/rancher-pushprox/charts"}
	for _, template := range templates {
		if err := suite.ParseTemplate(template); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}
	// Add templates that must not have errors while parsing
	strictTemplates := []string{"../../../packages/rancher-cis-benchmark/charts"}
	for _, template := range strictTemplates {
		if err := suite.ParseTemplateStrict(template); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}

	code := m.Run()
	os.Exit(code)
}
