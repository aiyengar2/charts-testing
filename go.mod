module github.com/rancher/charts

go 1.16

replace (
	golang.stackrox.io/kube-linter v0.0.0-20210316191552-241e80db4436 => github.com/aiyengar2/kube-linter v0.0.0-20210327012337-244a1733c800
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)

require (
	github.com/hashicorp/go-multierror v1.1.0
	github.com/rancher/cis-operator v1.0.3
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.7.0 // indirect
	golang.stackrox.io/kube-linter v0.0.0-20210316191552-241e80db4436
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kubectl v0.18.8
)
