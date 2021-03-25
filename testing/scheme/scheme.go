package scheme

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	cisv1 "github.com/rancher/cis-operator/pkg/apis/cis.cattle.io/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/scheme"
)

func Scheme() *runtime.Scheme {
	baseScheme := scheme.Scheme
	return AddCRDsToScheme(baseScheme)
}

func AddCRDsToScheme(s *runtime.Scheme) *runtime.Scheme {
	cisv1.AddToScheme(s)
	monitoringv1.AddToScheme(s)
	return s
}
