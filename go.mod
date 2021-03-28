module github.com/rancher/charts

go 1.16

replace (
	golang.stackrox.io/kube-linter v0.0.0-20210316191552-241e80db4436 => github.com/aiyengar2/kube-linter v0.0.0-20210328010234-9ca145667de4
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)

require (
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/swag v0.19.10 // indirect
	github.com/hashicorp/go-multierror v1.1.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	golang.stackrox.io/kube-linter v0.0.0-20210316191552-241e80db4436
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	helm.sh/helm/v3 v3.4.2
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kubectl v0.19.4
)
