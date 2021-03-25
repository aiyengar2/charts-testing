module github.com/rancher/charts

go 1.16

replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)

require (
	github.com/aws/aws-sdk-go v1.36.7 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gruntwork-io/terratest v0.32.16
	github.com/rancher/cis-operator v1.0.3
	github.com/rancher/wrangler v0.7.3-0.20210319211136-3eba78f45e7d // indirect
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.19.3
// indirect
)
