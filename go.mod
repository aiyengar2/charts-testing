module github.com/rancher/charts

go 1.16

replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)

require (
	github.com/fatih/color v1.9.0
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/swag v0.19.10 // indirect
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	golang.stackrox.io/kube-linter v0.0.0-20210316191552-241e80db4436
	k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver v0.19.2 // indirect
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v12.0.0+incompatible // indirect
)
