package tests

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rancher/charts/testing/builder"
	"github.com/sirupsen/logrus"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"helm.sh/helm/v3/pkg/cli/values"
)

const (
	valuesFilesDir = "values"
)

var (
	chdirOnce sync.Once
	chdirErr  error

	testsDir string
	repoRoot string
)

func setup() error {
	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	if err := chdirToRepositoryRoot(); err != nil {
		return err
	}
	return nil
}

func chdirToRepositoryRoot() error {
	chdirOnce.Do(func() {
		var err error
		testsDir, err = os.Getwd()
		if err != nil {
			chdirErr = err
		}
		repoRoot = filepath.Dir(filepath.Dir(filepath.Dir(testsDir)))
		if err = os.Chdir(repoRoot); err != nil {
			chdirErr = fmt.Errorf("Could not change working directory from %s to repository root at %s: %s", testsDir, repoRoot, err)
		}
	})
	return chdirErr
}

type Parser interface {
	ParseTemplate(template string, glob string) error
	ParseTemplateWithOptions(template string, glob string, options builder.ParseOptions) error
}

func ensureHelmChartExists(chartPath string) error {
	chartPathInfo, err := os.Stat(chartPath)
	if err != nil {
		return fmt.Errorf("Encountered error while trying to describe %s: %s", chartPath, err)
	}
	// Check if chartPath points to a tgz file
	if !chartPathInfo.IsDir() {
		if filepath.Ext(chartPath) != ".tgz" {
			return fmt.Errorf("Not a valid Helm chart: %s does not point to a chart archive", chartPath)
		}
		return nil
	}
	// Check if the directory has a Chart.yaml
	files, err := ioutil.ReadDir(chartPath)
	if err != nil {
		return fmt.Errorf("Encountered error while trying to read files in %s: %s", chartPath, err)
	}
	foundChartYaml := false
	for _, file := range files {
		if file.Name() == "Chart.yaml" {
			foundChartYaml = true
			break
		}
	}
	if !foundChartYaml {
		return fmt.Errorf("Not a valid Helm chart: %s does not have a Chart.yaml", chartPath)
	}
	return nil
}

func parseValuesForChart(p Parser, chartPath string, strict bool) error {
	chartPath = filepath.Clean(chartPath)
	logrus.Infof("Parsing templates for %s from values directory", chartPath)
	// Validate that the chartPath points to a single Helm chart
	if err := ensureHelmChartExists(chartPath); err != nil {
		return err
	}
	valuesFilesDirpath := filepath.Join(testsDir, valuesFilesDir)
	valuesFiles, err := ioutil.ReadDir(valuesFilesDirpath)
	if err != nil {
		return fmt.Errorf("Could not read values.yaml files in dir %s: %s", valuesFilesDirpath, err)
	}
	if len(valuesFiles) == 0 {
		logrus.Warnf("Could not find any values.yaml files in dir %s: %s", valuesFilesDirpath, err)
		return nil
	}
	for _, valuesFile := range valuesFiles {
		valuesAbsPath := filepath.Join(valuesFilesDirpath, valuesFile.Name())
		valuesPath := filepath.Clean(strings.TrimPrefix(valuesAbsPath, repoRoot+"/"))
		err := p.ParseTemplateWithOptions(
			valuesFile.Name(),
			chartPath,
			builder.ParseOptions{
				Strict: strict,
				Options: lintcontext.Options{
					HelmValuesOptions: values.Options{
						ValueFiles: []string{valuesPath},
					},
				},
			},
		)
		if err != nil {
			return fmt.Errorf("Failed to parse %s: %s", valuesFile.Name(), err)
		}
	}
	return nil
}
