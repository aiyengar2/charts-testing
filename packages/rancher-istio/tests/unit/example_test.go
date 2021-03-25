package integration

import (
	"bytes"
	"github.com/ghodss/yaml"
	"io"
	v1 "k8s.io/api/batch/v1"
	v12 "k8s.io/api/core/v1"
	"log"
	"testing"

	goyaml "github.com/go-yaml/yaml"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	cisv1 "github.com/rancher/cis-operator/pkg/apis/cis.cattle.io/v1"
)

type Kind struct {
	Kind string `yaml:kind,omitempty`
}

type TemplateData struct {
	Name string
	Kind string
	Data interface{}
}


func TestImageValues(t *testing.T) {
	helmChartPath := "packages/rancher-cis-benchmark/charts"
	releaseName := "rancher-cis-benchmark"
	namespaceName := "cis-operator-system"

	options := &helm.Options{
		SetValues:      map[string]string{"namespace": namespaceName},
		KubectlOptions: k8s.NewKubectlOptions("", "", namespaceName),
	}
	//renders helm templates
	output := helm.RenderTemplate(t, options, helmChartPath, releaseName, []string{})

	//decodes templates into byte array
	yamlResult, err := SplitYAML([]byte(output))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var templates []TemplateData
	for _, template := range yamlResult {
		var something Kind
		err := yaml.Unmarshal(template, &something)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		switch something.Kind {
		case "Job":
			var job v1.Job
			helm.UnmarshalK8SYaml(t, string(template), &job)
			td := TemplateData{
				Name: job.Name,
				Kind: job.Kind,
				Data: job,
			}
			templates = append(templates, td)
		case "ServiceAccount":
			var serviceAccount v12.ServiceAccount
			helm.UnmarshalK8SYaml(t, string(template), &serviceAccount)
			td := TemplateData{
				Name: serviceAccount.Name,
				Kind: serviceAccount.Kind,
				Data: serviceAccount,
			}
			templates = append(templates, td)
		case "ClusterScanProfile":
			var csProfile cisv1.ClusterScanProfile
			helm.UnmarshalK8SYaml(t, string(template), &csProfile)
			td := TemplateData{
				Name: csProfile.Name,
				Kind: csProfile.Kind,
				Data: csProfile,
			}
			templates = append(templates, td)
		default:
			continue
		}
	}
}


func SplitYAML(resources []byte) ([][]byte, error) {

	dec := goyaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := goyaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}