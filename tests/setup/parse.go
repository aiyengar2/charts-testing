package setup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rancher/charts/testing/builder"
	"github.com/sirupsen/logrus"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"helm.sh/helm/v3/pkg/cli/values"
)

type Parser interface {
	ParseTemplate(template string, glob string) error
	ParseTemplateWithOptions(template string, glob string, options builder.ParseOptions) error
}

func ParseValuesForChart(p Parser, chartPath string, strict bool) error {
	chartPath = filepath.Clean(chartPath)
	logrus.Infof("Parsing templates for %s", chartPath)
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
