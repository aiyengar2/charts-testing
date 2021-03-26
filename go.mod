module github.com/rancher/charts

go 1.16

replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)

require (
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.46.0
	github.com/rancher/cis-operator v1.0.3
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.7.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	helm.sh/helm/v3 v3.3.4
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/cli-runtime v0.20.4 // indirect
	k8s.io/kubectl v0.18.8
)
